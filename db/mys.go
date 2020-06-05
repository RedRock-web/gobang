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

func Start() {
	db, err := gorm.Open("mysql", "root:mima@/gobang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		logs.Error.Println(err)
	}
	//defer db.Close()

	if db.HasTable(&User{}) {
		db.AutoMigrate(&User{})
	} else {
		db.CreateTable(&User{})
	}

	MysqlClient = db
}
