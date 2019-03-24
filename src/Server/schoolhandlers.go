package server

import (
	"db"
)

func getArticle(SessionId string, ArticleId string) responseBody {
	Article, err := db.GetArticleFromArticleId(ArticleId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOKWithHtml(Article.Content)
}

//func getArticleList(SessionId string,BeginId string,NeedNumber string) responseBody{
//articles,err:=db.GetArticleList(SessionId,BeginId,NeedNumber)
//if err!=nil{
//	return responseInternalServerError(err)
//}

//}
