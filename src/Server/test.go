package server

import (
	"db"
	"github.com/gin-gonic/gin"
)

const (
	lmy789789 = "0c152176187ce61c9614c072e1a1e2f7"
	lmy456456 = "6a914208cfee6f3862144e5b40a83454"
)

func initDb(c *gin.Context) {
	db.AutoCreateTableTest()
	resp := responseOK()
	c.JSON(resp.Status, resp.Data)
}

func initArticle(c *gin.Context) {
	db.InitAllArticle()
	resp := responseOK()
	c.JSON(resp.Status, resp.Data)
}

func createNewTestUser(c *gin.Context) {
	err := db.CreateNewUser("11111111111", "0c152176187ce61c9614c072e1a1e2f7")
	var resp responseBody
	if err != nil {
		resp = responseNormalError("用户已存在")
	} else {
		resp = responseOK()
	}
	c.JSON(resp.Status, resp.Data)
}

func DatabaseTest() {
	//db.AutoCreateTableTest()
	//db.InitAllArticle()
	//db.InitUser()
	//db.IninArticleComment()
	//db.GetArticleCommentListTest()
	//db.GetSearchArticleListTest()
	//db.AddCollectArticleTest()
	db.RemoveCollectArticleTest()
}
