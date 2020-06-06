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
	"strconv"
	"time"
)

//RoomList 是房间列表
type RoomLists struct {
	Rooms map[int]*Room
}

// RoomList 是房间列表实例
var RoomList = &RoomLists{
	Rooms: make(map[int]*Room),
}

// Room 表示一个房间
type Room struct {
	id            int          // 房间 id
	owner         int          // 开房的玩家，默认为黑手
	anotherPlayer int          // 另一个玩家
	playerBlack   int          // 黑手玩家
	playerWhite   int          // 白手玩家
	holding       bool         // 当前该谁下，true 为黑色，默认为黑色
	playingColor  int          // 当前下的玩家颜色
	playingUid    int          // 当前下的玩家 id
	board         *Board       // 棋盘
	steps         int          // 下棋步数
	ready         map[int]bool // 玩家准备状态，true 表示已准备
	open          bool         // 表示房间是否已开
	start         bool         // 表示是否游戏已经开始
	winer         int          // 胜利者
	password      int          // 房间密码
	spectators    map[int]bool // 观众
	msg           string       // 发的消息
}

// NewRoom 在当前房间列表中新建一个房间
func NewRoom(uid int) *Room {
	room := &Room{
		id:            GetRandomRoomId(),
		owner:         uid,
		anotherPlayer: 0,
		playerBlack:   0,
		playerWhite:   0,
		holding:       true,
		playingColor:  Black,
		playingUid:    uid,
		board:         NewBoard(),
		steps:         0,
		ready:         map[int]bool{},
		open:          false,
		start:         false,
		winer:         0,
		password:      0,
		spectators:    map[int]bool{},
		msg:           "",
	}
	RoomList.Rooms[room.id] = room

	db.MysqlClient.Create(&db.Room{
		Rid:           room.id,
		Owner:         uid,
		AnotherPlayer: 0,
		PlayerBlack:   0,
		PlayerWhite:   0,
	})
	return room
}

//SetPassword 设置房间密码
func (room *Room) SetPassword(password int) {
	room.password = password
}

//IsPasswdOk 判断房间密码是否正确
func (room *Room) IsPasswdOk(passwd int) bool {
	return room.password == passwd
}

//Regret 用于玩家悔棋
func (room *Room) Regret() {
	room.board.cells[room.board.lastStepX][room.board.lastStepY] = Empty
}

//IsSpectators 用于判断当前用户是否是观众
func (room *Room) IsSpectators() bool {
	return room.spectators[configs.Uid]
}

//SetMsg
func (room *Room) SetMsg(msg string) {
	room.msg = msg
	db.MysqlClient.Create(&db.Message{
		Rid: configs.RoomId,
		Uid: configs.Uid,
		Msg: msg,
	})
}

func (room *Room) HavePassword() bool {
	return room.password == 0
}

//GetRandomRoomId  获取当前时间戳为房间 id
func GetRandomRoomId() int {
	return int(time.Now().Unix())
}

func (room *Room) IsStart() bool {
	return room.start
}

// canStart 判断是否能够开始房间
func (room *Room) canStart() bool {
	return room.IsAllReady() && room.playerBlack != 0 && room.playerWhite != 0
}

//AddSpectators 增加观众
func (room *Room) AddSpectators(uid int) {
	room.spectators[uid] = true
}

// Ready 用于玩家切换准备状态
func (room *Room) Ready(uid int) {
	a := room.ready[uid]
	room.ready[uid] = !a
}

//HasAllReady 用于判断所有玩家是否都已经准备
func (room *Room) IsAllReady() bool {
	return room.ready[room.playerBlack] && room.ready[room.playerWhite]
}

//SetAnotherPlayer 用于设置另一个玩家
func (room *Room) SetAnotherPlayer() {
	db.MysqlClient.Model(&db.Room{}).Where("rid = ?", configs.RoomId).Update("another_player", configs.Uid)
	RoomList.Rooms[room.id].anotherPlayer = configs.Uid
}

//SetPlayerBlack 用于设置黑色玩家，默认为开房玩家,默认未准备
func (room *Room) SetPlayerBlack() {
	room.ready[room.owner] = false
	db.MysqlClient.Model(&db.Room{}).Where("rid = ?", configs.RoomId).Update("player_black", room.owner)
	room.playerBlack = room.owner
}

