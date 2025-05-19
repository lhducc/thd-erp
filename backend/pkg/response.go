package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"statuscode"`
}

func (e *Response) Error() string {
	return fmt.Sprintf("[%d] %s", e.StatusCode, e.Message)
}

//	func ErrorResponseResult(c *gin.Context, message string, statusCode int, data interface{}) {
//		c.JSON(statusCode, Response{
//			Message:    message,
//			StatusCode: statusCode,
//		})
//	}
func ResponseMessage(c *gin.Context, message string, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	})
}
