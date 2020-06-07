//@program: gobang
//@author: edte
//@create: 2020-06-05 20:58

// package user 实现了用户的一些功能
//
//
package user

import (
	"github.com/gin-gonic/gin"
	"gobang/configs"
	"gobang/db"
	"gobang/jwts"
	"gobang/logs"
	"gobang/response"
	"time"
)

const (
	CookieName = "jwts"
	host       = "localhost:8080"
)

//Register
func Register(c *gin.Context) {
	var l configs.LoginForm

	if err := c.ShouldBindJSON(&l); err != nil {
		response.FormError(c)
		return
	}

	if IsRegister(l) {
		response.Error(c, 10001, "user exsits!")
		return
	} else {
		_ = AddUser(l)
	}

	jwt := jwts.NewJwt()
	data, err := jwt.Create(l, jwts.Key)
	if err != nil {
		logs.Error.Println(err)
		return
	}

	c.SetCookie(CookieName, data, 10000, "/", host, false, true)
	response.OkWithData(c, "register successful!")
}

//Login
func Login(c *gin.Context) {
	var l configs.LoginForm

	if err := c.ShouldBindJSON(&l); err != nil {
		response.FormError(c)
		return
	}

	if !IsRegister(l) {
		response.Error(c, 10002, "not registered!")
		return
	}

	jwt := jwts.NewJwt()
	data, err := jwt.Create(l, jwts.Key)

	if err != nil {
		logs.Error.Println(err)
		return
	}

	if PasswordIsOk(l) {
		c.SetCookie(CookieName, data, 1000, "/", host, false, true)
		response.OkWithData(c, "login successful!")
	} else {
		response.Error(c, 10003, "password not right!")
	}
}

//PasswordIsOk
func PasswordIsOk(l configs.LoginForm) bool {
	var u db.User
	db.MysqlClient.Where(db.User{Username: l.Username, Password: l.Password}).First(&u)
	return u.Gender != ""
}

//IsRegister
func IsRegister(l configs.LoginForm) bool {
	var u db.User
	db.MysqlClient.Where("username = ?", l.Username).First(&u)
	return u.Gender == "male" || u.Gender == "female"
}

//AddUser
func AddUser(l configs.LoginForm) error {
	return db.MysqlClient.Create(&db.User{
		Uid:      int(time.Now().Unix()) - 10000,
		Age:      18,
		Gender:   "male",
		Username: l.Username,
		Password: l.Password,
	}).Error
}

//GetInfo
func GetInfo(c *gin.Context) {
	var i configs.InfoForm

	if err := c.ShouldBindJSON(&i); err != nil {
		response.FormError(c)
		return
	}

	l := configs.LoginForm{
		Username: i.Username,
		Password: "",
	}

	if !IsRegister(l) {
		response.Error(c, 10002, "not registered!")
		return
	}

	var u db.User

	if err := db.MysqlClient.Where("username = ?", i.Username).First(&u).Error; err != nil {
		logs.Error.Println(err)
		return
	}

	response.OkWithData(c, gin.H{
		"username": u.Username,
		"gender":   u.Gender,
		"age":      u.Age,
	})
}

//ModifyInfo
func ModifyInfo(c *gin.Context) {
	var i configs.InfoForm

	if err := c.ShouldBindJSON(&i); err != nil {
		response.FormError(c)
		return
	}

	l := configs.LoginForm{
		Username: i.Username,
		Password: "",
	}

	if !IsRegister(l) {
		response.Error(c, 10002, "not registered!")
		return
	}

	if err := db.MysqlClient.Model(&db.User{}).Update(map[string]interface{}{"age": i.Age,
		"gender": i.Gender}).Error; err != nil {
		logs.Error.Println(err)
		return
	}
	response.Ok(c)
}

//GetUidByUsername
func GetUidByUsername(username string) int {
	var u db.User

	db.MysqlClient.Where("username = ?", username).First(&u)
	return u.Uid
}
