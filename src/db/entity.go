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
	Exp                     int            `gorm:"not null;default:'0'"`
	Level                   int            `gorm:"not null;default:'1'"`
	HeadPortraitUrl         string         ``
	FollowingUsers          []*User        `gorm:"many2many:user_following_ships;association_jointable_foreignkey:following_user_id"`
	FollowerUsers           []*User        `gorm:"many2many:user_follower_ships;association_jointable_foreignkey:follower_user_id"`
	CollectedArticles       []Article      `gorm:"many2many:user_collected_article"`
	CollectedTopics         []Topic        `gorm:"many2many:user_collected_topic"`
	JoinedGroups            []*FriendGroup `gorm:"many2many:user_joined_group"`
	BloodRecords            []BloodRecord
	HealthRecords           []HealthRecord
	FamilyMembers           []FamilyMember
	UserPrivacySetting      UserPrivacySetting
	UserPrivacySettingID    uint
	SugarGuideDietPlan      SugarGuideDietPlan
	SugarGuideDietPlanID    uint
	SugarGuideSportPlan     SugarGuideSportPlan
	SugarGuideSportPlanID   uint
	SugarGuideControlPlan   SugarGuideControlPlan
	SugarGuideControlPlanID uint
	// User healthy data
	Height float64
	Weight float64
	Area   string
	Job    string
}

type UserCheckIn struct {
	gorm.Model
	CheckTime time.Time `gorm:"type:date"`
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
	PlainContent    string `gorm:"not null;type:text;"`
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
	Content       string `gorm:"not null;type:text"`
	User          User
	UserID        int
	ThumbsUpCount int `gorm:"not null;default:'0'"`
	LordReplies   []TopicLordReply
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
	CallName    string
	User        User
	UserID      uint
}

type MessageU2u struct {
	gorm.Model
	SenderID uint
	Sender   User
	Content  string `gorm:"type:text"`
	TargetID uint
	Target   User
}

type FriendGroup struct {
	gorm.Model
	Name    string  `gorm:"unique;not null;"`
	Members []*User `gorm:"many2many:user_joined_group"`
	User    User
	UserID  uint
}

type MessageInGroup struct {
	gorm.Model
	SenderID uint
	Sender   User
	GroupID  uint
	Group    FriendGroup
	Content  string `gorm:"type:text"`
}

type SugarGuideDietPlan struct {
	gorm.Model
	Change     int `json:"change"`
	Cereals    int `json:"cereals"`
	Fruit      int `json:"fruit"`
	Meat       int `json:"meat"`
	Milk       int `json:"milk"`
	Fat        int `json:"fat"`
	Vegetables int `json:"vegetables"`
}

type SugarGuideSportPlan struct {
	gorm.Model
	Sport1 string `gorm:"type:text" json:"sport1"`
	Sport2 string `gorm:"type:text" json:"sport2"`
	Sport3 string `gorm:"type:text" json:"sport3"`
	Sport4 string `gorm:"type:text" json:"sport4"`
	Time1  int    `json:"time1"`
	Time2  int    `json:"time2"`
	Time3  int    `json:"time3"`
	Time4  int    `json:"time4"`
	Week1  string `json:"week1"`
	Week2  string `json:"week2"`
	Week3  string `json:"week3"`
	Week4  string `json:"week4"`
}

type SugarGuideControlPlan struct {
	gorm.Model
	Min1   float64 `json:"min1"`
	Min2   float64 `json:"min2"`
	Max1   float64 `json:"max1"`
	Max2   float64 `json:"max2"`
	Sleep1 float64 `json:"sleep1"`
	Sleep2 float64 `json:"sleep2"`
}

type UserReply struct {
	TopicLordReplyKey  uint      `gorm:"Column:topic_lord_reply_key"`
	TopicLayerReplyKey uint      `gorm:"Column:topic_layer_reply_key"`
	CreatedAt          time.Time `gorm:"Column:created_at"`
	Content            string    `gorm:"Column:content"`
	UserId             int       `gorm:"Column:user_id"`
	ThumbsUpCount      int       `gorm:"Column:thumbs_up_count"`
	FatherID           uint      `gorm:"Column:father_id"`
}

type UserU2uMessage struct {
	Id        int
	Content   string
	OtherId   int
	Other     User
	CreatedAt time.Time
}
