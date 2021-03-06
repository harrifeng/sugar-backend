package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type responseBody struct {
	Status int
	Data   gin.H
}

func responseInternalServerError(err error) responseBody {
	println("error:", err.Error())
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
	return responseBody{
		Status: http.StatusOK,
		Data: gin.H{
			"code": 0,
			"data": data,
		},
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
