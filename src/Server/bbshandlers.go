package server

import (
	"db"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"utils"
)

func publishTopic(SessionId string, Content string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	err = db.AddTopic(userId, Content)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func publishTopicLordReply(SessionId string, TopicId string, Content string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	err = db.AddTopicLordReply(userId, TopicId, Content)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func publishTopicLayerReply(SessionId string, TopicLordReplyId string, Content string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	err = db.AddTopicLayerReply(userId, TopicLordReplyId, Content)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func removeTopic(SessionId string, TopicId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	err = db.RemoveTopic(TopicId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func removeTopicLordReply(SessionId string, TopicLordReplyId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	err = db.RemoveTopicLordReply(TopicLordReplyId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func getLatestTopicList(SessionId string, TopicListString string, NeedNumber string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	var TopicList []string
	err = json.Unmarshal([]byte(TopicListString), &TopicList)
	if err != nil {
		return responseInternalServerError(err)
	}
	topics, err := db.GetLatestTopicList(TopicList, NeedNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respTopics := make([]gin.H, len(topics))
	for i, topic := range topics {
		count, err := db.GetTopicReplyCount(strconv.Itoa(int(topic.ID)))
		if err != nil {
			return responseInternalServerError(err)
		}
		respTopics[i] = gin.H{
			"topicId":    topic.ID,
			"userId":     topic.UserID,
			"username":   topic.User.UserName,
			"iconUrl":    topic.User.HeadPortraitUrl,
			"lastTime":   topic.UpdatedAt,
			"content":    utils.StringCut(topic.Content, 40),
			"replyCount": count,
		}
	}
	return responseOKWithData(respTopics)
}

func getTopic(SessionId string, TopicId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	topic, err := db.GetTopicFromTopicId(TopicId)
	if err != nil {
		return responseInternalServerError(err)
	}
	collected, err := db.CheckUserCollectedTopic(userId, TopicId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOKWithData(gin.H{
		"userId":    topic.UserID,
		"username":  topic.User.UserName,
		"iconUrl":   topic.User.HeadPortraitUrl,
		"topicTime": topic.CreatedAt,
		"favorite":  collected,
		"likes":     topic.ThumbsUpCount,
		"content":   topic.Content,
	})
}

func getTopicLordReplyList(SessionId string, TopicId string, BeginFloor string, NeedNumber string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	topicLordReplies, err := db.GetTopicLordReplyListFromTopicId(TopicId, BeginFloor, NeedNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respReplies := make([]gin.H, len(topicLordReplies))
	beginFloor, _ := strconv.Atoi(BeginFloor)
	for i, reply := range topicLordReplies {
		count, _ := db.GetTopicLayerReplyCountFromTopicLordReply(strconv.Itoa(int(reply.ID)))
		respReplies[i] = gin.H{
			"replyId":   reply.ID,
			"floor":     beginFloor + i + 1,
			"userId":    reply.User.ID,
			"username":  reply.User.UserName,
			"iconUrl":   reply.User.HeadPortraitUrl,
			"replyTime": reply.CreatedAt,
			"likes":     reply.ThumbsUpCount,
			"content":   reply.Content,
			"comNumber": count,
		}
	}
	return responseOKWithData(respReplies)
}

func getTopicLayerReplyList(SessionId string, TopicLordReplyId string, BeginFloor string, NeedNumber string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	topicLayerReplies, err := db.GetTopicLayerReplyList(TopicLordReplyId, BeginFloor, NeedNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respReplies := make([]gin.H, len(topicLayerReplies))
	for i, reply := range topicLayerReplies {
		respReplies[i] = gin.H{
			"subreplyId":   reply.ID,
			"userId":       reply.UserID,
			"iconUrl":      reply.User.HeadPortraitUrl,
			"username":     reply.User.UserName,
			"content":      reply.Content,
			"subreplyTime": reply.CreatedAt,
			"likes":        reply.ThumbsUpCount,
		}
	}
	return responseOKWithData(respReplies)
}

func valueTopic(SessionId string, TopicId string, Value string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	err = db.ValueTopic(TopicId, Value)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func valueTopicLordReply(SessionId string, TopicLordReplyId string, Value string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	err = db.ValueTopicLordReply(TopicLordReplyId, Value)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func valueTopicLayerReply(SessionId string, TopicLayerReplyId string, Value string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	err = db.ValueTopicLayerReply(TopicLayerReplyId, Value)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func addCollectedTopic(SessionId string, TopicId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	err = db.AddUserCollectedTopic(userId, TopicId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func removeCollectedTopic(SessionId string, TopicId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	err = db.RemoveUserCollectedTopic(userId, TopicId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}
