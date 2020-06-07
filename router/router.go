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
	r.GET("/", func(c *gin.Context) {
		c.String(200, "看到这个说明网站能用")
	})
	r.Use(middleware.Cors())

	users := r.Group("/user")
	{
		users.POST("/register", user.Register)
		users.POST("/login", user.Login)
		users.POST("/get", user.GetInfo)
		users.POST("/modify", middleware.LoginAuth(), user.ModifyInfo)
	}

	room := r.Group("/room", middleware.LoginAuth(), middleware.GetUid())
	{
		room.POST("/create", middleware.NotJoinRoomAuth(), gobang.CreateRoom)
		room.POST("/join", middleware.NotJoinRoomAuth(), middleware.GetRoomId(), middleware.RoomExistAuth(), middleware.PasswdAuth(), gobang.JoinRoom)
		room.POST("/exit", middleware.JoinRoomAuth(), gobang.ExitRoom)
		room.POST("/close", middleware.JoinRoomAuth(), middleware.OwnerAuth(), gobang.CloseRoom)
		room.POST("/ready", middleware.JoinRoomAuth(), gobang.Ready)
		room.POST("/password", middleware.JoinRoomAuth(), gobang.Password)
		room.POST("/chat", middleware.JoinRoomAuth(), gobang.Chat)
	}

	game := r.Group("/game", middleware.LoginAuth(), middleware.GetUid(), middleware.JoinRoomAuth())
	{
		game.POST("/start", gobang.StartGame)
		game.POST("/play", middleware.CheckPlayerSize(), middleware.HasStartGame(), middleware.PlayerAuth(), gobang.Play)
		game.POST("/peace", middleware.HasStartGame(), middleware.PlayerAuth(), gobang.Peace)
		game.POST("/confess", middleware.HasStartGame(), middleware.PlayerAuth(), gobang.Confess)
		game.POST("/regret", middleware.HasStartGame(), middleware.PlayerAuth(), gobang.Regret)
		game.POST("/chat", middleware.HasStartGame(), gobang.Chat)
		game.POST("/history", middleware.PlayerAuth(), gobang.GetHistory)
	}
}
