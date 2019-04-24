package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func urlsRoot(c *gin.Context) {
	c.String(http.StatusOK, "loading...")
}

func initRouter() {
	r.GET("/", urlsRoot)

	accountsGroup := r.Group("/accounts")
	{
		accountsGroup.GET("/code", accountSendVerificationCode)
		accountsGroup.POST("/register", accountRegister)
		accountsGroup.GET("/login", accountLogin)
		accountsGroup.POST("/alter/password", accountAlterPassword)
		accountsGroup.GET("/logout", accountLogout)
		accountsGroup.POST("/alter", LoginAuth(), accountAlterInformation)
		accountsGroup.GET("/info", LoginAuth(), accountGetUserInformation)
		accountsGroup.GET("/privacy", LoginAuth(), accountGetUserPrivacySetting)
		accountsGroup.POST("/alter/privacy", LoginAuth(), accountAlterUserPrivacySetting)
		accountsGroup.GET("/follower", LoginAuth(), accountGetUserFollowerList)
		accountsGroup.GET("/following", LoginAuth(), accountGetUserFollowingList)
		accountsGroup.POST("/following/follow", LoginAuth(), accountFollowUser)
		accountsGroup.POST("/following/ignore", LoginAuth(), accountIgnoreUser)

	}
	bbsGroup := r.Group("bbs")
	bbsGroup.Use(LoginAuth())
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
	schoolGroup.Use(LoginAuth())
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
	homeGroup := r.Group("home")
	homeGroup.Use(LoginAuth())
	{
		homeGroup.POST("/checkin", homeCheckInUser)
		homeGroup.GET("/families", homeGetUserFamilyMemberList)
		homeGroup.POST("/family/link", homeLinkNewFamilyMember)
		homeGroup.POST("/blood-sugar/record", homeRecordBloodSugar)
		homeGroup.GET("/blood-sugar/records", homeGetBloodSugarRecordList)
		homeGroup.GET("/blood-sugar/record", homeGetBloodSugarRecord)
		homeGroup.POST("/health/record", homeRecordHealth)
		homeGroup.GET("/health/records", homeGetHealthRecordList)
		homeGroup.POST("/health/voice", homeParseHealthRecordVoiceInput)
		homeGroup.POST("/blood-sugar/voice", homeParseBloodSugarRecordVoiceInput)
		homeGroup.GET("/sugar-guide", homeSugarGuideWebsocket)
		homeGroup.GET("/weekly-newspaper",homeWeeklyNewspaper)
	}
	socialGroup := r.Group("social")
	socialGroup.Use(LoginAuth())
	{
		socialGroup.GET("/group/members", socialGetUserInGroup)
		socialGroup.POST("/group/create", socialCreateGroup)
		socialGroup.POST("/group/remove", socialRemoveGroup)
		socialGroup.POST("/group/level", socialMemberLevelGroup)
		socialGroup.GET("/messages", socialGetMessageList)
		socialGroup.GET("/groups", socialGetUserJoinGroupList)
		socialGroup.POST("/group/chatting/send", socialSendMessageInGroup)
		socialGroup.GET("/group/chatting/records", socialGetHistoryMessageInGroupList)
		socialGroup.GET("/group/chatting/records/latest", socialGetLatestMessageInGroupList)
		socialGroup.GET("/chatting/records", socialGetHistoryMessageU2uList)
		socialGroup.GET("/chatting/records/latest", socialGetLatestMessageU2uList)
		socialGroup.POST("/chatting/send", socialSendMessageUser)
		socialGroup.GET("/recommend", socialRecommendUser)
	}
}
