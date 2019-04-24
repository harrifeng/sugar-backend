package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

func AddTopic(userId int, content string) error {
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return err
	}
	topic := Topic{
		Content: content,
		User:    user,
	}
	mysqlDb.Create(&topic)
	mysqlDb.Save(&topic)
	return nil
}

func AddTopicLordReply(userId int, topicId int, content string) error {
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return err
	}
	topicLordReply := TopicLordReply{
		Content: content,
		User:    user,
	}
	var topicTmp Topic
	return Transaction(func(db *gorm.DB) error {
		db.First(&topicTmp, topicId)
		err = db.Model(&topicTmp).Association("LordReplies").Append(topicLordReply).Error
		if err != nil {
			return err
		}
		topicTmp.UpdatedAt = time.Now()
		return db.Save(&topicTmp).Error
	})
}

func AddTopicLayerReply(userId int, topicLordReplyId int, content string) error {
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return err
	}
	topicLayerReply := TopicLayerReply{
		Content: content,
		User:    user,
	}
	topicLordReplyTmp, err := GetTopicLordReplyFromTopicLordReplyId(topicLordReplyId)
	if err != nil {
		return err
	}
	topicId := topicLordReplyTmp.TopicID
	topic, err := GetTopicFromTopicId(int(topicId))
	topic.UpdatedAt = time.Now()
	if err != nil {
		return err
	}
	return Transaction(func(db *gorm.DB) error {
		err = db.Model(&topicLordReplyTmp).Association("LayerReplies").Append(topicLayerReply).Error
		if err != nil {
			return err
		}
		return db.Save(&topic).Error
	})
}

func GetTopicFromTopicId(topicId int) (Topic, error) {
	var topic Topic
	err := mysqlDb.Preload("User").First(&topic, topicId).Error
	return topic, err
}

func GetTopicLordReplyFromTopicLordReplyId(topicLordReplyId int) (TopicLordReply, error) {
	var topicLordReply TopicLordReply
	err := mysqlDb.First(&topicLordReply, topicLordReplyId).Error
	return topicLordReply, err
}

func GetTopicLayerReplyFromTopicLayerReplyId(topicLayerReplyId int) (TopicLayerReply, error) {
	var topicLayerReply TopicLayerReply
	err := mysqlDb.First(&topicLayerReply, topicLayerReplyId).Error
	return topicLayerReply, err
}

func GetTopicLordReplyListFromTopicId(topicId int, beginId int, needNumber int) ([]TopicLordReply, error) {
	var topicLordReplies []TopicLordReply
	topic, err := GetTopicFromTopicId(topicId)
	if err != nil {
		return topicLordReplies, err
	}
	err = mysqlDb.Model(&topic).Preload("User").Offset(beginId).Limit(needNumber).
		Related(&topicLordReplies, "LordReplies").Error
	return topicLordReplies, err
}

func GetTopicLayerReplyListFromTopicLordReplyId(topicLordReplyId int, beginId int,
	needNumber int) ([]TopicLayerReply, error) {
	var topicLayerReplies []TopicLayerReply
	topicLordReply, err := GetTopicLordReplyFromTopicLordReplyId(topicLordReplyId)
	if err != nil {
		return topicLayerReplies, err
	}
	err = mysqlDb.Model(&topicLordReply).Preload("User").Offset(beginId).Limit(needNumber).
		Related(&topicLayerReplies, "LayerReplies").Error
	return topicLayerReplies, err
}

func RemoveTopic(topicId int) error {
	var topic Topic
	mysqlDb.Preload("LordReplies").First(&topic, topicId)
	for _, v := range topic.LordReplies {
		err := RemoveTopicLordReply(int(v.ID))
		if err != nil {
			return err
		}
	}
	err := mysqlDb.Delete(&topic).Error
	return err
}

func RemoveTopicLordReply(topicLordReplyId int) error {
	var topicLordReply TopicLordReply
	mysqlDb.Preload("LayerReplies").First(&topicLordReply, topicLordReplyId)
	for _, v := range topicLordReply.LayerReplies {
		err := RemoveTopicLayerReply(strconv.Itoa(int(v.ID)))
		if err != nil {
			return err
		}
	}
	err := mysqlDb.Delete(&topicLordReply).Error
	return err
}

func RemoveTopicLayerReply(TopicLayerReplyId string) error {
	topicLayerReplyId, _ := strconv.Atoi(TopicLayerReplyId)
	var topicLayerReply TopicLayerReply
	mysqlDb.First(&topicLayerReply, topicLayerReplyId)
	err := mysqlDb.Delete(&topicLayerReply).Error
	return err
}

func AddUserCollectedTopic(userId int, topicId int) error {
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return err
	}
	var topic Topic
	mysqlDb.First(&topic, topicId)
	err = mysqlDb.Model(&user).Association("CollectedTopics").Append(topic).Error
	return err
}

func RemoveUserCollectedTopic(userId int, topicId int) error {
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return err
	}
	var topic Topic
	mysqlDb.First(&topic, topicId)
	err = mysqlDb.Model(&user).Association("CollectedTopics").Delete(topic).Error
	return err
}

