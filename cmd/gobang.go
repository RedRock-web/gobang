//@program: gobang
//@author: edte
//@create: 2020-06-05 20:48
package main

import (
	"gobang/db"
	"gobang/router"
)

func main() {
	db.Start()
	router.Start()
}
