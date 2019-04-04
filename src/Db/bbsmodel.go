package db

import (
	"fmt"
	"strconv"
	"time"
)

func AddTopic(UserId string, Content string) error {
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return err
	}
	topic := Topic{
		Content: Content,
		User:    user,
	}
	mysqlDb.Create(&topic)
	mysqlDb.Save(&topic)
	return nil
}

func AddTopicLordReply(UserId string, TopicId string, Content string) error {
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return err
	}
	topicLordReply := TopicLordReply{
		Content: Content,
		User:    user,
	}
	var topicTmp Topic
	topicId, _ := strconv.Atoi(TopicId)
	mysqlDb.First(&topicTmp, topicId)
	tx := mysqlDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err = tx.Model(&topicTmp).Association("LordReplies").Append(topicLordReply).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	topicTmp.UpdatedAt = time.Now()
	err = tx.Save(&topicTmp).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	return err
}

func AddTopicLayerReply(UserId string, TopicLordReplyId string, Content string) error {
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return err
	}
	topicLayerReply := TopicLayerReply{
		Content: Content,
		User:    user,
	}
	topicLordReplyTmp, err := GetTopicLordReplyFromTopicLordReplyId(TopicLordReplyId)
	if err != nil {
		return err
	}
	topicId := topicLordReplyTmp.TopicID
	topic, err := GetTopicFromTopicId(strconv.Itoa(int(topicId)))
	topic.UpdatedAt = time.Now()
	if err != nil {
		return err
	}
	tx := mysqlDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err = tx.Model(&topicLordReplyTmp).Association("LayerReplies").Append(topicLayerReply).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&topic).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	return err
}

func GetTopicFromTopicId(TopicId string) (Topic, error) {
	var topic Topic
	topicId, _ := strconv.Atoi(TopicId)
	err := mysqlDb.Preload("User").First(&topic, topicId).Error
	return topic, err
}

func GetTopicLordReplyFromTopicLordReplyId(TopicLordReplyId string) (TopicLordReply, error) {
	var topicLordReply TopicLordReply
	topicLordReplyId, _ := strconv.Atoi(TopicLordReplyId)
	err := mysqlDb.First(&topicLordReply, topicLordReplyId).Error
	return topicLordReply, err
}

func GetTopicLayerReplyFromTopicLayerReplyId(TopicLayerReplyId string) (TopicLayerReply, error) {
	var topicLayerReply TopicLayerReply
	topicLayerReplyId, _ := strconv.Atoi(TopicLayerReplyId)
	err := mysqlDb.First(&topicLayerReply, topicLayerReplyId).Error
	return topicLayerReply, err
}

func GetTopicLordReplyListFromTopicId(TopicId string, BeginId string, NeedNumber string) ([]TopicLordReply, error) {
	var topicLordReplies []TopicLordReply
	topic, err := GetTopicFromTopicId(TopicId)
	if err != nil {
		return topicLordReplies, err
	}
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	err = mysqlDb.Model(&topic).Preload("User").Offset(beginId).Limit(needNumber).
		Related(&topicLordReplies, "LordReplies").Error
	return topicLordReplies, err
}

func GetTopicLayerReplyListFromTopicLordReplyId(TopicLordReplyId string, BeginId string,
	NeedNumber string) ([]TopicLayerReply, error) {
	var topicLayerReplies []TopicLayerReply
	topicLordReply, err := GetTopicLordReplyFromTopicLordReplyId(TopicLordReplyId)
	if err != nil {
		return topicLayerReplies, err
	}
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	err = mysqlDb.Model(&topicLordReply).Preload("User").Offset(beginId).Limit(needNumber).
		Related(&topicLayerReplies, "LayerReplies").Error
	return topicLayerReplies, err
}

func RemoveTopic(TopicId string) error {
	topicId, _ := strconv.Atoi(TopicId)
	var topic Topic
	mysqlDb.First(&topic, topicId)
	err := mysqlDb.Delete(&topic).Error
	return err
}

func RemoveTopicLordReply(TopicLordReplyId string) error {
	topicLordReplyId, _ := strconv.Atoi(TopicLordReplyId)
	var topicLordReply TopicLordReply
	mysqlDb.First(&topicLordReply, topicLordReplyId)
	err := mysqlDb.Delete(&topicLordReply).Error
	return err
}

func AddUserCollectedTopic(UserId string, TopicId string) error {
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return err
	}
	topicId, _ := strconv.Atoi(TopicId)
	var topic Topic
	mysqlDb.First(&topic, topicId)
	err = mysqlDb.Model(&user).Association("CollectedTopics").Append(topic).Error
	return err
}

func RemoveUserCollectedTopic(UserId string, TopicId string) error {
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return err
	}
	topicId, _ := strconv.Atoi(TopicId)
	var topic Topic
	mysqlDb.First(&topic, topicId)
	err = mysqlDb.Model(&user).Association("CollectedTopics").Delete(topic).Error
	return err
}

