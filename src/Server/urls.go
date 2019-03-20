package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func urlsRoot(c *gin.Context) {
	urls := map[string]string{
		"GET_code_url": "/code",
	}
	j, _ := json.Marshal(urls)
	r := string(j)
	c.String(http.StatusOK, r)
}

func initRouter() {
	if gin.Mode() == "debug" {
		//test start
		r.GET("/test/db/init", initDb)
		//test url end
	}

	r.GET("/", urlsRoot)

	// account start
	r.GET("/accounts/code", accountSendVerificationCode)
	r.POST("/accounts/create", accountRegister)
	r.GET("/accounts/login", accountLogin)
	r.POST("/accounts/alter", accountAlterInformation)
	// account end

}
