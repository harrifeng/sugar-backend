package main

import (
	"db"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"server"
)

func main() {
	var genFlag, runFlag, helpFlag,releaseFlag,initFlag bool
	var port uint
	flag.BoolVar(&genFlag, "gen", false, "generate a configuration example file.")
	flag.BoolVar(&initFlag, "init", false, "init the database")
	flag.BoolVar(&runFlag, "server", false, "run server on debug")
	flag.BoolVar(&releaseFlag, "release", false, "run server on release")
	flag.BoolVar(&helpFlag, "help", false, "cat help information")
	flag.UintVar(&port,"port",8080,"set port of server")
	flag.Parse()
	if helpFlag {
		flag.Usage()
	}else if initFlag{
		db.Init()
	} else if genFlag && !runFlag {
		genConfigurationFile()
	} else if !genFlag && runFlag {
		if releaseFlag{
			gin.SetMode(gin.ReleaseMode)
		}
		runServer(port)
	} else {
		fmt.Println("only can do one thing")
	}
}

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

// run server
func runServer(port uint) {
	// init configuration
	config, err := server.LoadConfiguration()
	if err != nil {
		fmt.Println(err)
		return
	}
	db.InitConfiguration(&config.MysqlConfig, &config.RedisConfig)

	// init mysql
	mysqlDb, err := db.InitMysql()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err = mysqlDb.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	// init redis
	redisPool := db.InitRedis()
	defer func() {
		err := redisPool.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	server.Start(port)
}
