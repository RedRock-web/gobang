//@program: gobang
//@author: edte
//@create: 2020-06-05 20:51
package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gobang/logs"
)

var MysqlClient = &gorm.DB{}

type User struct {
	gorm.Model
	Uid      int
	Age      int
	Gender   string
	Username string
	Password string
}

type Room struct {
	gorm.Model
	Rid           int // 房间 id
	Owner         int // 开房的玩家，默认为黑手
	AnotherPlayer int // 另一个玩家
	PlayerBlack   int // 黑手玩家
	PlayerWhite   int // 白手玩家
}

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

	MysqlClient = db
}
