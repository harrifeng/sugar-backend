package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func accountSendVerificationCode(c *gin.Context) {
	PhoneNumber := c.Query("phone_number")
	err := SendVerificationCode(PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.String(http.StatusOK, "ok")
}

func accountRegisterNewAccount(c *gin.Context) {
	//PhoneNumber :=c.PostForm("phone_number")
	//UserName := c.PostForm("username")
	//Password := c.PostForm("password")
	//Code :=c.PostForm("code")

}
