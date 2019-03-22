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

func accountLogout(c *gin.Context) {
	SessionId := c.Query("session_id")
	resp := logoutUser(SessionId)
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

func accountGetUserInformation(c *gin.Context) {
	SessionId := c.Query("session_id")
	OtherUserId := c.Query("other_user_id")
	var resp responseBody
	if OtherUserId == "" {
		resp = getUserInformationFromSessionId(SessionId)
	} else {
		resp = getUserInformationFromUserId(SessionId, OtherUserId)
	}
	c.JSON(resp.Status, resp.Data)
}

func accountAlterPassword(c *gin.Context) {
	PhoneNumber := c.PostForm("phone_number")
	Code := c.PostForm("code")
	NewPassword := c.PostForm("password")
	resp := alterPassword(PhoneNumber, Code, NewPassword)
	c.JSON(resp.Status, resp.Data)
}

func accountGetUserPrivacySetting(c *gin.Context) {
	SessionId := c.Query("session_id")
	resp := getUserPrivacySetting(SessionId)
	c.JSON(resp.Status, resp.Data)
}

func accountAlterUserPrivacySetting(c *gin.Context) {
	SessionId := c.PostForm("session_id")
	ShowPhoneNumber := c.PostForm("show_phone_number") == "1"
	ShowGender := c.PostForm("show_gender") == "1"
	ShowAge := c.PostForm("show_age") == "1"
	ShowHeight := c.PostForm("show_height") == "1"
	ShowWeight := c.PostForm("show_weight") == "1"
	ShowArea := c.PostForm("show_area") == "1"
	ShowJob := c.PostForm("show_job") == "1"
	resp := alterUserPrivacy(SessionId, ShowPhoneNumber, ShowGender,
		ShowAge, ShowHeight, ShowWeight, ShowArea, ShowJob)
	c.JSON(resp.Status, resp.Data)
}

func accountGetUserFollowingList(c *gin.Context) {
	SessionId := c.Query("session_id")
	BeginId := c.Query("begin_id")
	NeedNumber := c.Query("need_number")
	resp := getUserFollowingList(SessionId, BeginId, NeedNumber)
	c.JSON(resp.Status, resp.Data)
}

func accountGetUserFollowerList(c *gin.Context) {
	SessionId := c.Query("session_id")
	BeginId := c.Query("begin_id")
	NeedNumber := c.Query("need_number")
	resp := getUserFollowerList(SessionId, BeginId, NeedNumber)
	c.JSON(resp.Status, resp.Data)
}

func accountFollowUser(c *gin.Context) {
	SessionId := c.PostForm("session_id")
	UserId := c.PostForm("other_user_id")
	resp := followUser(SessionId, UserId)
	c.JSON(resp.Status, resp.Data)
}

func accountIgnoreUser(c *gin.Context) {
	SessionId := c.PostForm("session_id")
	UserId := c.PostForm("other_user_id")
	resp := ignoreUser(SessionId, UserId)
	c.JSON(resp.Status, resp.Data)
}

func bbsPublishTopic(c *gin.Context) {

}
