//@program: gobang
//@author: edte
//@create: 2020-06-05 21:38
package gobang

import (
	"time"
)

// Room 表示一个房间
type Room struct {
	roomId        int  // 房间 id
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
	room := &Room{
		roomId:      GetRandomRoomId(),
		owner:       uid,
		playerBlack: 0,
		playerWhite: 0,
		holding:     false,
		board:       NewBoard(),
		steps:       0,
	}
	roomList.rooms[room.roomId] = room
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

//SetPlayerBlack 用于设置黑色玩家，默认为开房玩家
func (room *Room) SetPlayerBlack() {
	room.playerBlack = room.owner
}

//SetPlayerBlack 用于设置白色玩家
func (room *Room) SetPlayerWhite() {
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
	roomList.rooms[room.roomId] = nil
}

//IsRoomExsit 用于判断房间是否存在
func IsRoomExsit(roomId int) bool {
	return roomList.rooms[roomId] == nil
}

//JoinPlayer 用于玩家加入房间
func JoinPlayer(roomId int, uid int) string {
	if IsRoomExsit(roomId) {
		return "room not exsit!"
	}

	if roomList.rooms[roomId].IsOpen() {
		return "room already open!"
	}
	if roomList.rooms[roomId].IsFullPlayer() {
		return "room already full player!"
	}

	roomList.rooms[roomId].anotherPlayer = uid
	return "successful!"
}
