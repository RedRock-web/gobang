//@program: gobang
//@author: edte
//@create: 2020-06-05 20:49

// middleware 用于存放一些需要的中间件
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gobang/app/gobang"
	"gobang/app/user"
	"gobang/configs"
	"gobang/db"
	"gobang/jwts"
	"gobang/logs"
	"gobang/response"
	"net/http"
)

//Auth 用于登录鉴权
func LoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(user.CookieName)

		if err != nil {
			response.Error(c, 10003, "needed login!")
			c.Abort()
			return
		}

		j := jwts.NewJwt()
		l, err := j.Check(token, jwts.Key)

		if err != nil {
			response.Error(c, 10002, "token is error!")
			c.Abort()
			return
		}

		configs.Uid = user.GetUidByUsername(l.Username)
	}
}

//NotJoinRoomAuth
func NotJoinRoomAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var r db.Room

		db.MysqlClient.Where("owner = ? or another_player = ?", configs.Uid, configs.Uid).First(&r)

		if r.Rid != 0 {
			response.Error(c, 10004, "already joined a room!")
			c.Abort()
			return
		}
		c.Next()
	}
}

//JoinRoomAuth
func JoinRoomAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var r db.Room

		db.MysqlClient.Where("owner = ? or another_player = ?", configs.Uid, configs.Uid).First(&r)

		if r.Rid == 0 {
			response.Error(c, 10004, "needed join a room!")
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetUid() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(user.CookieName)
		if err != nil {
			logs.Error.Println(err)
			return
		}

		j := jwts.NewJwt()
		l, err := j.Check(token, jwts.Key)

		if err != nil {
			response.Error(c, 10002, "token is error!")
			c.Abort()
			return
		}
		configs.Uid = user.GetUidByUsername(l.Username)
	}
}

//HasStartGame 用于游戏开始简权
func HasStartGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !gobang.RoomList.Rooms[configs.RoomId].IsStart() {
			response.OkWithData(c, "The game does not start!")
			c.Abort()
			return
		}
		c.Next()
	}
}

//CheckPlayerSize
func CheckPlayerSize() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := configs.PlayFrom{}

		if err := c.ShouldBindBodyWith(&p, binding.JSON); err != nil {
			logs.Error.Println(err)
			response.FormError(c)
			c.Abort()
			return
		}

		if !(p.X >= 0 && p.Y >= 0 && p.X <= 15 && p.Y <= 15) {
			response.FormError(c)
			c.Abort()
			return
		}
	}
}

//PasswdAuth 用于鉴权房间密码
func PasswdAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var p configs.PasswdFrom

		if err := c.ShouldBindBodyWith(&p, binding.JSON); err != nil {
			response.FormError(c)
			logs.Error.Println(err)
			c.Abort()
			return
		}
		logs.Info.Println(gobang.RoomList.Rooms[configs.RoomId].HavePassword())

		if !gobang.RoomList.Rooms[configs.RoomId].HavePassword() {
			c.Next()
			return
		}
		logs.Info.Println(gobang.RoomList.Rooms[configs.RoomId].IsPasswdOk(p.Password))

		if !gobang.RoomList.Rooms[configs.RoomId].IsPasswdOk(p.Password) {
			response.Error(c, 10005, "password is not right!")
			c.Abort()
			return
		}

		c.Next()
	}
}

//PlayerAuth 用于鉴别是否玩家
func PlayerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !gobang.RoomList.Rooms[configs.RoomId].IsPlayers() {
			response.OkWithData(c, "you are not a players")
			c.Abort()
			return
		}
		c.Next()
	}
}

//Cors 解决跨域问题
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

//RoomExistAuth
func RoomExistAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//logs.Info.Println("Begin...")
		//logs.Info.Println(configs.RoomId)
		//logs.Info.Println(gobang.RoomList.Rooms[configs.RoomId])
		//logs.Info.Println(gobang.RoomList.Rooms[configs.RoomId] == nil)
		//logs.Info.Println(gobang.RoomList.Rooms[configs.RoomId] == &gobang.Room{})
		//logs.Info.Println("End...")
		if gobang.IsRoomExist(configs.RoomId) {
			response.OkWithData(c, "room not exist!")
			c.Abort()
			return
		}
	}
}

//GetRoomId
func GetRoomId() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := configs.RoomForm{}

		if err := c.ShouldBindBodyWith(&r, binding.JSON); err != nil {
			logs.Error.Println(err)
			response.FormError(c)
			c.Abort()
			return
		}

		configs.RoomId = r.Rid
		c.Next()
	}
}

//OwnerAuth
func OwnerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !gobang.RoomList.Rooms[configs.RoomId].IsOwner() {
			response.OkWithData(c, "you are not the room owner!")
			c.Abort()
			return
		}
		c.Next()
	}
}
