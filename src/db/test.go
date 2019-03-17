package db

import "fmt"

func SetNewVerificationCodeTest() {
	err := SetNewVerificationCode("2323", "21321")
	if err != nil {
		fmt.Printf("%s", err)
	}
}

func AutoCreateTableTest() {
	mysqlDb.AutoMigrate(&User{})
	mysqlDb.AutoMigrate(&Article{})
	mysqlDb.AutoMigrate(&ArticleComment{})
	mysqlDb.AutoMigrate(&ArticleLabel{})
	mysqlDb.AutoMigrate(&Topic{})
	mysqlDb.AutoMigrate(&TopicLordReply{})
	mysqlDb.AutoMigrate(&TopicLayerReply{})
}

func Test() {
	mysqlDb.AutoMigrate(&User{})
	user := User{UserName: "haha", Age: 122, PhoneNumber: "213213213"}
	mysqlDb.Create(&user)
	mysqlDb.Save(&user)
}
