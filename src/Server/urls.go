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
		r.GET("/test/db/create", createNewTestUser)
		//test url end
	}

	r.GET("/", urlsRoot)

	// account start
	r.GET("/accounts/code", accountSendVerificationCode)
	r.POST("/accounts/create", accountRegister)
	r.GET("/accounts/login", accountLogin)
	r.POST("/accounts/alter", accountAlterInformation)
	r.GET("/accounts/info", accountGetUserInformation)
	r.POST("/accounts/alter/password", accountAlterPassword)
	r.GET("/accounts/privacy", accountGetUserPrivacySetting)
	r.POST("/accounts/alter/privacy", accountAlterUserPrivacySetting)
	r.GET("accounts/logout", accountLogout)
	// account end

}
