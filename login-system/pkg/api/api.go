package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindRequest[T any](ctx *gin.Context) (httpCode int, request T, response Response) {
	var (
		err error
	)
	httpCode = http.StatusOK
	defer func() {
		if err != nil {
			httpCode = http.StatusBadRequest
			response = Response{
				Status: StatusBadRequest,
				Error: &ErrorItem{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			}
			return
		}
	}()

	if ctx.Request.Method == http.MethodGet {
		err = ctx.ShouldBindQuery(&request)
	} else {
		err = ctx.ShouldBindJSON(&request)
	}
	return
}

func JSON(c *gin.Context, httpCode int, response Response) {
	c.JSON(httpCode, response)
}

type Response struct {
	Status  ResponseStatus `json:"status"`
	Message *string        `json:"message,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
	Error   *ErrorItem     `json:"error,omitempty"`
}
type ErrorItem struct {
	Code    int    `json:"code"`    // Error Code
	Message string `json:"message"` // Error Message
}

type ResponseStatus string

const (
	StatusSuccess      ResponseStatus = "SUCCESS"
	StatusError        ResponseStatus = "ERROR"
	StatusUnauthorized ResponseStatus = "UNAUTHORIZED"
	StatusNotFound     ResponseStatus = "NOT_FOUND"
	StatusBadRequest   ResponseStatus = "BAD_REQUEST"
)
