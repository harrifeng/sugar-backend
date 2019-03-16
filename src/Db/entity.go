package Db

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model

	// User Basic data
	UserName    string
	PhoneNumber string
	Password    string
	Gender      string
	Age         int

	// User Account data
	Exp               int
	Level             int
	HeadPortraitUrl   string
	FollowUsers       []*User   `gorm:"many2many:user_follow_ships;association_jointable_foreignkey:follow_user_id"`
	CollectedArticles []Article `gorm:"many2many:user_collected_article"`
	CollectedTopics   []Topic   `gorm:"many2many:user_collected_topic"`

	// User healthy data
	Height int
	Weight int
	Area   string
	Job    string
}

type ArticleLabel struct {
	gorm.Model
	Value    string
	Articles []*Article `gorm:"many2many:articles_to_labels;"`
}

type Article struct {
	gorm.Model
	Title         string
	Content       string
	Labels        []*ArticleLabel `gorm:"many2many:articles_to_labels;"`
	CoverImageUrl string
	ReadCount     int
}

type ArticleComment struct {
	gorm.Model
	Content       string
	ThumbsUpCount int
	Article       Article
	User          User
	ArticleID     int
	UserID        int
}

type Topic struct {
	gorm.Model
	Content       string
	User          User
	PictureUrls   string
	ThumbsUpCount int
	LordReplies   []TopicLordReply
}

type TopicLordReply struct {
	gorm.Model
	Content       string
	User          User
	PictureUrls   string
	ThumbsUpCount int
	TopicID       uint
	LayerReplies  []TopicLayerReply
}

type TopicLayerReply struct {
	gorm.Model
	Content          string
	User             User
	ThumbsUpCount    int
	TopicLordReplyID uint
}
