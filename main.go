package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"server"
)

func main() {
	var genFlag, runFlag, helpFlag, releaseFlag, initFlag bool
	var port uint
	flag.BoolVar(&genFlag, "gen", false, "generate a configuration example file.")
	flag.BoolVar(&initFlag, "init", false, "init the database")
	flag.BoolVar(&runFlag, "server", false, "run server on debug(default)")
	flag.BoolVar(&releaseFlag, "release", false, "run server on release")
	flag.BoolVar(&helpFlag, "help", false, "cat help information")
	flag.UintVar(&port, "port", 8080, "set port of server")
	flag.Parse()
	if helpFlag {
		flag.Usage()
	} else if initFlag {
		initDatabase()
	} else if genFlag && !runFlag {
		genConfigurationFile()
	} else if !genFlag && runFlag {
		if releaseFlag {
			gin.SetMode(gin.ReleaseMode)
		}
		runServer(port)
	} else {
		fmt.Println("flags error")
	}
}
// generate configuration file
func genConfigurationFile() {
	exeFile, err := exec.LookPath(os.Args[0])
	if err != nil {
		fmt.Println(err)
	}
	dir, _ := filepath.Split(exeFile)
	filePath := filepath.Join(dir, "config.example.json")
	fileJson, err := server.GenerateJSON()
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(filePath, fileJson, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
}

// init database
func initDatabase(){
	dbClose,err:=server.ConnectDatabase()
	if err!=nil{
		log.Fatal(err)
		return
	}
	defer dbClose()
	server.Init()
}

// run server
func runServer(port uint) {
	dbClose,err:=server.ConnectDatabase()
	if err!=nil{
		log.Fatal(err)
		return
	}
	defer dbClose()
	server.Start(port)
}
