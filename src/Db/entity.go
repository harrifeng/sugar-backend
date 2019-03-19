package db

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model

	// User Basic data
	UserName    string
	PhoneNumber string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	Gender      string
	Age         int

	// User Account data
	Exp               int       `gorm:"not null;default:'0'"`
	Level             int       `gorm:"not null;default:'1'"`
	HeadPortraitUrl   string    ``
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
	Title         string          `gorm:"not null"`
	Content       string          `gorm:"not null"`
	Labels        []*ArticleLabel `gorm:"many2many:articles_to_labels;"`
	CoverImageUrl string
	ReadCount     int `gorm:"not null;default:'0'"`
}

type ArticleComment struct {
	gorm.Model
	Content       string `gorm:"not null"`
	ThumbsUpCount int    `gorm:"not null;default:'0'"`
	Article       Article
	User          User
	ArticleID     int
	UserID        int
}

type Topic struct {
	gorm.Model
	Content       string `gorm:"not null"`
	User          User
	PictureUrls   string
	ThumbsUpCount int `gorm:"not null;default:'0'"`
	LordReplies   []TopicLordReply
}

type TopicLordReply struct {
	gorm.Model
	Content       string `gorm:"not null"`
	User          User
	PictureUrls   string
	ThumbsUpCount int `gorm:"not null;default:'0'"`
	TopicID       uint
	LayerReplies  []TopicLayerReply
}

type TopicLayerReply struct {
	gorm.Model
	Content          string `gorm:"not null"`
	User             User
	ThumbsUpCount    int `gorm:"not null;default:'0'"`
	TopicLordReplyID uint
}
