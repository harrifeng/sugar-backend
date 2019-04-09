package server

import (
	"db"
	"github.com/gin-gonic/gin"
)

func sendMessageToUser(userId int, content string, targetUserId int) responseBody {
	messageId, err := db.AddMessageToUser(userId, content, targetUserId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOKWithData(gin.H{
		"messageId": messageId,
	})
}

func getLatestMessageU2uList(userId int, targetUserId int, latestMessageId int, NeedNumber int) responseBody {
	messages, err := db.GetLatestMessageToUser(userId, targetUserId, latestMessageId, NeedNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respMessages := make([]gin.H, len(messages))
	for i, message := range messages {
		var imageUrl string
		if message.SenderID == uint(userId) {
			imageUrl = message.Sender.HeadPortraitUrl
		} else {
			imageUrl = message.Target.HeadPortraitUrl
		}
		respMessages[i] = gin.H{
			"content":   message.Content,
			"messageId": message.ID,
			"host":      message.SenderID == uint(userId),
			"createdAt": message.CreatedAt,
			"imageUrl":  imageUrl,
		}
	}
	return responseOKWithData(respMessages)
}

func getHistoryMessageU2uList(userId int, targetUserId int, oldestMessageId int, NeedNumber int) responseBody {
	messages, err := db.GetHistoryMessageToUser(userId, targetUserId, oldestMessageId, NeedNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respMessages := make([]gin.H, len(messages))
	for i, message := range messages {
		var imageUrl string
		if message.SenderID == uint(userId) {
			imageUrl = message.Sender.HeadPortraitUrl
		} else {
			imageUrl = message.Target.HeadPortraitUrl
		}
		respMessages[i] = gin.H{
			"content":   message.Content,
			"messageId": message.ID,
			"host":      message.SenderID == uint(userId),
			"createdAt": message.CreatedAt,
			"imageUrl":  imageUrl,
		}
	}
	return responseOKWithData(respMessages)
}

func getUserJoinGroupList(userId int, beginId int, needNumber int) responseBody {
	groups, count, err := db.GetUserJoinGroupList(userId, beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respGroups := make([]gin.H, len(groups))
	for i, group := range groups {
		respGroups[i] = gin.H{
			"groupId":     group.ID,
			"name":        group.Name,
			"memberCount": len(group.Members),
		}
	}
	return responseOKWithData(gin.H{
		"total": count,
		"data":  respGroups,
	})
}
