package db

func AddMessageToUser(userId int, content string, targetUserId int) (uint, error) {
	message := MessageU2u{
		SenderID: uint(userId),
		Content:  content,
		TargetID: uint(targetUserId),
	}
	err := mysqlDb.Create(&message).Save(&message).Error
	return message.ID, err
}

func GetLatestMessageToUser(userId int, targetUserId int, latestMessageId int, needNumber int) ([]MessageU2u, error) {
	var messages []MessageU2u
	err := mysqlDb.Preload("Sender").Preload("Target").
		Where("sender_id=? and target_id=? and id > ?", userId, targetUserId, latestMessageId).
		Or("sender_id=? and target_id=? and id > ?", targetUserId, userId, latestMessageId).
		Order("id desc").Limit(needNumber).Find(&messages).Error
	return messages, err
}

func GetHistoryMessageToUser(userId int, targetUserId int, oldestMessageId int, needNumber int) ([]MessageU2u, error) {
	var messages []MessageU2u
	var err error
	if oldestMessageId > 0 {
		err = mysqlDb.Preload("Sender").Preload("Target").
			Where("sender_id=? and target_id=? and id < ?", userId, targetUserId, oldestMessageId).
			Or("sender_id=? and target_id=? and id < ?", targetUserId, userId, oldestMessageId).
			Order("id desc").Limit(needNumber).Find(&messages).Error
	} else {
		err = mysqlDb.Preload("Sender").Preload("Target").
			Where("sender_id=? and target_id=?", userId, targetUserId).
			Or("sender_id=? and target_id=?", targetUserId, userId).
			Order("id desc").Limit(needNumber).Find(&messages).Error
	}
	return messages, err
}

func GetUserJoinGroupList(userId int, beginId int, needNumber int) ([]*FriendGroup, int, error) {
	var groups []*FriendGroup
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return groups, 0, err
	}
	count := mysqlDb.Model(&user).Association("JoinedGroups").Count()
	err = mysqlDb.Model(&user).Preload("Members").Offset(beginId).Limit(needNumber).
		Related(&groups, "JoinedGroups").Error
	return groups, count, err
}

func AddGroup(userId int, groupName string, groupMembers []int) error {
	var members []*User
	groupMembers = append(groupMembers, userId)
	err := mysqlDb.Where(groupMembers).Find(&members).Error
	if err != nil {
		return err
	}
	group := FriendGroup{
		UserID:  uint(userId),
		Name:    groupName,
		Members: members,
	}
	return mysqlDb.Create(&group).Save(&group).Error
}

func AddMessageInGroup(userId int, groupId int, content string) (uint, error) {
	message := MessageInGroup{
		SenderID: uint(userId),
		GroupID:  uint(groupId),
		Content:  content,
	}
	err := mysqlDb.Create(&message).Save(&message).Error
	return message.ID, err
}

func GetHistoryMessageInGroup(groupId int, oldestMessageId int, needNumber int) ([]MessageInGroup, error) {
	var messages []MessageInGroup
	var err error
	if oldestMessageId > 0 {
		err = mysqlDb.Preload("Sender").Where("group_id=? and id < ?", groupId, oldestMessageId).
			Order("id desc").Limit(needNumber).Find(&messages).Error
	} else {
		err = mysqlDb.Preload("Sender").Where("group_id=?", groupId).
			Order("id desc").Limit(needNumber).Find(&messages).Error
	}
	return messages, err
}

func GetLatestMessageInGroup(groupId int, latestMessageId int, needNumber int) ([]MessageInGroup, error) {
	var messages []MessageInGroup
	err := mysqlDb.Preload("Sender").Where("group_id=? and id > ?", groupId, latestMessageId).
		Order("id desc").Limit(needNumber).Find(&messages).Error
	println(messages)
	return messages, err
}
