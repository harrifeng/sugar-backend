package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type responseSimpleUser struct {
	followId uint
	username string
	iconUrl  string
}

type responseBody struct {
	Status int
	Data   interface{}
}

func responseInternalServerError(err error) responseBody {
	return responseBody{
		Status: http.StatusInternalServerError,
		Data: gin.H{
			"error": err.Error(),
		},
	}
}

func responseNormalError(errorMessage string) responseBody {
	return responseBody{
		Status: http.StatusOK,
		Data: gin.H{
			"code": 1,
			"msg":  errorMessage,
		},
	}
}

func responseOKWithData(data interface{}) responseBody {
	dataMap := data.(gin.H)
	dataMap["code"] = 0
	return responseBody{
		Status: http.StatusOK,
		Data:   dataMap,
	}
}

func responseOK() responseBody {
	return responseBody{
		Status: http.StatusOK,
		Data: gin.H{
			"code": 0,
		},
	}
}
