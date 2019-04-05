package db

import (
	"errors"
	"fmt"
	"math/rand"
)

func AddNewArticle(title string, content string, plainContent string) error {
	imageNum := rand.Intn(155) + 1
	coverImageUrl := fmt.Sprintf("/static/articleImg/%d.jpg", imageNum)
	article := Article{
		Title:         title,
		Content:       content,
		PlainContent:  plainContent,
		CoverImageUrl: coverImageUrl,
	}
	mysqlDb.Create(&article)
	mysqlDb.Save(&article)
	return nil
}

func GetArticleFromArticleId(articleId int) (Article, error) {
	var articleTmp Article
	if mysqlDb.First(&articleTmp, articleId).RecordNotFound() {
		return articleTmp, errors.New("this article has not existed")
	}
	return articleTmp, nil
}

func GetArticleList(beginId int, needNumber int) ([]Article, error) {
	var articles []Article
	err := mysqlDb.Offset(beginId).Limit(needNumber).Find(&articles).Error
	return articles, err
}

func AddArticleComment(userId int, articleId int, content string) error {
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return err
	}
	articleComment := ArticleComment{
		Content: content,
		User:    user,
	}
	var articleTmp Article
	mysqlDb.First(&articleTmp, articleId)
	err = mysqlDb.Model(&articleTmp).Association("ArticleComments").Append(articleComment).Error
	return err
}

func GetArticleCommentListFromArticleId(articleId int, beginId int, needNumber int) ([]ArticleComment, error) {
	var comments []ArticleComment
	err := mysqlDb.Preload("User").Where(&ArticleComment{ArticleID: articleId}).
		Offset(beginId).Limit(needNumber).Find(&comments).Error
	return comments, err
}

func GetArticleCommentListFromUserId(userId int, beginId int, needNumber int) ([]ArticleComment, error) {
	var comments []ArticleComment
	err := mysqlDb.Preload("Article").Where(&ArticleComment{UserID: userId}).
		Offset(beginId).Limit(needNumber).Find(&comments).Error
	return comments, err
}

func GetSearchArticleList(searchContent string, beginId int, needNumber int) ([]Article, error) {
	var articles []Article
	err := mysqlDb.Where("title LIKE ?",
		fmt.Sprintf("%%%s%%", searchContent)).Offset(beginId).Limit(needNumber).Find(&articles).Error
	return articles, err
}

func AddUserCollectedArticle(userId int, articleId int) error {
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return err
	}
	article, err := GetArticleFromArticleId(articleId)
	if err != nil {
		return err
	}
	err = mysqlDb.Model(&user).Association("CollectedArticles").Append(article).Error
	return err
}

func RemoveUserCollectedArticle(userId int, articleId int) error {
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return err
	}
	article, err := GetArticleFromArticleId(articleId)
	if err != nil {
		return err
	}
	err = mysqlDb.Model(&user).Association("CollectedArticles").Delete(article).Error
	return err
}

func GetUserCollectedArticleList(userId int, beginId int, needNumber int) ([]Article, error) {
	var articles []Article
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return articles, err
	}
	err = mysqlDb.Model(&user).Offset(beginId).Limit(needNumber).
		Related(&articles, "CollectedArticles").Error
	return articles, err
}

func GetUserCollectedArticleCount(userId int) (int, error) {
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return 0, err
	}
	count := mysqlDb.Model(&user).Association("CollectedArticles").Count()
	return count, nil
}

func GetUserArticleCommentCount(userId int) (int, error) {
	var count int
	err := mysqlDb.Model(&ArticleComment{}).Where(&ArticleComment{UserID: userId}).Count(&count).Error
	return count, err
}

func GetArticleCommentCount(articleId int) (int, error) {
	var count int
	err := mysqlDb.Model(&ArticleComment{}).Where(&ArticleComment{ArticleID: articleId}).Count(&count).Error
	return count, err
}

func CheckUserCollectedArticle(userId int, articleId int) (bool, error) {
	var count int
	err := mysqlDb.Table("user_collected_article").
		Where("user_id=? and article_id=?", userId, articleId).Count(&count).Error
	return count > 0, err
}

func getArticleCommentFromArticleCommentId(articleCommentId int) (ArticleComment, error) {
	var articleCommentTmp ArticleComment
	if mysqlDb.First(&articleCommentTmp, articleCommentId).RecordNotFound() {
		return articleCommentTmp, errors.New("this article comment has not existed")
	}
	return articleCommentTmp, nil
}

func ValueArticleComment(articleCommentId int, value int) error {
	articleComment, err := getArticleCommentFromArticleCommentId(articleCommentId)
	if err != nil {
		return err
	}
	articleComment.ThumbsUpCount += value
	return mysqlDb.Save(&articleComment).Error
}

func AddArticleReadCount(articleId int) error {
	article, err := GetArticleFromArticleId(articleId)
	if err != nil {
		return err
	}
	article.ReadCount++
	return mysqlDb.Save(&article).Error
}
