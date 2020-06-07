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
	peaceFlag     map[int]bool // 求和
	confessFlag   map[int]bool // 认输
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
		peaceFlag:     map[int]bool{},
		confessFlag:   map[int]bool{},
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

//IsBlackPlayer 用于当前用用户判断是否是黑色玩家 
func (room *Room) IsBlackPlayer() bool {
	return configs.Uid == room.playerBlack
}

//SetPeaceFlag 用于同意求和
func (room *Room) SetPeaceFlag() {
	room.peaceFlag[configs.Uid] = true
}

//SetConfessFlag 用于同意求和
func (room *Room) SetConfessFlag() {
	room.confessFlag[configs.Uid] = true
}

//IsAllPeace 用于判断是否两个玩家都同意和平
func (room *Room) IsAllPeace() bool {
	return room.peaceFlag[room.playerBlack] && room.peaceFlag[room.playerBlack]
}

//IsAllConfess 用于判断是否两个玩家都同意和平
func (room *Room) IsAllConfess() bool {
	return room.confessFlag[room.playerBlack] && room.confessFlag[room.playerBlack]
}

//IsWhitePlayer 用于当前用用户判断是否是白色玩家
func (room *Room) IsWhitePlayer() bool {
	return configs.Uid == room.playerWhite
}

//Regret 用于玩家悔棋
func (room *Room) Regret() {
	room.board.cells[room.board.lastStepX][room.board.lastStepY] = Empty
}

//IsSpectators 用于判断当前用户是否是观众
func (room *Room) IsSpectators() bool {
	return room.spectators[configs.Uid]
}

//IsPlayers 用于判断当前用户是否是玩家
func (room *Room) IsPlayers() bool {
	return configs.Uid == room.playerBlack || configs.Uid == room.playerWhite
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

//gameOver 结束游戏
func (room *Room) gameOver() {
	room.start = false
}

//CloseRoom 关闭房间
func (room *Room) CloseRoom() {
	room.open = false
}

//ExitRoom 退出房间
func (room *Room) ExitRoom() {
	RoomList.Rooms[room.id] = nil
}

//IsOwner 判断用户是否是房间拥有者
func (room *Room) IsOwner() bool {
	return configs.Uid == room.owner
}

//IsRoomExist 用于判断房间是否存在
func IsRoomExist(roomId int) bool {
	return RoomList.Rooms[roomId] == nil //&Room{}
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

	if RoomList.Rooms[configs.RoomId].IsOpen() {
		response.OkWithData(c, "room already open!")
		return
	}

	if RoomList.Rooms[configs.RoomId].IsFullPlayer() {
		response.OkWithData(c, "room already full player!")
		return
	}

	RoomList.Rooms[configs.RoomId].SetAnotherPlayer()
	RoomList.Rooms[configs.RoomId].SetPlayerWhite()

	response.OkWithData(c, "successful!")
}

//CreateRoom 用于创建房间
func CreateRoom(c *gin.Context) {
	room := NewRoom(configs.Uid)
	configs.RoomId = room.id

	room.SetPlayerBlack()

	response.OkWithData(c, "create room successful,room id is "+strconv.Itoa(room.id))
}

func ExitRoom(c *gin.Context) {
	RoomList.Rooms[configs.RoomId].ExitRoom()
	response.OkWithData(c, "Successfully exit the room！")
}

func CloseRoom(c *gin.Context) {
	RoomList.Rooms[configs.RoomId].CloseRoom()
	response.OkWithData(c, "Successfully close the room！")
}

func Ready(c *gin.Context) {
	red := RoomList.Rooms[configs.RoomId].ready[configs.Uid]

	RoomList.Rooms[configs.RoomId].Ready(configs.Uid)

	if red {
		response.OkWithData(c, "already cancel prepare")
		return
	}

	if RoomList.Rooms[configs.RoomId].IsAllReady() {
		response.OkWithData(c, "already prepared, all player has prepared,can start the game!")
	} else {
		response.OkWithData(c, "already prepared, another player is not prepare!")
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
	err := db.MysqlClient.Delete(&db.Round{
		Rid: configs.RoomId,
		Uid: configs.Uid,
		X:   RoomList.Rooms[configs.RoomId].board.lastStepX,
		Y:   RoomList.Rooms[configs.RoomId].board.lastStepY,
	}).Error
	if err != nil {
		logs.Error.Println(err)
		return
	}
	ViewStatus(c)
}
