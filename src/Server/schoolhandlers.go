package server

import (
	"db"
	"fmt"
	"github.com/gin-gonic/gin"
)

//WARNING
func getArticle(SessionId string, ArticleId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	article, err := db.GetArticleFromArticleId(ArticleId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOKWithData(gin.H{
		"title":      article.Title,
		"contentUrl": fmt.Sprintf("/article/%d", article.ID),
	})
}

func getArticlePage(ArticleId string) (string, error) {
	article, err := db.GetArticleFromArticleId(ArticleId)
	if err != nil {
		return "", err
	}
	return article.Content, err
}

func getArticleList(SessionId string, BeginId string, NeedNumber string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if userId == "" {
		return responseNormalError("请先登录")
	}
	articles, err := db.GetArticleList(BeginId, NeedNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respArticles := make([]gin.H, len(articles))
	for i := 0; i < len(articles); i++ {
		respArticles[i] = gin.H{
			"articleId":   articles[i].ID,
			"title":       articles[i].Title,
			"content":     articles[i].Content[:40],
			"articleTime": articles[i].CreatedAt,
			"imgUrl":      articles[i].CoverImageUrl,
			"views":       articles[i].ReadCount,
		}
	}
	return responseOKWithData(gin.H{
		"data": respArticles,
	})
}

func createArticleComment(SessionId string, ArticleId string, Content string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	err = db.AddArticleComment(userId, ArticleId, Content)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func getArticleCommentList(SessionId string, ArticleId string, BeginId string, NeedNumber string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	comments, err := db.GetArticleCommentListFromArticleId(ArticleId, BeginId, NeedNumber)
	respComments := make([]gin.H, len(comments))
	for i := 0; i < len(comments); i++ {
		fmt.Println(comments[i].User)
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

func getSearchArticleList(SessionId string, SearchContent string, BeginId string, NeedNumber string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	if SearchContent == "" {
		return responseNormalError("关键词不能为空")
	}
	articles, err := db.GetSearchArticleList(SearchContent, BeginId, NeedNumber)
	respArticles := make([]gin.H, len(articles))
	for i := 0; i < len(articles); i++ {
		respArticles[i] = gin.H{
			"articleId":   articles[i].ID,
			"title":       articles[i].Title,
			"content":     articles[i].Content[:40],
			"articleTime": articles[i].CreatedAt,
			"imgUrl":      articles[i].CoverImageUrl,
			"views":       articles[i].ReadCount,
		}
	}
	return responseOKWithData(respArticles)
}

func getUserCollectedArticleList(SessionId string, BeginId string, NeedNumber string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	articles, err := db.GetUserCollectedArticleList(userId, BeginId, NeedNumber)
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
			"content":   articles[i].Content[:40],
		}
	}
	return responseOKWithData(gin.H{
		"total": count,
		"data":  respArticles,
	})
}

func getUserArticleCommentList(SessionId string, BeginId string, NeedNumber string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	comments, err := db.GetArticleCommentListFromUserId(userId, BeginId, NeedNumber)
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

func addCollectedArticle(SessionId string, ArticleId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	err = db.AddUserCollectedArticle(userId, ArticleId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func removeCollectedArticle(SessionId string, ArticleId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if userId == "" {
		return responseNormalError("请先登录")
	}
	err = db.RemoveUserCollectedArticle(userId, ArticleId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}
