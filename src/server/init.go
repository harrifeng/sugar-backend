package server

import (
	"db"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"runtime"
	"strings"
	"utils"
)

func InitAllArticle() {
	ys137Path := filepath.Join(serverConfig.DataDirPath, "articles", "ys137", "*")
	ys137Files, err := filepath.Glob(ys137Path)
	if err != nil {
		log.Fatal(err)
		return
	}
	for i, file := range ys137Files {
		var s []string
		if runtime.GOOS == "windows" {
			s = strings.Split(file, "\\")
		} else {
			s = strings.Split(file, "/")
		}
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
		_ = db.AddNewArticle(title, string(content), string(plainContent))
		fmt.Printf("article.%d created successfully!\n ", i)
	}
}

func InitUser() {
	for i := 1; i <= 100; i++ {
		_ = db.CreateNewUser(utils.RandPhoneNumber(), utils.RandUsername(), "0c152176187ce61c9614c072e1a1e2f7")
		fmt.Printf("user.%d created successfully!\n ", i)
	}
}

func InitArticleComment() {
	for i := 1; i <= 1000; i++ {
		userId := rand.Intn(10) + 1
		articleId := rand.Intn(100) + 1
		_ = db.AddArticleComment(userId, articleId, utils.RandWords())
		fmt.Printf("articleComment.%d created successfully!\n ", i)
	}
}

func AutoCreateTable() {
	db.AutoCreateTable()
}

func InitTopicTest() {
	for i := 1; i <= 500; i++ {
		userId := rand.Intn(100) + 1
		_ = db.AddTopic(userId, utils.RandWords())
		fmt.Printf("topic.%d created successfully!\n", i)
	}
}

func InitTopicLordReplyTest() {
	for i := 1; i <= 500; i++ {
		userId := rand.Intn(100) + 1
		topicId := rand.Intn(100) + 1
		_ = db.AddTopicLordReply(userId, topicId, utils.RandWords())
		fmt.Printf("topicLordReply.%d created successfully!\n", i)
	}
}

func InitTopicLayerReplyTest() {
	for i := 1; i <= 500; i++ {
		userId := rand.Intn(100) + 1
		topicLordReplyId := rand.Intn(100) + 1
		_ = db.AddTopicLayerReply(userId, topicLordReplyId, utils.RandWords())
		fmt.Printf("topicLordReply.%d created successfully!\n", i)
	}
}

func Init() {
	AutoCreateTable()
	InitAllArticle()
	InitUser()
	InitArticleComment()
	InitTopicTest()
	InitTopicLordReplyTest()
	InitTopicLayerReplyTest()
}
