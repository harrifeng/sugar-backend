package main

import (
	"Db"
	"Server"
	"fmt"
)

func DbTest() {
	Db.AutoCreateTableTest()
}

func main() {
	// init db
	db, err := Db.Init()
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
	Server.Start()

}
