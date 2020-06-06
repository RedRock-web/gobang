//@program: gobang
//@author: edte
//@create: 2020-06-05 21:13
package configs

var (
	Uid    int
	RoomId int
)

type LoginForm struct {
	Username string
	Password string
}

type InfoForm struct {
	Username string
	Age      int
	Gender   string
}

type RoomForm struct {
	Id int
}
