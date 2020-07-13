//@program: gobang
//@author: edte
//@create: 2020-06-05 21:38
package gobang

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"gobang/configs"
	"gobang/db"
	"gobang/logs"
	"gobang/response"
)

func StartGame(c *gin.Context) {
	if RoomList.Rooms[configs.RoomId].IsStartGame() {
		response.OkWithData(c, "game has start!")
		return
	}

	if !RoomList.Rooms[configs.RoomId].IsFullPlayer() {
		response.OkWithData(c, "player not full!")
		return
	}

	if !RoomList.Rooms[configs.RoomId].IsAllReady() {
		response.OkWithData(c, "player not all prepared")
		return
	}

	RoomList.Rooms[configs.RoomId].startGame()

	if RoomList.Rooms[configs.RoomId].playingUid == configs.Uid {
		response.OkWithData(c, "The game starts successfully, you are the white pawn, "+
			"please  perform your turn")
	} else {
		response.OkWithData(c, "The game starts successfully, you are the white pawn, "+
			"please wait for the opponent to perform its turn")
	}
}

//Play 用于开启游戏，玩家下棋，观众则默认观察最新棋局
func Play(c *gin.Context) {
	if RoomList.Rooms[configs.RoomId].IsSpectators() {
		ViewStatus(c)
		return
	}
	PlayChess(c)
}

//PlayChess 用于玩家下棋
func PlayChess(c *gin.Context) {
	if !RoomList.Rooms[configs.RoomId].start {
		response.OkWithData(c, "game not start!")
		return
	}

	p := configs.PlayFrom{}

	if err := c.ShouldBindBodyWith(&p, binding.JSON); err != nil {
		logs.Error.Println(err)
		response.FormError(c)
		return
	}

	if configs.Uid != RoomList.Rooms[configs.RoomId].playingUid {
		response.OkWithData(c, "Your round is over！")
		return
	}

	if !RoomList.Rooms[configs.RoomId].board.IsNullCell(p.X, p.Y) {
		response.OkWithData(c, "cell has Occupied!")
		return
	}

	if RoomList.Rooms[configs.RoomId].winer != 0 {
		if configs.Uid == RoomList.Rooms[configs.RoomId].winer {
			response.OkWithData(c, "you are the winner")
		} else {
			response.OkWithData(c, "you lost!")
		}
	}

	RoomList.Rooms[configs.RoomId].board.playChess(p.X, p.Y)
	RoomList.Rooms[configs.RoomId].board.PrintStatus()
	response.OkWithData(c, "you are win!")
	if RoomList.Rooms[configs.RoomId].board.checkWin() {

		RoomList.Rooms[configs.RoomId].winer = configs.Uid
	}

	//response.OkWithData(c, RoomList.rooms[configs.RoomId].board.GetStatusByDb())
	ViewStatus(c)
}

//ViewStatus
func ViewStatus(c *gin.Context) {
	c.String(200, RoomList.Rooms[configs.RoomId].board.GetStatus())
}

//GetChats
func GetChats() (data []db.Message) {
	db.MysqlClient.Where("rid = ?", configs.RoomId).Find(&data)
	return
}

//Chat
func Chat(c *gin.Context) {
	var m configs.MsgForm

	if err := c.BindWith(&m, binding.JSON); err != nil {
		response.FormError(c)
		return
	}

	RoomList.Rooms[configs.RoomId].SetMsg(m.Msg)

	response.OkWithData(c, GetChats())
}

//Peace
func Peace(c *gin.Context) {
	RoomList.Rooms[configs.RoomId].SetPeaceFlag()

	if !RoomList.Rooms[configs.RoomId].IsAllPeace() {
		response.OkWithData(c, "Waiting for another player to agree")
		return
	}

	RoomList.Rooms[configs.RoomId].start = false
	response.OkWithData(c, "Successful summation！")
}

//Confess
func Confess(c *gin.Context) {
	RoomList.Rooms[configs.RoomId].SetConfessFlag()

	if !RoomList.Rooms[configs.RoomId].IsAllConfess() {
		response.OkWithData(c, "Waiting for another player to agree")
		return
	}

	RoomList.Rooms[configs.RoomId].start = false
	response.OkWithData(c, "Successful confess！")
}

//GetHistory
func GetHistory(c *gin.Context) {
	for _, v := range configs.History {
		c.String(200, v)
	}
}
