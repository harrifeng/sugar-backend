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
		Limit(needNumber).Find(&messages).Error
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

func GetUserJoinGroupList(userId int, beginId int, needNumber int) ([]FriendGroup, int, error) {
	var groups []FriendGroup
	var count int
	mysqlDb.Model(&FriendGroup{}).Where("user_id=?", userId).Count(&count)
	err := mysqlDb.Where("user_id=?", userId).Offset(beginId).Limit(needNumber).Find(&groups).Error
	return groups, count, err
}
