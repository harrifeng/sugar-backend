package main

import (
	"db"
	"fmt"
	"server"
)

func DbTest() {
	db.AutoCreateTableTest()
}

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
	server.Start()
	//conn :=redisPool.Get()
	//server.SendMessageTest("18061532353")
	//r,err:=conn.Do("EXISTS","ptc_18061532353")
	//fmt.Printf("%d",r.(int64))

}
