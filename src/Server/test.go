package server

import (
	"db"
	"github.com/gin-gonic/gin"
)

func initDb(c *gin.Context) {
	db.AutoCreateTableTest()
}
