package server

import (
	"github.com/gin-gonic/gin"
	"strconv"
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

func accountAlterInformation(c *gin.Context) {
	SessionId := c.PostForm("session_id")
	UserName := c.PostForm("username")
	Gender := c.PostForm("gender")
	Height, _ := strconv.ParseFloat(c.PostForm("height"), 64)
	Weight, _ := strconv.ParseFloat(c.PostForm("weight"), 64)
	Area := c.PostForm("area")
	Job := c.PostForm("job")
	Age, _ := strconv.Atoi(c.PostForm("age"))
	resp := alterUserInformation(SessionId, UserName, Gender, Height, Weight, Area, Job, Age)
	c.JSON(resp.Status, resp.Data)
}
