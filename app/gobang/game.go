//@program: gobang
//@author: edte
//@create: 2020-06-05 21:38
package gobang

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gobang/configs"
	"gobang/logs"
	"gobang/response"
)

func StartGame(c *gin.Context) {
	if roomList.rooms[configs.RoomId].IsStartGame() {
		response.OkWithData(c, "game has start!")
		return
	}

	if !roomList.rooms[configs.RoomId].IsFullPlayer() {
		response.OkWithData(c, "player not full!")
		return
	}

	if !roomList.rooms[configs.RoomId].IsAllReady() {
		response.OkWithData(c, "player not all prepared")
		return
	}

	roomList.rooms[configs.RoomId].startGame()
	response.OkWithData(c, "game start successful!")

	if roomList.rooms[configs.RoomId].playingUid == configs.Uid {
		response.OkWithData(c, "you are the black player, Please start your round")
	} else {
		response.OkWithData(c, "you are the white player, Please wait the balck player start her round")
	}
}

//PlayChess 用于玩家下棋
func PlayChess(c *gin.Context) {
	if !roomList.rooms[configs.RoomId].start {
		response.OkWithData(c, "game not start!")
		return
	}

	p := configs.PlayFrom{}

	if err := c.ShouldBindBodyWith(&p, binding.JSON); err != nil {
		logs.Error.Println(err)
		response.FormError(c)
		return
	}

	logs.Info.Println(p)

	if configs.Uid != roomList.rooms[configs.RoomId].playingUid {
		response.OkWithData(c, "Your round is over！")
		return
	}

	if !roomList.rooms[configs.RoomId].board.IsNullCell() {
		response.OkWithData(c, "cell has Occupied!")
		return
	}

	roomList.rooms[configs.RoomId].board.playChess(p.X, p.Y)

	if roomList.rooms[configs.RoomId].board.checkWin() {
		response.OkWithData(c, "you are win!")
		return
	}

	data := roomList.rooms[configs.RoomId].board.GetStatus()

	response.OkWithData(c, data)
}
