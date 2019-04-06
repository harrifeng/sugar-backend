package db

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"utils"
)

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
		s[len(s)-3] = "plain_articles"
		plainFile := strings.Join(s, "\\")
		plainContent, err := ioutil.ReadFile(plainFile)
		_ = AddNewArticle(title, string(content), string(plainContent))
		fmt.Printf("article.%d created successfully!\n ", i)
	}
}

func InitUser() {
	for i := 1; i <= 100; i++ {
		_ = CreateNewUser(utils.RandPhoneNumber(), utils.RandUserName(), "0c152176187ce61c9614c072e1a1e2f7")
		fmt.Printf("user.%d created successfully!\n ", i)
	}
}

func IninArticleComment() {
	for i := 1; i <= 1000; i++ {
		userId := rand.Intn(10) + 1
		articleId := rand.Intn(100) + 1
		_ = AddArticleComment(userId, articleId, utils.RandWords())
		fmt.Printf("articleComment.%d created successfully!\n ", i)
	}
}

func AutoCreateTableTest() {
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
	fmt.Printf("tables created successfully!\n ")
}

func InitTopicTest() {
	for i := 1; i <= 500; i++ {
		userId := rand.Intn(100) + 1
		_ = AddTopic(userId, utils.RandWords())
		fmt.Printf("topic.%d created successfully!\n", i)
	}
}

func InitTopicLordReplyTest() {
	for i := 1; i <= 500; i++ {
		userId := rand.Intn(100) + 1
		topicId := rand.Intn(100) + 1
		_ = AddTopicLordReply(userId, topicId, utils.RandWords())
		fmt.Printf("topicLordReply.%d created successfully!\n", i)
	}
}

func InitTopicLayerReplyTest() {
	for i := 1; i <= 500; i++ {
		userId := rand.Intn(100) + 1
		topicLordReplyId := rand.Intn(100) + 1
		_ = AddTopicLayerReply(userId, topicLordReplyId, utils.RandWords())
		fmt.Printf("topicLordReply.%d created successfully!\n", i)
	}
}

func Init() {
	AutoCreateTableTest()
	//InitAllArticle()
	//InitUser()
	//IninArticleComment()
	//InitTopicTest()
	//InitTopicLordReplyTest()
	//InitTopicLayerReplyTest()
}
