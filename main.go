package main

import (
	"Db"
	"fmt"
)

func DbTest() {
	Db.AutoCreateTableTest()
}

func main() {
	// init mysql
	db, err := Db.InitMysql()
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Printf("%s", err)
			return
		}
	}()

	// init redis
	redisPool := Db.InitRedis()
	defer func() {
		err := redisPool.Close()
		if err != nil {
			fmt.Printf("%s", err)
			return
		}

	}()

	//Server.Start()
	Db.SetNewVerificationCodeTest()
}