//SetPlayerBlack 用于设置白色玩家
func (room *Room) SetPlayerWhite() {
	room.ready[room.anotherPlayer] = false
	db.MysqlClient.Model(&db.Room{}).Where("rid = ?", configs.RoomId).Update("player_white", room.anotherPlayer)
	room.playerWhite = room.anotherPlayer
}

//IsOpen 用于判断房间是否已开
func (room *Room) IsOpen() bool {
	return room.open
}

func (room *Room) IsStartGame() bool {
	return room.start
}

//IsFullPlayer 用于判断是否玩家已经满员
func (room *Room) IsFullPlayer() bool {
	return room.owner != 0 && room.anotherPlayer != 0
}

// startGame 用于开启游戏
func (room *Room) startGame() {
	if !room.canStart() {
		return
	}
	room.steps = 0
	RoomList.Rooms[configs.RoomId].start = true
}

//gameOver 结束游戏，并关闭房间
func (room *Room) gameOver() {
	if !room.open {
		return
	}
	room.open = false
	RoomList.Rooms[room.id] = nil
}

//IsRoomExist 用于判断房间是否存在
func IsRoomExist(roomId int) bool {
	return RoomList.Rooms[roomId] == nil
}

//JoinRoom 用于玩家加入房间，满员则加入观众
func JoinRoom(c *gin.Context) {
	if RoomList.Rooms[configs.RoomId].IsFullPlayer() {
		JoinSpectators(c)
		return
	}

	JoinPlayer(c)
}

//JoinSpectators 用于玩家加入观众
func JoinSpectators(c *gin.Context) {
	RoomList.Rooms[configs.RoomId].AddSpectators(configs.Uid)
}

//JoinPlayer 用于玩家加入房间
func JoinPlayer(c *gin.Context) {
	r := configs.RoomForm{}

	if err := c.ShouldBindJSON(&r); err != nil {
		response.FormError(c)
		return
	}

	if IsRoomExist(r.Id) {
		response.OkWithData(c, "room not exist!")
		return
	}

	if RoomList.Rooms[r.Id].IsOpen() {
		response.OkWithData(c, "room already open!")
		return
	}

	if RoomList.Rooms[r.Id].IsFullPlayer() {
		response.OkWithData(c, "room already full player!")
		return
	}

	RoomList.Rooms[r.Id].SetAnotherPlayer()
	RoomList.Rooms[r.Id].SetPlayerWhite()

	response.OkWithData(c, "successful!")
}

//CreateRoom 用于创建房间
func CreateRoom(c *gin.Context) {
	room := NewRoom(configs.Uid)

	if IsRoomExist(room.id) {
		response.Error(c, 10003, "room is exist!")
		return
	}

	configs.RoomId = room.id

	room.SetPlayerBlack()

	response.OkWithData(c, "create room successful,room id is "+strconv.Itoa(room.id))
}

func ExitRoom(c *gin.Context) {
	r := configs.RoomForm{}
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FormError(c)
		return
	}

	RoomList.Rooms[r.Id].gameOver()
	response.OkWithData(c, "Successfully left the room！")
}

func Ready(c *gin.Context) {
	red := RoomList.Rooms[configs.RoomId].ready[configs.Uid]

	RoomList.Rooms[configs.RoomId].Ready(configs.Uid)

	if !red == true {
		response.OkWithData(c, "already prepared")
		if RoomList.Rooms[configs.RoomId].IsAllReady() {
			response.OkWithData(c, "all player has prepared,can start the game!")
		} else {
			response.OkWithData(c, "another player is not prepare!")
		}
	} else {
		response.OkWithData(c, "already cancel prepare")
	}
}

func Password(c *gin.Context) {
	var p configs.PasswdFrom

	if err := c.BindWith(&p, binding.JSON); err != nil {
		logs.Error.Println(err)
		return
	}

	if RoomList.Rooms[configs.RoomId].owner != configs.Uid {
		response.OkWithData(c, "Non-homeowner, unable to change password！")
		return
	}

	RoomList.Rooms[configs.RoomId].SetPassword(p.Password)

	response.OkWithData(c, "set password successfully!")
}

// RoomChat
func RoomChat(c *gin.Context) {
	var m configs.MsgForm

	if err := c.BindWith(&m, binding.JSON); err != nil {
		response.FormError(c)
		return
	}

	RoomList.Rooms[configs.RoomId].SetMsg(m.Msg)
}

//Regret
func Regret(c *gin.Context) {
	RoomList.Rooms[configs.RoomId].Regret()
}
