package server

import (
	"github.com/gin-gonic/gin"
)

func accountSendVerificationCode(c *gin.Context) {
	PhoneNumber := c.Query("phone_number")
	resp := sendVerificationCode(PhoneNumber)
	c.JSON(resp.Status, resp.Data)
}

func accountRegister(c *gin.Context) {
	PhoneNumber := c.PostForm("phone_number")
	Password := c.PostForm("password")
	Code := c.PostForm("code")
	resp := registerNewUser(PhoneNumber, Password, Code)
	c.JSON(resp.Status, resp.Data)
}

func accountLogin(c *gin.Context) {
	PhoneNumber := c.Query("phone_number")
	Password := c.Query("password")
	resp := loginUser(PhoneNumber, Password)
	c.JSON(resp.Status, resp.Data)
}
