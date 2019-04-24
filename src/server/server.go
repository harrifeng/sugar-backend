package server

import (
	"db"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var r *gin.Engine
var serverConfig *MainConfiguration

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnectDatabase() (func(), error) {
	config, err := loadConfiguration()
	if err != nil {
		return func() {}, err
	}
	serverConfig = &config.ServerConfig
	db.InitConfiguration(&config.MysqlConfig, &config.RedisConfig)

	// init mysql
	mysqlDb, err := db.InitMysql()
	if err != nil {
		fmt.Println(err)
		return func() {}, err
	}
	// init redis
	redisPool := db.InitRedis()

	// return defer function
	return func() {
		_ = mysqlDb.Close()
		_ = redisPool.Close()
	}, err
}

func Start(port uint) {
	//init gin http-server
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static")
	initRouter()
	log.Fatal(r.Run(fmt.Sprintf(":%d", port)))
}
