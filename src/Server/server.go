package Server

import (
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Start() {
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	initRouter()
	_ = r.Run(":8080")
}
