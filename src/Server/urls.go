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

	accountsGroup := r.Group("/accounts")
	{
		accountsGroup.GET("/code", accountSendVerificationCode)
		accountsGroup.POST("/register", accountRegister)
		accountsGroup.GET("/login", accountLogin)
		accountsGroup.POST("/alter", accountAlterInformation)
		accountsGroup.GET("/info", accountGetUserInformation)
		accountsGroup.POST("/alter/password", accountAlterPassword)
		accountsGroup.GET("/privacy", accountGetUserPrivacySetting)
		accountsGroup.POST("/alter/privacy", accountAlterUserPrivacySetting)
		accountsGroup.GET("/logout", accountLogout)
		accountsGroup.GET("/follower", accountGetUserFollowerList)
		accountsGroup.GET("/following", accountGetUserFollowingList)
		accountsGroup.POST("/following/follow", accountFollowUser)
		accountsGroup.POST("/following/ignore", accountIgnoreUser)
	}
	bbsGroup:=r.Group("bbs")
	{
		bbsGroup.POST("/topic/publish", bbsPublishTopic)
		bbsGroup.POST("/topic/lord-reply/publish", bbsPublishTopicLordReply)
		bbsGroup.POST("/topic/layer-reply/publish", bbsPublishTopicLayerReply)
		bbsGroup.POST("/topic/remove", bbsRemoveTopic)
		bbsGroup.POST("/topic/lord-reply/remove", bbsRemoveTopicLordReply)
		bbsGroup.GET("/topics", bbsGetLatestTopicList)
		bbsGroup.GET("/topic", bbsGetTopic)
		bbsGroup.GET("/topic/lord-reply", bbsGetTopicLordReplyList)
		bbsGroup.GET("/topic/layer-reply", bbsGetTopicLayerReplyList)
		bbsGroup.POST("/topic/value", bbsValueTopic)
		bbsGroup.POST("/topic/lord-reply/value", bbsValueTopicLordReply)
		bbsGroup.POST("/topic/layer-reply/value", bbsValueTopicLayerReply)
		bbsGroup.GET("/topics/search", bbsSearchTopic)
		bbsGroup.POST("/topic/user-collect", bbsCollectTopic)
		bbsGroup.POST("/topic/user-cancel-collect", bbsCancelCollectTopic)
		bbsGroup.GET("/topics/user-collected", bbsGetUserCollectedTopicList)
		bbsGroup.GET("/topics/user-published", bbsGetUserPublishedTopicList)
		bbsGroup.GET("/topics/user-replies", bbsGetUserTopicReplyList)
	}
	schoolGroup := r.Group("school")
	{
		schoolGroup.GET("/article-page/:id", schoolGetArticlePage)
		schoolGroup.GET("/article", schoolGetArticle)
		schoolGroup.GET("/articles", schoolGetArticleList)
		schoolGroup.POST("/article/comment", schoolPublishArticleComment)
		schoolGroup.GET("/article/comments", schoolGetArticleCommentList)
		schoolGroup.GET("/articles/search", schoolSearchArticle)
		schoolGroup.GET("/articles/user-collected", schoolGetUserCollectedArticleList)
		schoolGroup.GET("/articles/user-comments", schoolGetUserArticleCommentList)
		schoolGroup.POST("/articles/user-collect", schoolCollectArticle)
		schoolGroup.POST("/articles/user-cancel-collect", schoolCancelCollectArticle)
		schoolGroup.POST("/article/comments/value", schoolValueArticle)
	}
}
