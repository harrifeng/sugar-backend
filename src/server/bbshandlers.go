package server

import (
	"db"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"utils"
)

func publishTopic(userId int, content string) responseBody {
	err := db.AddTopic(userId, content)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func publishTopicLordReply(userId int, topicId int, content string) responseBody {
	err := db.AddTopicLordReply(userId, topicId, content)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func publishTopicLayerReply(userId int, topicLordReplyId int, content string) responseBody {
	err := db.AddTopicLayerReply(userId, topicLordReplyId, content)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func removeTopic(topicId int) responseBody {
	err := db.RemoveTopic(topicId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func removeTopicLordReply(topicLordReplyId int) responseBody {
	err := db.RemoveTopicLordReply(topicLordReplyId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

// wait to trans[]
func getLatestTopicList(topicListString string, needNumber int) responseBody {
	var TopicList []string
	err := json.Unmarshal([]byte(topicListString), &TopicList)
	if err != nil {
		return responseInternalServerError(err)
	}
	topics, err := db.GetLatestTopicList(TopicList, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respTopics := make([]gin.H, len(topics))
	for i, topic := range topics {
		count, err := db.GetTopicReplyCount(int(topic.ID))
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

func getTopic(userId int, topicId int) responseBody {
	topic, err := db.GetTopicFromTopicId(topicId)
	if err != nil {
		return responseInternalServerError(err)
	}
	collected, err := db.CheckUserCollectedTopic(userId, topicId)
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

func getTopicLordReplyList(topicId int, beginFloor int, needNumber int) responseBody {
	topicLordReplies, err := db.GetTopicLordReplyListFromTopicId(topicId, beginFloor, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respReplies := make([]gin.H, len(topicLordReplies))
	for i, reply := range topicLordReplies {
		count, _ := db.GetTopicLayerReplyCountFromTopicLordReply(int(reply.ID))
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

func getTopicLayerReplyList(topicLordReplyId int, beginFloor int, needNumber int) responseBody {
	topicLayerReplies, err := db.GetTopicLayerReplyListFromTopicLordReplyId(topicLordReplyId, beginFloor, needNumber)
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

func valueTopic(topicId int, value int) responseBody {
	err := db.ValueTopic(topicId, value)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func valueTopicLordReply(topicLordReplyId int, value int) responseBody {
	err := db.ValueTopicLordReply(topicLordReplyId, value)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func valueTopicLayerReply(topicLayerReplyId int, value int) responseBody {
	err := db.ValueTopicLayerReply(topicLayerReplyId, value)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func addCollectedTopic(userId int, topicId int) responseBody {
	err := db.AddUserCollectedTopic(userId, topicId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func removeCollectedTopic(userId int, TopicId int) responseBody {
	err := db.RemoveUserCollectedTopic(userId, TopicId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func getSearchTopicList(searchContent string, beginId int, needNumber int) responseBody {
	topics, err := db.GetSearchTopicList(searchContent, beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respTopics := make([]gin.H, len(topics))
	for i, topic := range topics {
		respTopics[i] = gin.H{
			"topicId":  topic.ID,
			"userId":   topic.UserID,
			"username": topic.User.UserName,
			"iconUrl":  topic.User.HeadPortraitUrl,
			"lastTime": topic.UpdatedAt,
			"content":  utils.StringCut(topic.Content, 40),
		}
	}
	return responseOKWithData(respTopics)
}

func getUserCollectedTopicList(userId int, beginId int, needNumber int) responseBody {
	topics, count, err := db.GetUserCollectedTopicList(userId, beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respTopics := make([]gin.H, len(topics))
	for i, topic := range topics {
		lastFloor, err := db.GetLastFloorFromTopicId(strconv.Itoa(int(topic.ID)))
		if err != nil {
			return responseInternalServerError(err)
		}
		respTopics[i] = gin.H{
			"topicId":   topic.ID,
			"userId":    topic.UserID,
			"username":  topic.User.UserName,
			"iconUrl":   topic.User.HeadPortraitUrl,
			"topicTime": topic.CreatedAt,
			"content":   utils.StringCut(topic.Content, 40),
			"lastFloor": lastFloor,
		}
	}
	return responseOKWithData(gin.H{
		"data":  respTopics,
		"total": count,
	})
}

func getUserPublishedTopicList(userId int, beginId int, needNumber int) responseBody {
	topics, count, err := db.GetUserPublishedTopicList(userId, beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respTopics := make([]gin.H, len(topics))
	for i, topic := range topics {
		replyCount, err := db.GetTopicReplyCount(int(topic.ID))
		if err != nil {
			return responseInternalServerError(err)
		}
		collectingCount, err := db.GetTopicCollectingUserCount(strconv.Itoa(int(topic.ID)))
		respTopics[i] = gin.H{
			"topicId":     topic.ID,
			"userId":      topic.UserID,
			"username":    topic.User.UserName,
			"iconUrl":     topic.User.HeadPortraitUrl,
			"topicTime":   topic.CreatedAt,
			"content":     utils.StringCut(topic.Content, 40),
			"likes":       topic.ThumbsUpCount,
			"replyCount":  replyCount,
			"favoriteNum": collectingCount,
		}
	}
	return responseOKWithData(gin.H{
		"data":  respTopics,
		"total": count,
	})
}

func getUserReplyList(userId int, beginId int, needNumber int) responseBody {
	replies, count, err := db.GetUserReplyList(userId, beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respReplies := make([]gin.H, len(replies))
	for i, reply := range replies {
		if reply.TopicLordReplyKey == 0 && reply.TopicLayerReplyKey != 0 {
			topicLordReply, err := db.GetTopicLordReplyFromTopicLordReplyId(int(reply.FatherID))
			if err != nil {
				return responseInternalServerError(err)
			}
			floor, err := db.GetTopicLordReplyFloor(int(reply.FatherID))
			if err != nil {
				return responseInternalServerError(err)
			}
			respReplies[i] = gin.H{
				"type":            "subreply",
				"subreplyId":      reply.TopicLayerReplyKey,
				"replyId":         reply.FatherID,
				"topicId":         topicLordReply.TopicID,
				"time":            reply.CreatedAt,
				"subreplyContent": reply.Content,
				"replyContent":    topicLordReply.Content,
				"likes":           reply.ThumbsUpCount,
				"floor":           floor,
			}
		} else if reply.TopicLordReplyKey != 0 && reply.TopicLayerReplyKey == 0 {
			topic, err := db.GetTopicFromTopicId(int(reply.FatherID))
			if err != nil {
				return responseInternalServerError(err)
			}
			respReplies[i] = gin.H{
				"type":         "reply",
				"replyId":      reply.TopicLordReplyKey,
				"topicId":      reply.FatherID,
				"time":         reply.CreatedAt,
				"replyContent": reply.Content,
				"topicContent": topic.Content,
				"likes":        reply.ThumbsUpCount,
			}
		} else {
			return responseInternalServerError(errors.New("invalid reply"))
		}
	}
	return responseOKWithData(gin.H{
		"data":  respReplies,
		"total": count,
	})
}
