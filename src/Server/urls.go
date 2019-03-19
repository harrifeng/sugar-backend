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
	r.GET("/", urlsRoot)

	// account start
	r.GET("/accounts/code", accountSendVerificationCode)
	r.POST("/accounts/create", accountRegisterNewAccount)
	// account end

}
