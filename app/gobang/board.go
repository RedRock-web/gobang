//@program: gobang
//@author: edte
//@create: 2020-06-05 21:38

// package gobang 实现了五子棋功能

package gobang

import (
	"fmt"

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
	RoomList.Rooms[configs.RoomId].steps++

	err := db.MysqlClient.Create(&db.Round{
		Steps: RoomList.Rooms[configs.RoomId].steps,
		Rid:   configs.RoomId,
		Uid:   configs.Uid,
		X:     x,
		Y:     y,
	}).Error

	if err != nil {
		logs.Error.Println(err)
	}

	RoomList.Rooms[configs.RoomId].board.lastStepX = x
	RoomList.Rooms[configs.RoomId].board.lastStepY = y
	RoomList.Rooms[configs.RoomId].board.lastColor = RoomList.Rooms[configs.RoomId].playingColor

	board.setCells()

	RoomList.Rooms[configs.RoomId].holding = !RoomList.Rooms[configs.RoomId].holding

	if RoomList.Rooms[configs.RoomId].holding {
		RoomList.Rooms[configs.RoomId].playingColor = Black
		RoomList.Rooms[configs.RoomId].playingUid = RoomList.Rooms[configs.RoomId].playerBlack
	} else {
		RoomList.Rooms[configs.RoomId].playingColor = White
		RoomList.Rooms[configs.RoomId].playingUid = RoomList.Rooms[configs.RoomId].playerWhite
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

//GetStatusByDb 按时间/步骤获取下棋信息
func (board *Board) GetStatusByDb() []db.Round {
	var b []db.Round

	if err := db.MysqlClient.Where("rid = ?", configs.RoomId).Find(&b).Error; err != nil {
		logs.Error.Println(err)
	}

	return b
}

//IsNullCell 用于判断该坐标是否为空
func (board *Board) IsNullCell(x int, y int) bool {
	return board.cells[x][y] == Empty
}

//PrintStatus 在控制台输出实时下棋信息
func (board *Board) PrintStatus() {
	for i := 0; i < Size; i++ {
		for j := 0; j < Size; j++ {
			if board.cells[i][j] == 1 {
				fmt.Printf("* ")
			} else if board.cells[i][j] == 2 {
				fmt.Print("# ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}
}

func (board *Board) GetStatus() (s string) {
	for i := 0; i < Size; i++ {
		for j := 0; j < Size; j++ {
			if board.cells[i][j] == 1 {
				s = s + "* "
			} else if board.cells[i][j] == 2 {
				s = s + "# "
			} else {
				s = s + "0 "
			}
		}
		s = s + "\n"
	}
	configs.History = append(configs.History, s)
	return
}
