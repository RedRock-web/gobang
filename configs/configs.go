//@program: gobang
//@author: edte
//@create: 2020-06-05 21:13

// package configs 用于存放一些简单的类型，用于解决循环导包的问题

package configs

var (
	Uid    int
	RoomId int
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
	Id int `json:"id"`
}

type PlayFrom struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PasswdFrom struct {
	Rid      int `json:"rid"`
	Password int `json:"password"`
}

type MsgForm struct {
	Msg string
}
