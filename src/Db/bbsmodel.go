package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
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
	return Transaction(func(db *gorm.DB) error {
		err = db.Model(&topicLordReplyTmp).Association("LayerReplies").Append(topicLayerReply).Error
		if err != nil {
			return err
		}
		return db.Save(&topic).Error
	})
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
	mysqlDb.Preload("LordReplies").First(&topic, topicId)
	for _,v:=range topic.LordReplies{
		err := RemoveTopicLordReply(strconv.Itoa(int(v.ID)))
		if err!=nil{
			return err
		}
	}
	err := mysqlDb.Delete(&topic).Error
	return err
}

func RemoveTopicLordReply(TopicLordReplyId string) error {
	topicLordReplyId, _ := strconv.Atoi(TopicLordReplyId)
	var topicLordReply TopicLordReply
	mysqlDb.Preload("LayerReplies").First(&topicLordReply, topicLordReplyId)
	for _,v:=range topicLordReply.LayerReplies{
		err := RemoveTopicLayerReply(strconv.Itoa(int(v.ID)))
		if err!=nil{
			return err
		}
	}
	err := mysqlDb.Delete(&topicLordReply).Error
	return err
}

func RemoveTopicLayerReply(TopicLayerReplyId string)error{
	topicLayerReplyId, _ := strconv.Atoi(TopicLayerReplyId)
	var topicLayerReply TopicLayerReply
	mysqlDb.First(&topicLayerReply, topicLayerReplyId)
	err := mysqlDb.Delete(&topicLayerReply).Error
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
	err = mysqlDb.Model(&user).Preload("User").Offset(beginId).Limit(needNumber).
		Related(&topics, "CollectedTopics").Error
	if err != nil {
		return topics, 0, err
	}
	count := mysqlDb.Model(&user).Association("CollectedTopics").Count()
	return topics, count, err
}

func GetUserPublishedTopicList(UserId string, BeginId string, NeedNumber string) ([]Topic, int, error) {
	var topics []Topic
	var count int
	userId, _ := strconv.Atoi(UserId)
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	err := mysqlDb.Table("topics").Where("user_id=? and deleted_at is null", userId).Count(&count).
		Offset(beginId).Limit(needNumber).Find(&topics).Error
	return topics, count, err
}


func GetUserReplyList(UserId string, BeginId string, NeedNumber string) ([]UserReply, int, error) {

	var replies []UserReply
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	err := mysqlDb.Raw(
		`(select id as topic_lord_reply_key,NULL as topic_layer_reply_key,created_at,content,user_id,thumbs_up_count,
			topic_id as father_id from topic_lord_replies where user_id=? and deleted_at is null)
			union all
			(select  NULL as topic_lord_reply_key,id as topic_layer_reply_key,created_at,content,user_id,thumbs_up_count,
			topic_lord_reply_id as father_id from topic_layer_replies where user_id=? and deleted_at is null)`,
		UserId, UserId).Scan(&replies).Error
	if beginId + needNumber>len(replies){
		return replies[beginId : ], len(replies), err
	}
	return replies[beginId : beginId + needNumber], len(replies), err
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
	count := mysqlDb.Preload("LayerReplies").Model(&topic).
		Association("LordReplies").Find(&topicLordReplies).Count()
	for _, topic := range topicLordReplies {
		count += len(topic.LayerReplies)
	}
	return count, nil
}

func CheckUserCollectedTopic(UserId string, TopicId string) (bool, error) {
	var count int
	err := mysqlDb.Table("user_collected_topic").
		Where("user_id=? and topic_id=?", UserId, TopicId).Count(&count).Error
	return count > 0, err
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
