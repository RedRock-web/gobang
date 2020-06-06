//@program: gobang
//@author: edte
//@create: 2020-06-05 20:49

// package response 简单封装了一下 请求后的 response
package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ok(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 10000})
}

func FormError(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 10001, "message": "request form error!"})
}

func OkWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 10000, "message": "ok", "data": data})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{"code": code, "message": msg})
}
