package server

import (
	"db"
	"github.com/gin-gonic/gin"
)

func LoginAuth()gin.HandlerFunc{
	return func(c *gin.Context) {
		var sessionId string
		if c.Request.Method == "GET"{
			sessionId = c.Query("session_id")
		}else{
			sessionId = c.PostForm("session_id")
		}
		var resp responseBody
		if sessionId == ""{
			resp = responseNormalError("请先登录")
		}else{
			userId,err:= db.GetNowSessionId(sessionId)
			if err!=nil{
				resp = responseInternalServerError(err)
			}else if userId == ""{
				resp = responseNormalError("请先登录")
			}else{
				c.Set("user_id",userId)
				c.Next()
				return
			}
		}
		c.JSON(resp.Status,resp.Data)
	}
}