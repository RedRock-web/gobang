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

//HasJoinRoom
func HasJoinRoom() gin.HandlerFunc {
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

func NeedJoinRoom() gin.HandlerFunc {
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
		if !gobang.RoomList.Rooms[configs.Uid].IsStart() {
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

		if err := c.BindWith(&p, binding.JSON); err != nil {
			response.FormError(c)
			logs.Error.Println(err)
			c.Abort()
			return
		}

		if !gobang.RoomList.Rooms[configs.RoomId].HavePassword() {
			c.Next()
			return
		}

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
