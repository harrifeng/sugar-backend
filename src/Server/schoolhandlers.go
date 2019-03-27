package server

import (
	"db"
	"github.com/gin-gonic/gin"
)

func getArticle(SessionId string, ArticleId string) responseBody {
	//Article, err := db.GetArticleFromArticleId(ArticleId)
	//if err != nil {
	//	return responseInternalServerError(err)
	//}
	return responseBody{}
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
			"contentUrl":  articles[i].Content,
			"articleTime": articles[i].CreatedAt,
			"imgUrl":      articles[i].CoverImageUrl,
			"views":       articles[i].ReadCount,
		}
	}
	return responseOKWithData(gin.H{
		"data": respArticles,
	})
}
