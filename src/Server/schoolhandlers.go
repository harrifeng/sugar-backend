package server

import (
	"db"
	"fmt"
	"github.com/gin-gonic/gin"
	"utils"
)

// 获取单个文章信息
func getArticle(userId int, articleId int) responseBody {
	article, err := db.GetArticleFromArticleId(articleId)
	if err != nil {
		return responseInternalServerError(err)
	}
	collected, err := db.CheckUserCollectedArticle(userId, articleId)
	if err != nil {
		return responseInternalServerError(err)
	}
	count, err := db.GetArticleCommentCount(articleId)
	if err != nil {
		return responseInternalServerError(err)
	}
	err = db.AddArticleReadCount(articleId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOKWithData(gin.H{
		"contentUrl": fmt.Sprintf("/school/article-page/%d", article.ID),
		"collected":  collected,
		"comNumber":  count,
	})
}

// 获取文章HTML页面
func getArticlePage(articleId int) (string, error) {
	article, err := db.GetArticleFromArticleId(articleId)
	if err != nil {
		return "", err
	}
	return article.Content, err
}

// 获取文章列表
func getArticleList(beginId int, needNumber int) responseBody {
	articles, err := db.GetArticleList(beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respArticles := make([]gin.H, len(articles))
	for i := 0; i < len(articles); i++ {
		respArticles[i] = gin.H{
			"articleId":   articles[i].ID,
			"title":       articles[i].Title,
			"content":     utils.StringCut(articles[i].Content, 40),
			"articleTime": articles[i].CreatedAt,
			"imgUrl":      articles[i].CoverImageUrl,
			"views":       articles[i].ReadCount,
		}
	}
	return responseOKWithData(respArticles)
}

// 创建文章评论
func createArticleComment(userId int, articleId int, content string) responseBody {
	err := db.AddArticleComment(userId, articleId, content)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

// 获取文章的评论列表
func getArticleCommentList(articleId int, beginId int, needNumber int) responseBody {
	comments, err := db.GetArticleCommentListFromArticleId(articleId, beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respComments := make([]gin.H, len(comments))
	for i := 0; i < len(comments); i++ {
		respComments[i] = gin.H{
			"commentId":   comments[i].ID,
			"content":     comments[i].Content,
			"commentTime": comments[i].CreatedAt,
			"userId":      comments[i].UserID,
			"likes":       comments[i].ThumbsUpCount,
			"username":    comments[i].User.UserName,
			"iconUrl":     comments[i].User.HeadPortraitUrl,
		}
	}
	return responseOKWithData(respComments)
}

func getSearchArticleList(searchContent string, beginId int, needNumber int) responseBody {
	if searchContent == "" {
		return responseNormalError("关键词不能为空")
	}
	articles, err := db.GetSearchArticleList(searchContent, beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respArticles := make([]gin.H, len(articles))
	for i := 0; i < len(articles); i++ {
		respArticles[i] = gin.H{
			"articleId":   articles[i].ID,
			"title":       articles[i].Title,
			"content":     utils.StringCut(articles[i].Content, 40),
			"articleTime": articles[i].CreatedAt,
			"imgUrl":      articles[i].CoverImageUrl,
			"views":       articles[i].ReadCount,
		}
	}
	return responseOKWithData(respArticles)
}

func getUserCollectedArticleList(userId int, beginId int, needNumber int) responseBody {
	articles, err := db.GetUserCollectedArticleList(userId, beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	count, err := db.GetUserCollectedArticleCount(userId)
	if err != nil {
		return responseInternalServerError(err)
	}
	respArticles := make([]gin.H, len(articles))
	for i := 0; i < len(articles); i++ {
		respArticles[i] = gin.H{
			"articleId": articles[i].ID,
			"title":     articles[i].Title,
			"content":   utils.StringCut(articles[i].Content, 40),
		}
	}
	return responseOKWithData(gin.H{
		"total": count,
		"data":  respArticles,
	})
}

func getUserArticleCommentList(userId int, beginId int, needNumber int) responseBody {
	comments, err := db.GetArticleCommentListFromUserId(userId, beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	count, err := db.GetUserArticleCommentCount(userId)
	if err != nil {
		return responseInternalServerError(err)
	}
	respComments := make([]gin.H, len(comments))
	for i := 0; i < len(comments); i++ {
		respComments[i] = gin.H{
			"commentId":   comments[i].ID,
			"content":     comments[i].Content,
			"title":       comments[i].Article.Title,
			"articleId":   comments[i].ArticleID,
			"likes":       comments[i].ThumbsUpCount,
			"commentTime": comments[i].CreatedAt,
		}
	}
	return responseOKWithData(gin.H{
		"total": count,
		"data":  respComments,
	})
}

func addCollectedArticle(userId int, articleId int) responseBody {
	err := db.AddUserCollectedArticle(userId, articleId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func removeCollectedArticle(userId int, articleId int) responseBody {
	err := db.RemoveUserCollectedArticle(userId, articleId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func valueArticleComment(articleCommentId int, value int) responseBody {
	err := db.ValueArticleComment(articleCommentId, value)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}
