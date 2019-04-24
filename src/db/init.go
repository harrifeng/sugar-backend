package db

import (
	"fmt"
)

func AutoCreateTable() {
	mysqlDb.AutoMigrate(&User{})
	mysqlDb.AutoMigrate(&Article{})
	mysqlDb.AutoMigrate(&ArticleComment{})
	mysqlDb.AutoMigrate(&Topic{})
	mysqlDb.AutoMigrate(&TopicLordReply{})
	mysqlDb.AutoMigrate(&TopicLayerReply{})
	mysqlDb.AutoMigrate(&UserPrivacySetting{})
	mysqlDb.AutoMigrate(&UserCheckIn{})
	mysqlDb.AutoMigrate(&FamilyMember{})
	mysqlDb.AutoMigrate(&BloodRecord{})
	mysqlDb.AutoMigrate(&HealthRecord{})
	mysqlDb.AutoMigrate(&MessageU2u{})
	mysqlDb.AutoMigrate(&MessageInGroup{})
	mysqlDb.AutoMigrate(&FriendGroup{})
	mysqlDb.AutoMigrate(&SugarGuideDietPlan{})
	mysqlDb.AutoMigrate(&SugarGuideSportPlan{})
	mysqlDb.AutoMigrate(&SugarGuideControlPlan{})
	fmt.Printf("tables created successfully!\n ")
}
