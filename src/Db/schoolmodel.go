package db

import (
	"errors"
	"fmt"
	"strconv"
)

func AddNewArticle(Title string, Content string) error {
	article := Article{
		Title:   Title,
		Content: Content,
	}
	mysqlDb.Create(&article)
	mysqlDb.Save(&article)
	return nil
}

func GetArticleFromArticleId(ArticleId string) (Article, error) {
	var articleTmp Article
	articleId, _ := strconv.Atoi(ArticleId)
	if mysqlDb.First(&articleTmp, articleId).RecordNotFound() {
		return articleTmp, errors.New("this article has not existed")
	}
	return articleTmp, nil
}

func GetArticleList(BeginId string, NeedNumber string) ([]Article, error) {
	var articles []Article
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	mysqlDb.Offset(beginId).Limit(needNumber).Find(&articles)
	return articles, nil
}

func AddArticleComment(UserId string, ArticleId string, Content string) error {
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return err
	}
	articleComment := ArticleComment{
		Content: Content,
		User:    user,
	}
	var articleTmp Article
	articleId, _ := strconv.Atoi(ArticleId)
	mysqlDb.First(&articleTmp, articleId)
	mysqlDb.Model(&articleTmp).Association("ArticleComments").Append(articleComment)
	return nil
}

func GetArticleCommentListFromArticleId(ArticleId string, BeginId string, NeedNumber string) ([]ArticleComment, error) {
	var comments []ArticleComment
	articleId, _ := strconv.Atoi(ArticleId)
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	mysqlDb.Where(&ArticleComment{ArticleID: articleId}).Offset(beginId).Limit(needNumber).Find(&comments)
	return comments, nil
}

func GetArticleCommentListFromUserId(UserId string, BeginId string, NeedNumber string) ([]ArticleComment, error) {
	var comments []ArticleComment
	userId, _ := strconv.Atoi(UserId)
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	mysqlDb.Where(&ArticleComment{UserID: userId}).Offset(beginId).Limit(needNumber).Find(&comments)
	return comments, nil
}

func GetSearchArticleList(SearchContent string, BeginId string, NeedNumber string) ([]Article, error) {
	var articles []Article
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	mysqlDb.Where("title LIKE ?",
		fmt.Sprintf("%%%s%%", SearchContent)).Offset(beginId).Limit(needNumber).Find(&articles)
	return articles, nil
}

func AddUserCollectedArticle(UserId string, ArticleId string) error {
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return err
	}
	article, err := GetArticleFromArticleId(ArticleId)
	if err != nil {
		return err
	}
	mysqlDb.Model(&user).Association("CollectedArticles").Append(article)
	return nil
}

func RemoveUserCollectedArticle(UserId string, ArticleId string) error {
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return err
	}
	article, err := GetArticleFromArticleId(ArticleId)
	if err != nil {
		return err
	}
	mysqlDb.Model(&user).Association("CollectedArticles").Delete(article)
	return nil
}

// WARNING
func GetUserCollectedArticleList(UserId string, BeginId string, NeedNumber string) ([]Article, error) {
	var articles []Article
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return articles, err
	}
	mysqlDb.Model(&user).Association("CollectedArticles").Find(&articles)
	return articles, nil
}
