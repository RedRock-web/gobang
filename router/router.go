//@program: gobang
//@author: edte
//@create: 2020-06-05 21:04
package router

import (
	"github.com/gin-gonic/gin"
	"gobang/app/user"
	"gobang/logs"
	"gobang/middleware"
)

func Start() {
	r := gin.Default()
	SetRouter(r)
	err := r.Run()
	logs.Error.Println(err)
}

func SetRouter(r *gin.Engine) {
	users := r.Group("")
	{
		users.POST("/register", user.Register)
		users.POST("/login", user.Login)
		users.POST("/get-info", user.GetInfo)
		users.POST("/modify-info", middleware.Auth(), user.ModifyInfo)
	}

	room := r.Group("", middleware.Auth())
	{
		room.POST("/create-room")
		room.POST("/join-room")
		room.POST("/out-room")
		room.POST("ready")
	}

	game := r.Group("")
	{
		game.POST("start-game")
		game.POST("end-game")
		game.POST("play-chess")
	}
}
