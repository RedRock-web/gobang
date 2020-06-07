//@program: gobang
//@author: edte
//@create: 2020-06-05 21:13

// package configs 用于存放一些简单的类型，用于解决循环导包的问题

package configs

var (
	Uid     int
	RoomId  int
	History []string
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type InfoForm struct {
	Username string `json:"username"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
}

type RoomForm struct {
	Rid int
}

type PlayFrom struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PasswdFrom struct {
	Rid      int
	Password string
}

type MsgForm struct {
	Msg string
}
