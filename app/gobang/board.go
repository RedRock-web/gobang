//@program: gobang
//@author: edte
//@create: 2020-06-05 21:38
package gobang

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
	color     int             //最近下棋的颜色
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
	board.cells[board.lastStepX][board.lastStepY] = board.color
}

//CheckWin 判断哪一个玩家胜利了
//func (board *Board) CheckWin() int {
//	count := 0
//	winFlag := 1
//	i := board.lastStepX - 1
//	j := board.lastStepY
//
//	for ; i >= 0 && count < Num; i-- {
//		count++
//		if board.cells[i][j] == board.color {
//			winFlag++
//		} else {
//			break
//		}
//	}
//
//}

func (board *Board) getTimes(cx, cy, dx, dy, c int) int {
	if c == Empty {
		return 0
	}
	if dx == 0 && dy == 0 {
		return 0
	}
	times := 0
	for i := 1; i <= 5; i++ {
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

func (board *Board) checkWin(x, y, d int, color int) bool {
	if d == Empty {
		return false
	}
	if (board.getTimes(x, y, 0, 1, d)+board.getTimes(x, y, 0, -1, d)) >= 4 || (board.getTimes(x, y, 1, 0, d)+board.getTimes(x, y, -1, 0, d)) >= 4 || (board.getTimes(x, y, 1, 1, d)+board.getTimes(x, y, -1, -1, d)) >= 4 || (board.getTimes(x, y, 1, -1, d)+board.getTimes(x, y, -1, 1, d)) >= 4 {
		return d == color
	}
	return false
}
