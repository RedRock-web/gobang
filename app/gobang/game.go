//@program: gobang
//@author: edte
//@create: 2020-06-05 21:38
package gobang

import (
	"github.com/gin-gonic/gin"
	"gobang/configs"
	"gobang/response"
)

func StartGame(c *gin.Context) {
	r := configs.RoomForm{}
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FormError(c)
		return
	}
	if roomList.rooms[r.Id].canStart() {

	}
}
