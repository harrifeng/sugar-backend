package server

import (
	"db"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Start() {
	// init database
	db.AutoCreateTableTest()
	//init gin http server
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static")
	initRouter()
	_ = r.Run(":8080")
}
