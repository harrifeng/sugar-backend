package server

import (
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Start() {
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/*")
	initRouter()
	_ = r.Run(":8080")
}
