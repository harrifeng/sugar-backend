package db

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

func AddNewArticle(Title string, Content string,PlainContent string) error {
	imageNum := rand.Intn(155) + 1
	coverImageUrl := fmt.Sprintf("/static/articleImg/%d.jpg",imageNum)
	article := Article{
		Title:   Title,
		Content: Content,
		PlainContent:PlainContent,
		CoverImageUrl:coverImageUrl,
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
	err := mysqlDb.Offset(beginId).Limit(needNumber).Find(&articles).Error
	return articles, err
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
	err = mysqlDb.Model(&articleTmp).Association("ArticleComments").Append(articleComment).Error
	return err
}

func GetArticleCommentListFromArticleId(ArticleId string, BeginId string, NeedNumber string) ([]ArticleComment, error) {
	var comments []ArticleComment
	articleId, _ := strconv.Atoi(ArticleId)
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	err := mysqlDb.Preload("User").Where(&ArticleComment{ArticleID: articleId}).
		Offset(beginId).Limit(needNumber).Find(&comments).Error
	return comments, err
}

func GetArticleCommentListFromUserId(UserId string, BeginId string, NeedNumber string) ([]ArticleComment, error) {
	var comments []ArticleComment
	userId, _ := strconv.Atoi(UserId)
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	err := mysqlDb.Preload("Article").Where(&ArticleComment{UserID: userId}).Offset(beginId).Limit(needNumber).Find(&comments).Error
	return comments, err
}

func GetSearchArticleList(SearchContent string, BeginId string, NeedNumber string) ([]Article, error) {
	var articles []Article
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	err := mysqlDb.Where("title LIKE ?",
		fmt.Sprintf("%%%s%%", SearchContent)).Offset(beginId).Limit(needNumber).Find(&articles).Error
	return articles, err
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
	err = mysqlDb.Model(&user).Association("CollectedArticles").Append(article).Error
	return err
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
	err = mysqlDb.Model(&user).Association("CollectedArticles").Delete(article).Error
	return err
}

func GetUserCollectedArticleList(UserId string, BeginId string, NeedNumber string) ([]Article, error) {
	var articles []Article
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return articles, err
	}
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	err = mysqlDb.Model(&user).Offset(beginId).Limit(needNumber).
		Related(&articles, "CollectedArticles").Error
	return articles, err
}

func GetUserCollectedArticleCount(UserId string) (int, error) {
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return 0, err
	}
	count := mysqlDb.Model(&user).Association("CollectedArticles").Count()
	return count, nil
}

func GetUserArticleCommentCount(UserId string) (int, error) {
	userId, _ := strconv.Atoi(UserId)
	var count int
	err := mysqlDb.Model(&ArticleComment{}).Where(&ArticleComment{UserID: userId}).Count(&count).Error
	return count, err
}

func GetArticleCommentCount(ArticleId string) (int, error) {
	articleId, _ := strconv.Atoi(ArticleId)
	var count int
	err := mysqlDb.Model(&ArticleComment{}).Where(&ArticleComment{ArticleID: articleId}).Count(&count).Error
	return count, err
}

func CheckUserCollectedArticle(UserId string, ArticleId string) (bool, error) {
	var count int
	err := mysqlDb.Table("user_collected_article").
		Where("user_id=? and article_id=?", UserId, ArticleId).Count(&count).Error
	return count > 0, err
}

func getArticleCommentFromArticleCommentId(ArticleCommentId string) (ArticleComment, error) {
	var articleCommentTmp ArticleComment
	articleCommentId, _ := strconv.Atoi(ArticleCommentId)
	if mysqlDb.First(&articleCommentTmp, articleCommentId).RecordNotFound() {
		return articleCommentTmp, errors.New("this article comment has not existed")
	}
	return articleCommentTmp, nil
}

func ValueArticleComment(ArticleCommentId string, Value string) error {
	articleComment, err := getArticleCommentFromArticleCommentId(ArticleCommentId)
	if err != nil {
		return err
	}
	value, _ := strconv.Atoi(Value)
	articleComment.ThumbsUpCount += value
	return mysqlDb.Save(&articleComment).Error
}

func AddArticleReadCount(ArticleId string) error {
	article, err := GetArticleFromArticleId(ArticleId)
	if err != nil {
		return err
	}
	article.ReadCount++
	return mysqlDb.Save(&article).Error
}
