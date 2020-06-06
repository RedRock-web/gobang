//@program: gobang
//@author: edte
//@create: 2020-06-05 21:04

// package router 用于初始化路由
package router

import (
	"github.com/gin-gonic/gin"
	"gobang/app/gobang"
	"gobang/app/user"
	"gobang/logs"
	"gobang/middleware"
)

//Start
func Start() {
	r := gin.Default()
	SetRouter(r)
	err := r.Run()
	logs.Error.Println(err)
}

//SetRouter 设置路由
func SetRouter(r *gin.Engine) {
	users := r.Group("")
	{
		users.POST("/register", user.Register)
		users.POST("/login", user.Login)
		users.POST("/get", user.GetInfo)
		users.POST("/modify", middleware.LoginAuth(), user.ModifyInfo)
	}

	room := r.Group("/room", middleware.LoginAuth(), middleware.GetUid())
	{
		room.POST("/create", middleware.HasJoinRoom(), gobang.CreateRoom)
		room.POST("/join", middleware.HasJoinRoom(), middleware.PasswdAuth(), gobang.JoinRoom)
		room.POST("/exit", middleware.NeedJoinRoom(), gobang.ExitRoom)
		room.POST("/ready", middleware.NeedJoinRoom(), gobang.Ready)
		room.POST("/password", middleware.NeedJoinRoom(), gobang.Password)
		room.POST("/chat", middleware.NeedJoinRoom(), gobang.Chat)
	}

	game := r.Group("/game", middleware.LoginAuth(), middleware.GetUid(), middleware.NeedJoinRoom())
	{
		game.POST("/start", gobang.StartGame)
		game.POST("/play", middleware.CheckPlayerSize(), middleware.HasStartGame(), gobang.Play)
		game.POST("/peace", middleware.HasStartGame())
		game.POST("/confess", middleware.HasStartGame())
		game.POST("/regret", middleware.HasStartGame(), gobang.Regret)
		room.POST("/chat", middleware.HasStartGame(), gobang.Chat)
	}
}
