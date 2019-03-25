package db

func AddTopic(UserId string, Content string, PictureBase64 []string) error {
	return nil
}

func AddTopicLordReply(UserId string, TopicId string, Content string, PictureBase64 []string) error {
	return nil
}

func AddTopicLayerReply(UserId string, TopicLordReplyId string, Content string) error {
	return nil
}

func GetTopicFromTopicId(TopicId string) (Topic, error) {
	return Topic{}, nil
}

func GetTopicLordReplyFromTopicId(TopicId string, BeginId string, NeedNumber string) ([]TopicLordReply, error) {
	var topicLordReplies []TopicLordReply
	return topicLordReplies, nil
}

func GetTopicLayerReplyFromTopicLordReplyId(TopicLordReplyId string, BeginId string, NeedNumber string) ([]TopicLayerReply, error) {
	var topicLayerReplies []TopicLayerReply
	return topicLayerReplies, nil
}

func RemoveTopic(TopicId string) error {
	return nil
}

func RemoveTopicLordReply(TopicLordReplyId string) error {
	return nil
}

func AddUserCollectedTopic(UserId string, TopicId string) error {
	return nil
}

func RemoveUserCollectedTopic(UserId string, TopicId string) error {
	return nil
}

func GetSearchTopicList(SearchContent string, BeginId string, NeedNumber string) ([]Topic, error) {
	var topics []Topic
	return topics, nil
}

func GetUserCollectedTopicList(UserId string, BeginId string, NeedNumber string) ([]Topic, error) {
	var topics []Topic
	return topics, nil
}

func GetUserPublishedTopicList(UserId string, BeginId string, NeedNumber string) ([]Topic, error) {
	var topics []Topic
	return topics, nil
}

func GetUserReplyList(UserId string, BeginId string, NeedNumber string) ([]interface{}, error) {
	var replies []interface{}
	return replies, nil
}
