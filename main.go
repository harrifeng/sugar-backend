package main

import (
	"db"
	"fmt"
	"server"
)

func main() {
	// init mysql
	mysqlDb, err := db.InitMysql()
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer func() {
		err = mysqlDb.Close()
		if err != nil {
			fmt.Printf("%s", err)
			return
		}
	}()

	// init redis
	redisPool := db.InitRedis()
	defer func() {
		err := redisPool.Close()
		if err != nil {
			fmt.Printf("%s", err)
			return
		}

	}()
	//server.Start()
	server.DatabaseTest()
}
