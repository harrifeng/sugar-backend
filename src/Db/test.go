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
	for i := 0; i < 100; i++ {
		_ = CreateNewUser(utils.RandPhoneNumber(), "0c152176187ce61c9614c072e1a1e2f7")
		fmt.Printf("user.%d created successfully!\n ", i)
	}
}

func IninArticleComment() {
	for i := 0; i <= 1000; i++ {
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
	mysqlDb.AutoMigrate(&ArticleLabel{})
	mysqlDb.AutoMigrate(&Topic{})
	mysqlDb.AutoMigrate(&TopicLordReply{})
	mysqlDb.AutoMigrate(&TopicLayerReply{})
	mysqlDb.AutoMigrate(&UserPrivacySetting{})
	fmt.Printf("tables created successfully!\n ")
}

func UserFollowTest() {
	_ = AddUserFollowing("2", "3")
	_ = AddUserFollowing("3", "2")
	_ = AddUserFollowing("5", "2")
	_ = AddUserFollowing("7", "3")
	_ = AddUserFollowing("2", "87")
	_ = AddUserFollowing("8", "75")
	_ = AddUserFollowing("3", "23")
	_ = AddUserFollowing("23", "54")
}

func GetUserFollowerListTest() {
	users, err := GetUserFollowerList("2", "0", "10")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, user := range users {
		fmt.Println(user.ID)
	}
}

func GetUserFollowingListTest() {
	users, err := GetUserFollowingList("3", "0", "10")
	if err != nil {
		fmt.Println(err)
		return
	}
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
	err := AddUserCollectedArticle("3", "4")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func RemoveCollectArticleTest() {
	err := RemoveUserCollectedArticle("3", "4")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Test() {
	var users []User
	mysqlDb.Preload("FollowingUsers").Find(&users)
	for _, user := range users {
		fmt.Println(user.ID)
	}
}
