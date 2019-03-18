package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func urlsRoot(c *gin.Context){
	urls := map[string]string{
		"GET_code_url": "/code",
	}
	j,_ :=json.Marshal(urls)
	r :=string(j)
	c.String(http.StatusOK,r)
}

func initRouter() {
	r.GET("/", urlsRoot)

	// account start
	accountGroup := r.Group("accounts")
	{
		r.GET("code",accountSendVerificationCode)
		r.POST("create",accountRegisterNewAccount)
	}
	// account end

}