func GetSearchTopicList(searchContent string, beginId int, needNumber int) ([]Topic, error) {
	var topics []Topic
	err := mysqlDb.Where("content LIKE ?",
		fmt.Sprintf("%%%s%%", searchContent)).Offset(beginId).Limit(needNumber).Find(&topics).Error
	return topics, err
}

func GetUserCollectedTopicList(userId int, beginId int, needNumber int) ([]Topic, int, error) {
	var topics []Topic
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return topics, 0, err
	}
	err = mysqlDb.Model(&user).Preload("User").Offset(beginId).Limit(needNumber).
		Related(&topics, "CollectedTopics").Error
	if err != nil {
		return topics, 0, err
	}
	count := mysqlDb.Model(&user).Association("CollectedTopics").Count()
	return topics, count, err
}

func GetUserPublishedTopicList(userId int, beginId int, needNumber int) ([]Topic, int, error) {
	var topics []Topic
	var count int
	err := mysqlDb.Table("topics").Where("user_id=? and deleted_at is null", userId).Count(&count).
		Offset(beginId).Limit(needNumber).Find(&topics).Error
	return topics, count, err
}

func GetUserReplyList(userId int, beginId int, needNumber int) ([]UserReply, int, error) {
	var replies []UserReply
	err := mysqlDb.Raw(
		`(select id as topic_lord_reply_key,NULL as topic_layer_reply_key,created_at,content,user_id,thumbs_up_count,
			topic_id as father_id from topic_lord_replies where user_id=? and deleted_at is null)
			union all
			(select  NULL as topic_lord_reply_key,id as topic_layer_reply_key,created_at,content,user_id,thumbs_up_count,
			topic_lord_reply_id as father_id from topic_layer_replies where user_id=? and deleted_at is null)`,
		userId, userId).Scan(&replies).Error
	if beginId+needNumber > len(replies) {
		return replies[beginId:], len(replies), err
	}
	return replies[beginId : beginId+needNumber], len(replies), err
}

func GetLatestTopicList(TopicList []string, needNumber int) ([]Topic, error) {
	var topics []Topic
	topicList := make([]int, len(TopicList))
	for i, topicId := range TopicList {
		topicList[i], _ = strconv.Atoi(topicId)
	}
	err := mysqlDb.Preload("User").Not(topicList).Order("updated_at desc,id desc").
		Limit(needNumber).Find(&topics).Error
	return topics, err
}

func GetTopicReplyCount(topicId int) (int, error) {
	topic, err := GetTopicFromTopicId(topicId)
	if err != nil {
		return 0, err
	}
	var topicLordReplies []TopicLordReply
	count := mysqlDb.Preload("LayerReplies").Model(&topic).
		Association("LordReplies").Find(&topicLordReplies).Count()
	for _, topic := range topicLordReplies {
		count += len(topic.LayerReplies)
	}
	return count, nil
}

func CheckUserCollectedTopic(userId int, topicId int) (bool, error) {
	var count int
	err := mysqlDb.Table("user_collected_topic").
		Where("user_id=? and topic_id=?", userId, topicId).Count(&count).Error
	return count > 0, err
}

func GetTopicLayerReplyCountFromTopicLordReply(topicLordReplyId int) (int, error) {
	topicLordReply, err := GetTopicLordReplyFromTopicLordReplyId(topicLordReplyId)
	if err != nil {
		return 0, err
	}
	count := mysqlDb.Model(&topicLordReply).Association("LayerReplies").Count()
	return count, nil
}

func ValueTopic(topicId int, value int) error {
	topic, err := GetTopicFromTopicId(topicId)
	if err != nil {
		return err
	}
	topic.ThumbsUpCount += value
	return mysqlDb.Save(&topic).Error
}

func ValueTopicLordReply(topicLordReplyId int, value int) error {
	topicLordReply, err := GetTopicLordReplyFromTopicLordReplyId(topicLordReplyId)
	if err != nil {
		return err
	}
	topicLordReply.ThumbsUpCount += value
	return mysqlDb.Save(&topicLordReply).Error
}

func ValueTopicLayerReply(topicLayerReplyId int, value int) error {
	topicLayerReply, err := GetTopicLayerReplyFromTopicLayerReplyId(topicLayerReplyId)
	if err != nil {
		return err
	}
	topicLayerReply.ThumbsUpCount += value
	return mysqlDb.Save(&topicLayerReply).Error
}

func GetTopicLordReplyFloor(topicLordReplyId int) (int, error) {
	topicLordReply, err := GetTopicLordReplyFromTopicLordReplyId(topicLordReplyId)
	if err != nil {
		return 0, err
	}
	topic, err := GetTopicFromTopicId(int(topicLordReply.TopicID))
	var topicLordRelies []TopicLordReply
	err = mysqlDb.Unscoped().Model(&topic).Related(&topicLordRelies, "LordReplies").
		Where("id <= ?", topicLordReplyId).Error
	return len(topicLordRelies), err
}

func GetTopicCollectingUserCount(TopicId string) (int, error) {
	var count int
	err := mysqlDb.Table("user_collected_topic").Where("topic_id=?", TopicId).Count(&count).Error
	return count, err
}

func GetLastFloorFromTopicId(TopicId string) (int, error) {
	var count int
	err := mysqlDb.Table("topic_lord_replies").Where("topic_id=?", TopicId).Count(&count).Error
	return count, err
}
