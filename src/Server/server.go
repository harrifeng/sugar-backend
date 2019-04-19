package server

import (
	"db"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

var r *gin.Engine

func Start(port uint) {
	// init database
	db.AutoCreateTableTest()
	//init gin http server
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static")
	initRouter()
	log.Fatal(r.Run(fmt.Sprintf(":%d",port)))
}
