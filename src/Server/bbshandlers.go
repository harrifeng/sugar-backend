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
