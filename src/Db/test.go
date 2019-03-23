package db

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func SetNewVerificationCodeTest() {
	err := SetNewVerificationCode("2323", "21321")
	if err != nil {
		fmt.Printf("%s", err)
	}
}

func InitAllArticle() {
	ys137Path := filepath.Join(TestArticleDirPath, "ys137", "*")
	ys137Files, _ := filepath.Glob(ys137Path)
	for i, file := range ys137Files {
		if i > 100 {
			break
		}
		s := strings.Split(file, "\\")
		fileShortName := s[len(s)-1]
		title := fileShortName[:len(fileShortName)-4]
		content, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		_ = AddNewArticle(title, string(content))
		fmt.Printf("finish %d !\n ", i)
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
	mysqlDb.AutoMigrate(&UserPrivacySetting{})
}

func Test() {

	//AutoCreateTableTest()
	//user := User{UserName: "haha", Age: 122, PhoneNumber: "213213213"}
	//err:=CreateNewUser(user)
	//if err!=nil{
	//	fmt.Println(err)
	//}
}
