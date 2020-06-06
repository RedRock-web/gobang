//@program: gobang
//@author: edte
//@create: 2020-06-05 21:38
package gobang

import (
	"github.com/gin-gonic/gin"
	"gobang/configs"
	"gobang/db"
	"gobang/logs"
	"gobang/response"
	"strconv"
	"time"
)

// Room 表示一个房间
type Room struct {
	id            int  // 房间 id
	owner         int  // 开房的玩家，默认为黑手
	anotherPlayer int  // 另一个玩家
	playerBlack   int  // 黑手玩家
	playerWhite   int  // 白手玩家
	holding       bool // 当前该谁下， true for black, false for white
	//playing       int          // 当前下的玩家
	board *Board       // 板子
	steps int          // 下棋步数
	ready map[int]bool // 玩家准备状态，true 表示已准备
	open  bool         // 表示房间是否已开
}

//RoomList 是房间列表
type RoomList struct {
	rooms map[int]*Room
}

// roomList 是房间列表实例
var roomList = &RoomList{
	rooms: make(map[int]*Room),
}

// NewRoom 在当前房间列表中新建一个房间
func NewRoom(uid int) *Room {
	logs.Info.Println(uid)
	room := &Room{
		id:            GetRandomRoomId(),
		owner:         uid,
		anotherPlayer: 0,
		playerBlack:   0,
		playerWhite:   0,
		holding:       false,
		board:         NewBoard(),
		steps:         0,
		ready:         nil,
		open:          false,
	}
	roomList.rooms[room.id] = room

	db.MysqlClient.Create(&db.Room{
		Rid:           room.id,
		Owner:         uid,
		AnotherPlayer: 0,
		PlayerBlack:   0,
		PlayerWhite:   0,
	})
	return room
}

//GetRandomRoomId  获取当前时间戳为房间 id
func GetRandomRoomId() int {
	return int(time.Now().Unix())
}

// canStart 判断是否能够开始房间
func (room *Room) canStart() bool {
	return room.HasAllReady() && room.playerBlack != 0 && room.playerWhite != 0
}

// Ready 用于玩家切换准备状态
func (room *Room) Ready(uid int) {
	a := room.ready[uid]
	room.ready[uid] = !a
}

//HasAllReady 用于判断所有玩家是否都已经准备
func (room *Room) HasAllReady() bool {
	return room.ready[room.playerBlack] && room.ready[room.playerWhite]
}

//SetAnotherPlayer 用于设置另一个玩家
func (room *Room) SetAnotherPlayer() {
	db.MysqlClient.Model(&db.Room{}).Where("rid = ?", configs.RoomId).Update("another_player", configs.Uid)
	roomList.rooms[room.id].anotherPlayer = configs.Uid
}

//SetPlayerBlack 用于设置黑色玩家，默认为开房玩家
func (room *Room) SetPlayerBlack() {
	db.MysqlClient.Model(&db.Room{}).Where("rid = ?", configs.RoomId).Update("player_black", room.owner)
	room.playerBlack = room.owner
}

//SetPlayerBlack 用于设置白色玩家
func (room *Room) SetPlayerWhite() {
	db.MysqlClient.Model(&db.Room{}).Where("rid = ?", configs.RoomId).Update("player_white", room.anotherPlayer)
	room.playerWhite = room.anotherPlayer
}

//IsOpen 用于判断房间是否已开
func (room *Room) IsOpen() bool {
	return room.open
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
}

//gameOver 结束游戏，并关闭房间
func (room *Room) gameOver() {
	if !room.open {
		return
	}
	room.open = false
	roomList.rooms[room.id] = nil
}

//IsRoomExsit 用于判断房间是否存在
func IsRoomExsit(roomId int) bool {
	return roomList.rooms[roomId] == nil
}

//JoinPlayer 用于玩家加入房间
func JoinPlayer(c *gin.Context) {
	r := configs.RoomForm{}
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FormError(c)
		return
	}

	if IsRoomExsit(r.Id) {
		response.OkWithData(c, "room not exsit!")
		return
	}

	if roomList.rooms[r.Id].IsOpen() {
		response.OkWithData(c, "room already open!")
		return
	}
	if roomList.rooms[r.Id].IsFullPlayer() {
		response.OkWithData(c, "room already full player!")
		return
	}

	roomList.rooms[r.Id].SetAnotherPlayer()
	roomList.rooms[r.Id].SetPlayerWhite()

	response.OkWithData(c, "successful!")
}

//CreateRoom 用于创建房间
func CreateRoom(c *gin.Context) {
	room := NewRoom(configs.Uid)
	if IsRoomExsit(room.id) {
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

	roomList.rooms[r.Id].gameOver()
}

func Ready(c *gin.Context) {
	r := configs.RoomForm{}
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FormError(c)
		return
	}
	roomList.rooms[r.Id].Ready(configs.Uid)
}
