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
	{
		r.POST("/register", user.Register)
		r.POST("/login", user.Login)
		r.POST("/get-info", user.GetInfo)
		r.POST("/modify-info", middleware.Auth(), user.ModifyInfo)
	}

}
