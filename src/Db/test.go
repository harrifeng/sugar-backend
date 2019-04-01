package db

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"utils"
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
		_ = AddArticleComment(strconv.Itoa(userId), strconv.Itoa(articleId), utils.RandWords())
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
	fmt.Printf("tables created successfully!\n ")
}

func UserFollowTest() {
	_ = AddUserFollowing("1", "2")
	_ = AddUserFollowing("103", "5")
	_ = AddUserFollowing("103", "6")
	_ = AddUserFollowing("1", "3")
	_ = AddUserFollowing("103", "8")
	_ = AddUserFollowing("1", "12")
	_ = AddUserFollowing("1", "56")
	_ = AddUserFollowing("2", "4")
	_ = AddUserFollowing("2", "66")
	_ = AddUserFollowing("2", "23")
	_ = AddUserFollowing("2", "103")
	_ = AddUserFollowing("3", "8")
	_ = AddUserFollowing("3", "103")
	_ = AddUserFollowing("3", "78")
	_ = AddUserFollowing("3", "4")
	_ = AddUserFollowing("3", "1")
	_ = AddUserFollowing("3", "9")
	_ = AddUserFollowing("3", "7")
	_ = AddUserFollowing("3", "2")
	_ = AddUserFollowing("7", "103")
	_ = AddUserFollowing("7", "2")
	_ = AddUserFollowing("8", "103")
	fmt.Println("user relation pairs created successfully!")
}

func GetUserFollowerListTest() {
	users, total, err := GetUserFollowerList("2", "0", "1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("total:%d\n", total)
	for _, user := range users {
		fmt.Println(user.ID)
	}
}

func GetUserFollowingListTest() {
	users, total, err := GetUserFollowingList("2", "1", "2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("total:%d\n", total)
	for _, user := range users {
		fmt.Println(user.ID)
	}
}

func GetArticleCommentListTest() {
	articleComments, err := GetArticleCommentListFromArticleId("2", "2", "6")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, comment := range articleComments {
		fmt.Printf("%d %d %d %s\n", comment.ID, comment.ArticleID, comment.UserID, comment.Content)
	}
}

func GetSearchArticleListTest() {
	searchArticles, err := GetSearchArticleList("糖尿病", "1", "1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, comment := range searchArticles {
		fmt.Printf("%d %s\n", comment.ID, comment.Content)
	}
}

func AddCollectArticleTest() {
	_ = AddUserCollectedArticle("3", "4")
	_ = AddUserCollectedArticle("3", "2")
	_ = AddUserCollectedArticle("3", "6")
	_ = AddUserCollectedArticle("3", "7")
	_ = AddUserCollectedArticle("3", "9")
	_ = AddUserCollectedArticle("3", "2")
	_ = AddUserCollectedArticle("1", "3")
	_ = AddUserCollectedArticle("1", "2")
	_ = AddUserCollectedArticle("1", "5")
}

func RemoveCollectArticleTest() {
	_ = RemoveUserCollectedArticle("3", "4")
}

func GetUserCollectedArticleListTest() {
	articles, err := GetUserCollectedArticleList("3", "2", "3")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, article := range articles {
		fmt.Println(article.ID)
	}
}

func InitTopicTest() {
	for i := 1; i <= 500; i++ {
		userId := rand.Intn(100) + 1
		_ = AddTopic(strconv.Itoa(userId), utils.RandWords())
		fmt.Printf("topic.%d created successfully!\n", i)
	}
}

func InitTopicLordReplyTest() {
	for i := 1; i <= 500; i++ {
		userId := rand.Intn(100) + 1
		topicId := rand.Intn(100) + 1
		_ = AddTopicLordReply(strconv.Itoa(userId), strconv.Itoa(topicId), utils.RandWords())
		fmt.Printf("topicLordReply.%d created successfully!\n", i)
	}
}

func InitTopicLayerReplyTest() {
	for i := 1; i <= 500; i++ {
		userId := rand.Intn(100) + 1
		topicLordReplyId := rand.Intn(100) + 1
		_ = AddTopicLayerReply(strconv.Itoa(userId), strconv.Itoa(topicLordReplyId), utils.RandWords())
		fmt.Printf("topicLordReply.%d created successfully!\n", i)
	}
}

func GetUserReplyListTest() {
	replies, count, err := GetUserReplyList("33", "0", "4")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(count)
	for _, reply := range replies {
		fmt.Println(reply)
	}
}

func Test() {
	var users []User
	mysqlDb.Preload("FollowingUsers").Find(&users)
	for _, user := range users {
		fmt.Println(user.ID)
	}
}
