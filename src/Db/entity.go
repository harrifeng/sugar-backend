package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model

	// User Basic data
	UserName    string `gorm:"type:text;"`
	PhoneNumber string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	Gender      string `gorm:"type:text;"`
	Age         int

	// User Account data
	Exp                  int       `gorm:"not null;default:'0'"`
	Level                int       `gorm:"not null;default:'1'"`
	HeadPortraitUrl      string    ``
	FollowingUsers       []*User   `gorm:"many2many:user_following_ships;association_jointable_foreignkey:following_user_id"`
	FollowerUsers        []*User   `gorm:"many2many:user_follower_ships;association_jointable_foreignkey:follower_user_id"`
	CollectedArticles    []Article `gorm:"many2many:user_collected_article"`
	CollectedTopics      []Topic   `gorm:"many2many:user_collected_topic"`
	BloodRecords         []BloodRecord
	HealthRecords        []HealthRecord
	FamilyMembers        []FamilyMember
	CheckIns             []UserCheckIn
	UserPrivacySetting   UserPrivacySetting
	UserPrivacySettingID uint

	// User healthy data
	Height float64
	Weight float64
	Area   string
	Job    string
}

type UserCheckIn struct {
	gorm.Model
	CheckTime time.Time
	User      User
	UserID    uint
}

type UserPrivacySetting struct {
	gorm.Model
	ShowPhoneNumber bool
	ShowGender      bool
	ShowAge         bool
	ShowHeight      bool
	ShowWeight      bool
	ShowArea        bool
	ShowJob         bool
}

type Article struct {
	gorm.Model
	Title           string `gorm:"not null"`
	Content         string `gorm:"not null;type:text;"`
	CoverImageUrl   string
	ReadCount       int `gorm:"not null;default:'0'"`
	ArticleComments []ArticleComment
}

type ArticleComment struct {
	gorm.Model
	Content       string `gorm:"not null;type:text"`
	ThumbsUpCount int    `gorm:"not null;default:'0'"`
	User          User
	UserID        int
	Article       Article
	ArticleID     int
}

type Topic struct {
	gorm.Model
	Content         string `gorm:"not null;type:text"`
	User            User
	UserID          int
	ThumbsUpCount   int `gorm:"not null;default:'0'"`
	LordReplies     []TopicLordReply
}

type TopicLordReply struct {
	gorm.Model
	Content       string `gorm:"not null;type:text"`
	User          User
	UserID        int
	ThumbsUpCount int `gorm:"not null;default:'0'"`
	TopicID       uint
	LayerReplies  []TopicLayerReply
}

type TopicLayerReply struct {
	gorm.Model
	Content          string `gorm:"not null;type:text"`
	User             User
	UserID           int
	ThumbsUpCount    int `gorm:"not null;default:'0'"`
	TopicLordReplyID uint
}

type BloodRecord struct {
	gorm.Model
	Level      string
	RecordTime string
	RecordDate time.Time `gorm:"type:date;"`
	User       User
	UserID     uint
}

type HealthRecord struct {
	gorm.Model
	Insulin       string
	SportTime     string
	Weight        string
	BloodPressure string
	RecordTime    string
	RecordDate    time.Time `gorm:"type:date;"`
	User          User
	UserID        uint
}

type FamilyMember struct {
	gorm.Model
	PhoneNumber string
	Call        string
	User        User
	UserID      uint
}

type UserReply struct {
	gorm.Model
	Content          string
	UserId           int
	ThumbsUpCount    int
	TopicID          uint
	TopicLordReplyID uint
}
