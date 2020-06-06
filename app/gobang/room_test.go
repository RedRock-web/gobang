//@program: gobang
//@author: edte
//@create: 2020-06-05 21:50
package gobang

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestA(t *testing.T) {
	r := gin.Default()

	room := r.Group("/rom")
	{
		room.GET("a")
		room.GET("b")
	}

	r.Run()

}
