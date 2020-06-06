//@program: gobang
//@author: edte
//@create: 2020-06-06 14:38
package gobang

import (
	"github.com/gin-gonic/gin"
	"gobang/db"
	"testing"
)

func TestPlayChess(t *testing.T) {
	db.Start()
	r := gin.Default()
	r.POST("/", PlayChess)
	r.Run()
}
