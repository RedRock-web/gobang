//@program: gobang
//@author: edte
//@create: 2020-06-05 21:38
package gobang

import (
	"gobang/configs"
	"gobang/db"
	"gobang/logs"
)

const (
	Empty = 0  //表示棋盘该坐标为空
	Size  = 15 //棋盘大小
	Black = 1  //表示黑色棋子
	White = 2  //表示白色棋子
	Num   = 5
)

//Board 表示棋盘
type Board struct {
	lastStepX int             //最近下棋的横坐标
	lastStepY int             //最近下棋的纵坐标
	lastColor int             //最近下棋的颜色
	cells     [Size][Size]int //棋盘各坐标
}

//NewBoard 新生成一个棋盘
func NewBoard() *Board {
	var cell [Size][Size]int

	for i := 0; i < Size; i++ {
		for j := 0; j < Size; j++ {
			cell[i][j] = Empty
		}
	}
	return &Board{cells: cell}
}

func (board *Board) setCells() {
	board.cells[board.lastStepX][board.lastStepY] = board.lastColor
}

func (board *Board) playChess(x int, y int) {
	roomList.rooms[configs.RoomId].steps++

	err := db.MysqlClient.Create(&db.Round{
		Steps: roomList.rooms[configs.RoomId].steps,
		Rid:   configs.RoomId,
		Uid:   configs.Uid,
		X:     x,
		Y:     y,
	}).Error

	if err != nil {
		logs.Error.Println(err)
	}

	roomList.rooms[configs.RoomId].board.lastStepX = x
	roomList.rooms[configs.RoomId].board.lastStepY = y
	roomList.rooms[configs.RoomId].board.lastColor = roomList.rooms[configs.RoomId].playingColor

	board.setCells()

	roomList.rooms[configs.RoomId].holding = !roomList.rooms[configs.RoomId].holding

	if roomList.rooms[configs.RoomId].holding {
		roomList.rooms[configs.RoomId].playingColor = Black
		roomList.rooms[configs.RoomId].playingUid = roomList.rooms[configs.RoomId].playerBlack
	} else {
		roomList.rooms[configs.RoomId].playingColor = White
		roomList.rooms[configs.RoomId].playingUid = roomList.rooms[configs.RoomId].playerWhite
	}
}

// getTimes 根据初始坐标和最近的一次坐标获取判断是否连子
func (board *Board) getTimes(cx, cy, dx, dy, c int) int {
	if c == Empty {
		return 0
	}
	if dx == 0 && dy == 0 {
		return 0
	}
	times := 0
	for i := 1; i <= Num; i++ {
		nx := cx + (dx * i)
		ny := cy + (dy * i)
		if nx < 0 || ny < 0 || nx >= len(board.cells) || ny >= len(board.cells[0]) {
			continue
		}
		nc := board.cells[nx][ny]
		if nc == Empty || c != nc {
			break
		}
		times++
	}
	return times
}

// checkWin 检查最后一次下的棋是否胜利
func (board *Board) checkWin() bool {
	return (board.getTimes(board.lastStepX, board.lastStepY, 0, 1, board.lastColor)+board.getTimes(board.lastStepX, board.lastStepY, 0, -1, board.lastColor)) >= 4 || (board.getTimes(board.lastStepX, board.lastStepY, 1, 0, board.lastColor)+board.getTimes(board.lastStepX, board.lastStepY, -1, 0, board.lastColor)) >= 4 || (board.getTimes(board.lastStepX, board.lastStepY, 1, 1, board.lastColor)+board.getTimes(board.lastStepX, board.lastStepY, -1, -1, board.lastColor)) >= 4 || (board.getTimes(board.lastStepX, board.lastStepY, 1, -1, board.lastColor)+board.getTimes(board.lastStepX, board.lastStepY, -1, 1, board.lastColor)) >= 4
}

func (board *Board) GetStatus() []db.Round {
	var b []db.Round

	if err := db.MysqlClient.Where("rid = ?", configs.RoomId).Find(&b).Error; err != nil {
		logs.Error.Println(err)
	}

	return b
}

func (board *Board) IsNullCell() bool {
	logs.Info.Println(board.lastStepX, board.lastStepY)
	logs.Info.Println(board.cells[board.lastStepX][board.lastStepY])
	logs.Info.Println(board.cells[board.lastStepX][board.lastStepY] == Empty)
	return board.cells[board.lastStepX][board.lastStepY] == Empty
}
