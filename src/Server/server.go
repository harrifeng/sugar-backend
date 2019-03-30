package server

import (
	"db"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Start() {
	db.AutoCreateTableTest()
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static")
	initRouter()
	_ = r.Run(":8080")
}
