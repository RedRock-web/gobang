//@program: gobang
//@author: edte
//@create: 2020-06-05 20:51

// package db 简单实现了一些对数据库的使用
package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gobang/logs"
)

var MysqlClient = &gorm.DB{}

//User 表示一个用户
type User struct {
	gorm.Model
	Uid      int
	Age      int
	Gender   string
	Username string
	Password string
}

//Room 表示一个房间
type Room struct {
	gorm.Model
	Rid           int // 房间 id
	Password      int // 密码
	Owner         int // 开房的玩家，默认为黑手
	AnotherPlayer int // 另一个玩家
	PlayerBlack   int // 黑手玩家
	PlayerWhite   int // 白手玩家
}

//Round 表示一个回合
type Round struct {
	gorm.Model
	Steps int // 回合数
	Rid   int // 房间 id
	Uid   int // 玩家 id
	X     int // 横坐标
	Y     int // 纵坐标
}

//Message
type Message struct {
	gorm.Model
	Rid int
	Uid int
	Msg string
}

//Start 初始化 mysql 数据库
func Start() {
	db, err := gorm.Open("mysql", "root:mima@/gobang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		logs.Error.Println(err)
	}
	//defer db.Close()

	if db.HasTable(&Room{}) {
		db.AutoMigrate(&Room{})
	} else {
		db.CreateTable(&Room{})
	}

	if db.HasTable(&User{}) {
		db.AutoMigrate(&User{})
	} else {
		db.CreateTable(&User{})
	}

	if db.HasTable(&Round{}) {
		db.AutoMigrate(&Round{})
	} else {
		db.CreateTable(&Round{})
	}

	if db.HasTable(&Message{}) {
		db.AutoMigrate(&Message{})
	} else {
		db.CreateTable(&Message{})
	}

	MysqlClient = db
}
