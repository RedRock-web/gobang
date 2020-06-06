//@program: gobang
//@author: edte
//@create: 2020-06-05 20:49
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

func HasStartGame() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func CheckPlayeSize() gin.HandlerFunc {
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
