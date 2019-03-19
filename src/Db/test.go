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

	//AutoCreateTableTest()
	//user := User{UserName: "haha", Age: 122, PhoneNumber: "213213213"}
	//err:=CreateNewUser(user)
	//if err!=nil{
	//	fmt.Println(err)
	//}
}
