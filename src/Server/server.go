package server

import (
	"db"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

var r *gin.Engine

func ConnectDatabase() (func(), error) {
	config, err := loadConfiguration()
	if err != nil {
		return func() {}, err
	}
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
