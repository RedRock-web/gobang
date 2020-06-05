//@program: gobang
//@author: edte
//@create: 2020-06-05 21:38
package gobang

const (
	Empty = 0  //表示棋盘该坐标为空
	Size  = 15 //棋盘大小
	Black = 1  //表示黑色棋子
	White = 2  //表示白色棋子
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
func (board *Board) CheckWin() int {
}
