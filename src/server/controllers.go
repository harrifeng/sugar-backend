package server

import "github.com/gin-gonic/gin"

func index(c *gin.Context) {
	c.String(200, "pong")
}
