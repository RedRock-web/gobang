//@program: gobang
//@author: edte
//@create: 2020-06-05 20:49
package middleware

import (
	"github.com/gin-gonic/gin"
	"gobang/app/user"
	"gobang/configs"
	"gobang/jwts"
	"gobang/response"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var i configs.InfoForm

		if err := c.ShouldBindJSON(&i); err != nil {
			response.FormError(c)
			return
		}

		token, err := c.Cookie(user.CookieName)

		if err != nil {
			panic(err)
			return
		}

		j := jwts.NewJwt()
		_, err = j.Check(token, jwts.Key)

		if err != nil {
			c.Abort()
			response.Error(c, 10002, "token is error!")
		}
	}
}
