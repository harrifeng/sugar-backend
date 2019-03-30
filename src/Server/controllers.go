package server

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
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
	UserName := c.PostForm("username")
	Code := c.PostForm("code")
	resp := registerNewUser(PhoneNumber, Password, UserName, Code)
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

func schoolGetArticle(c *gin.Context) {
	SessionId := c.Query("session_id")
	ArticleId := c.Query("article_id")
	resp := getArticle(SessionId, ArticleId)
	c.JSON(resp.Status, resp.Data)
}

func schoolGetArticlePage(c *gin.Context) {
	ArticleId := c.Param("id")
	content, err := getArticlePage(ArticleId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.HTML(http.StatusOK, "article.tmpl", gin.H{
		"content": template.HTML(content),
	})
}

func schoolGetArticleList(c *gin.Context) {
	SessionId := c.Query("session_id")
	BeginId := c.Query("begin_id")
	NeedNumber := c.Query("need_number")
	resp := getArticleList(SessionId, BeginId, NeedNumber)
	c.JSON(resp.Status, resp.Data)
}

func schoolPublishArticleComment(c *gin.Context) {
	SessionId := c.PostForm("session_id")
	ArticleId := c.PostForm("id")
	Content := c.PostForm("content")
	resp := createArticleComment(SessionId, ArticleId, Content)
	c.JSON(resp.Status, resp.Data)
}

func schoolGetArticleCommentList(c *gin.Context) {
	SessionId := c.Query("session_id")
	ArticleId := c.Query("id")
	BeginId := c.Query("begin_id")
	NeedNumber := c.Query("need_number")
	resp := getArticleCommentList(SessionId, ArticleId, BeginId, NeedNumber)
	c.JSON(resp.Status, resp.Data)
}

func schoolSearchArticle(c *gin.Context) {
	SessionId := c.Query("session_id")
	SearchContent := c.Query("content")
	BeginId := c.Query("begin_id")
	NeedNumber := c.Query("need_number")
	resp := getSearchArticleList(SessionId, SearchContent, BeginId, NeedNumber)
	c.JSON(resp.Status, resp.Data)
}

func schoolGetUserCollectedArticleList(c *gin.Context) {
	SessionId := c.Query("session_id")
	BeginId := c.Query("begin_id")
	NeedNumber := c.Query("need_number")
	resp := getUserCollectedArticleList(SessionId, BeginId, NeedNumber)
	c.JSON(resp.Status, resp.Data)
}

func schoolGetUserArticleCommentList(c *gin.Context) {
	SessionId := c.Query("session_id")
	BeginId := c.Query("begin_id")
	NeedNumber := c.Query("need_number")
	resp := getUserArticleCommentList(SessionId, BeginId, NeedNumber)
	c.JSON(resp.Status, resp.Data)
}

func schoolCollectArticle(c *gin.Context) {
	SessionId := c.PostForm("session_id")
	ArticleId := c.PostForm("article_id")
	resp := addCollectedArticle(SessionId, ArticleId)
	c.JSON(resp.Status, resp.Data)
}

func schoolCancelCollectArticle(c *gin.Context) {
	SessionId := c.PostForm("session_id")
	ArticleId := c.PostForm("article_id")
	resp := removeCollectedArticle(SessionId, ArticleId)
	c.JSON(resp.Status, resp.Data)
}

func bbsPublishTopic(c *gin.Context) {
	SessionId := c.PostForm("session_id")
	Content := c.PostForm("content")
	resp := publishTopic(SessionId, Content)
	c.JSON(resp.Status, resp.Data)
}

func bbsPublishTopicLordReply(c *gin.Context) {
	SessionId := c.PostForm("session_id")
	TopicId := c.PostForm("topic_id")
	Content := c.PostForm("content")
	resp := publishTopicLordReply(SessionId, TopicId, Content)
	c.JSON(resp.Status, resp.Data)
}

func bbsPublishTopicLayerReply(c *gin.Context) {
	SessionId := c.PostForm("session_id")
	TopicLordReplyId := c.PostForm("topic_lord_reply_id")
	Content := c.PostForm("content")
	resp := publishTopicLayerReply(SessionId, TopicLordReplyId, Content)
	c.JSON(resp.Status, resp.Data)
}

func bbsRemoveTopic(c *gin.Context) {
	SessionId := c.PostForm("session_id")
	TopicId := c.PostForm("topic_id")
	resp := removeTopic(SessionId, TopicId)
	c.JSON(resp.Status, resp.Data)
}

func bbsRemoveTopicLordReply(c *gin.Context) {
	SessionId := c.PostForm("session_id")
	TopicLordReplyId := c.PostForm("topic_lord_reply_id")
	resp := removeTopicLordReply(SessionId, TopicLordReplyId)
	c.JSON(resp.Status, resp.Data)
}

func bbsGetLatestTopicList(c *gin.Context) {
	SessionId := c.PostForm("session_id")
	TopicList := c.PostForm("topic_id_list")
	NeedNumber := c.PostForm("need_number")
	resp := getLatestTopicList(SessionId, TopicList, NeedNumber)
	c.JSON(resp.Status, resp.Data)
}

func bbsGetTopic(c *gin.Context) {
	SessionId := c.Query("session_id")
	TopicId := c.Query("topic_id")
	resp := getTopic(SessionId, TopicId)
	c.JSON(resp.Status, resp.Data)
}

func bbsGetTopicLordReplyList(c *gin.Context) {
	SessionId := c.Query("session_id")
	TopicId := c.Query("topic_id")
	BeginFloor := c.Query("begin_floor")
	NeedNumber := c.Query("need_number")
	resp := getTopicLordReplyList(SessionId, TopicId, BeginFloor, NeedNumber)
	c.JSON(resp.Status, resp.Data)
}

func bbsGetTopicLayerReplyList(c *gin.Context) {
	SessionId := c.Query("session_id")
	TopicLordReplyId := c.Query("topic_lord_reply_id")
	BeginFloor := c.Query("begin_floor")
	NeedNumber := c.Query("need_number")
	resp := getTopicLayerReplyList(SessionId, TopicLordReplyId, BeginFloor, NeedNumber)
	c.JSON(resp.Status, resp.Data)
}

func bbsValueTopic(c *gin.Context) {
	SessionId := c.Query("session_id")
	TopicId := c.Query("topic_id")
	Value := c.Query("value")
	resp := valueTopic(SessionId, TopicId, Value)
	c.JSON(resp.Status, resp.Data)
}

func bbsValueTopicLordReply(c *gin.Context) {
	SessionId := c.Query("session_id")
	TopicLordReplyId := c.Query("topic_lord_reply_id")
	Value := c.Query("value")
	resp := valueTopicLordReply(SessionId, TopicLordReplyId, Value)
	c.JSON(resp.Status, resp.Data)
}

func bbsValueTopicLayerReply(c *gin.Context) {
	SessionId := c.Query("session_id")
	TopicLayerReplyId := c.Query("topic_layer_reply_id")
	Value := c.Query("value")
	resp := valueTopicLayerReply(SessionId, TopicLayerReplyId, Value)
	c.JSON(resp.Status, resp.Data)
}
