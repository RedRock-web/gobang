//@program: gobang
//@author: edte
//@create: 2020-06-05 21:04
package router

import (
	"github.com/gin-gonic/gin"
	"gobang/app/gobang"
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
		users.POST("/modify-info", middleware.LoginAuth(), user.ModifyInfo)
	}

	room := r.Group("", middleware.LoginAuth())
	{
		room.POST("/create-room", middleware.GetUid(), middleware.HasJoinRoom(), gobang.CreateRoom)
		room.POST("/join-room", middleware.GetUid(), middleware.HasJoinRoom(), gobang.JoinPlayer)
		room.POST("/exit-room", middleware.NeedJoinRoom(), gobang.ExitRoom)
		room.POST("ready", middleware.NeedJoinRoom(), gobang.Ready)
	}

	game := r.Group("", middleware.LoginAuth())
	{
		game.POST("start-game")
		game.POST("end-game")
		game.POST("play-chess")
	}
}