func GetSearchTopicList(SearchContent string, BeginId string, NeedNumber string) ([]Topic, error) {
	var topics []Topic
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	err := mysqlDb.Where("content LIKE ?",
		fmt.Sprintf("%%%s%%", SearchContent)).Offset(beginId).Limit(needNumber).Find(&topics).Error
	return topics, err
}

func GetUserCollectedTopicList(UserId string, BeginId string, NeedNumber string) ([]Topic, int, error) {
	var topics []Topic
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return topics, 0, err
	}
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	err = mysqlDb.Model(&user).Offset(beginId).Limit(needNumber).Related(&topics, "CollectedTopics").Error
	if err != nil {
		return topics, 0, err
	}
	count := mysqlDb.Model(&user).Association("CollectedTopics").Count()
	return topics, count, err
}

func GetUserPublishedTopicList(UserId string, BeginId string, NeedNumber string) ([]Topic, int, error) {
	var topics []Topic
	userId, _ := strconv.Atoi(UserId)
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	err := mysqlDb.Preload("CollectingUsers").Where(&Topic{UserID: userId}).
		Offset(beginId).Limit(needNumber).Find(&topics).Error
	if err != nil {
		return topics, 0, err
	}
	var count int
	err = mysqlDb.Where(&Topic{UserID: userId}).Count(&count).Error
	return topics, count, err
}

func GetUserReplyList(UserId string, BeginId string, NeedNumber string) ([]UserReply, int, error) {
	var replies []UserReply
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	var count int
	err := mysqlDb.Raw(
		`(select content,user_id,thumbs_up_count,topic_id,NULL as topic_lord_reply_id from topic_lord_replies
	where user_id=?)
	union all
	(select content,user_id,thumbs_up_count,NULL as topic_id,topic_lord_reply_id from topic_layer_replies
		where user_id=?)`, UserId, UserId).Count(&count).Offset(beginId).Limit(needNumber).Scan(&replies).Error
	return replies, count, err
}

func GetLatestTopicList(TopicList []string, NeedNumber string) ([]Topic, error) {
	var topics []Topic
	needNumber, _ := strconv.Atoi(NeedNumber)
	topicList := make([]int, len(TopicList))
	for i, topicId := range TopicList {
		topicList[i], _ = strconv.Atoi(topicId)
	}
	err := mysqlDb.Preload("User").Not(topicList).Order("updated_at desc,id desc").
		Limit(needNumber).Find(&topics).Error
	return topics, err
}

func GetTopicReplyCount(TopicId string) (int, error) {
	topic, err := GetTopicFromTopicId(TopicId)
	if err != nil {
		return 0, err
	}
	var topicLordReplies []TopicLordReply
	count := mysqlDb.Preload("LayerReplies").Model(&topic).Association("LordReplies").Find(&topicLordReplies).Count()
	for _, topic := range topicLordReplies {
		count += len(topic.LayerReplies)
	}
	return count, nil
}

func CheckUserCollectedTopic(UserId string, TopicId string) (bool, error) {
	var count int
	err:=mysqlDb.Table("user_collected_topic").
		Where("user_id=? and topic_id=?",UserId,TopicId).Count(&count).Error
	return count>0, err
}

func GetTopicLayerReplyCountFromTopicLordReply(TopicLordReplyId string) (int, error) {
	topicLordReply, err := GetTopicLordReplyFromTopicLordReplyId(TopicLordReplyId)
	if err != nil {
		return 0, err
	}
	count := mysqlDb.Model(&topicLordReply).Association("LayerReplies").Count()
	return count, nil
}

func ValueTopic(TopicId string, Value string) error {
	topic, err := GetTopicFromTopicId(TopicId)
	if err != nil {
		return err
	}
	value, _ := strconv.Atoi(Value)
	topic.ThumbsUpCount += value
	return mysqlDb.Save(&topic).Error
}

func ValueTopicLordReply(TopicLordReplyId string, Value string) error {
	topicLordReply, err := GetTopicLordReplyFromTopicLordReplyId(TopicLordReplyId)
	if err != nil {
		return err
	}
	value, _ := strconv.Atoi(Value)
	topicLordReply.ThumbsUpCount += value
	return mysqlDb.Save(&topicLordReply).Error
}

func ValueTopicLayerReply(TopicLayerReplyId string, Value string) error {
	topicLayerReply, err := GetTopicLayerReplyFromTopicLayerReplyId(TopicLayerReplyId)
	if err != nil {
		return err
	}
	value, _ := strconv.Atoi(Value)
	topicLayerReply.ThumbsUpCount += value
	return mysqlDb.Save(&topicLayerReply).Error
}

func GetTopicLordReplyFloor(TopicLordReplyId string) (int, error) {
	topicLordReply, err := GetTopicLordReplyFromTopicLordReplyId(TopicLordReplyId)
	if err != nil {
		return 0, err
	}
	topicId := strconv.Itoa(int(topicLordReply.TopicID))
	topic, err := GetTopicFromTopicId(topicId)
	var topicLordRelies []TopicLordReply
	err = mysqlDb.Unscoped().Model(&topic).Related(&topicLordRelies, "LordReplies").
		Where("id <= ?", TopicLordReplyId).Error
	return len(topicLordRelies), err

}
