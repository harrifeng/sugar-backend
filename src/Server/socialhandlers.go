package server

import (
	"db"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"utils"
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

func getLatestMessageU2uList(userId int, targetUserId int, latestMessageId int, needNumber int) responseBody {
	messages, err := db.GetLatestMessageToUser(userId, targetUserId, latestMessageId, needNumber)
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

func getHistoryMessageU2uList(userId int, targetUserId int, oldestMessageId int, needNumber int) responseBody {
	messages, err := db.GetHistoryMessageToUser(userId, targetUserId, oldestMessageId, needNumber)
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
			"host":        userId == int(group.UserID),
		}
	}
	return responseOKWithData(gin.H{
		"total": count,
		"data":  respGroups,
	})
}

func createGroup(userId int, groupName string, groupMembers string) responseBody {
	if groupName == "" {
		return responseNormalError("群组名不能为空")
	}
	var MemberList []string
	err := json.Unmarshal([]byte(groupMembers), &MemberList)
	if err != nil {
		return responseInternalServerError(err)
	}
	iMembers := make([]int, len(MemberList))
	for i, v := range MemberList {
		iMembers[i], err = strconv.Atoi(v)
	}
	err = db.AddGroup(userId, groupName, iMembers)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func sendMessageInGroup(userId int, groupId int, content string) responseBody {
	messageId, err := db.AddMessageInGroup(userId, groupId, content)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOKWithData(gin.H{
		"messageId": messageId,
	})
}

func getHistoryMessageInGroupList(userId int, groupId int, oldestMessageId int, needNumber int) responseBody {
	messages, err := db.GetHistoryMessageInGroup(groupId, oldestMessageId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respMessages := make([]gin.H, len(messages))
	for i, message := range messages {
		respMessages[i] = gin.H{
			"content":        message.Content,
			"messageId":      message.ID,
			"senderId":       message.SenderID,
			"senderUserName": message.Sender.UserName,
			"host":           message.SenderID == uint(userId),
			"createdAt":      message.CreatedAt,
			"imageUrl":       message.Sender.HeadPortraitUrl,
		}
	}
	return responseOKWithData(respMessages)
}

func getLatestMessageInGroupList(userId int, groupId int, latestMessageId int, needNumber int) responseBody {
	messages, err := db.GetLatestMessageInGroup(groupId, latestMessageId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respMessages := make([]gin.H, len(messages))
	for i, message := range messages {
		respMessages[i] = gin.H{
			"content":   message.Content,
			"messageId": message.ID,
			"host":      message.SenderID == uint(userId),
			"createdAt": message.CreatedAt,
			"imageUrl":  message.Sender.HeadPortraitUrl,
		}
	}
	return responseOKWithData(respMessages)
}

func getMessageList(userId int, existList string, needNumber int) responseBody {
	existIds := make(map[string][]int)
	err := json.Unmarshal([]byte(existList), &existIds)
	if err != nil {
		return responseInternalServerError(err)
	}
	groupNeedNumber := needNumber / 2
	u2uNeedNumber := needNumber - groupNeedNumber
	groupMessages, err := db.GetMessageInGroup(userId, existIds["groupIds"], groupNeedNumber)
	if len(groupMessages) < groupNeedNumber {
		u2uNeedNumber = needNumber - len(groupMessages)
	}
	u2uMessages, err := db.GetMessageU2u(userId, existIds["u2uIds"], u2uNeedNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	// build response data
	respGroupMessages := make([]gin.H, len(groupMessages))
	for i, message := range groupMessages {
		respGroupMessages[i] = gin.H{
			"groupId":        message.GroupID,
			"content":        message.Content,
			"groupName":      message.Group.Name,
			"senderUserName": message.Sender.UserName,
			"updatedTime":    utils.GoTimeToESTime(message.CreatedAt),
		}
	}
	respU2uMessages := make([]gin.H, len(u2uMessages))
	for i, message := range u2uMessages {
		respU2uMessages[i] = gin.H{
			"otherId":       message.OtherId,
			"content":       message.Content,
			"otherImageUrl": message.Other.HeadPortraitUrl,
			"otherUserName": message.Other.UserName,
			"updatedTime":   utils.GoTimeToESTime(message.CreatedAt),
		}
	}
	return responseOKWithData(gin.H{
		"groupMessages": respGroupMessages,
		"u2uMessages":   respU2uMessages,
	})
}

func getUserListInGroup(userId int, groupId int) responseBody {
	members, err := db.GetUserListInGroup(userId, groupId)
	if err != nil {
		return responseInternalServerError(err)
	}
	respMembers := make([]gin.H, len(members))
	for i, member := range members {
		respMembers[i] = gin.H{
			"userId":       member.ID,
			"userName":     member.UserName,
			"userImageUrl": member.HeadPortraitUrl,
		}
	}
	creatorId, err := db.GetHostInGroup(groupId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOKWithData(gin.H{
		"total": len(respMembers),
		"data":  respMembers,
		"host":  creatorId == uint(userId),
	})
}

func removeMemberInGroup(groupId int, memberId int) responseBody {
	err := db.RemoveMemberInGroup(groupId, memberId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func removeGroup(userId int, groupId int) responseBody {
	err := db.ReomveGroup(userId, groupId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func recommendUser(userId int) responseBody {
	users, err := db.GetRecommendUserList(userId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOKWithData(users)
}
