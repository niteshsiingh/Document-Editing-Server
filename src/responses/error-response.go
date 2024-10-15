package responses

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Success    bool        `json:"success"`
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	ServerTime string      `json:"serverTime"`
	Data       interface{} `json:"data,omitempty"`
}

type ErrorResponseInterface interface {
	Throw(ctx *gin.Context)
}

func NewErrorResponse(data interface{}, code int) (r ErrorResponse) {

	fmt.Println(data)

	message := "Some error occurred."

	if err, ok := data.(error); ok {
		message = err.Error()
	} else if err, ok := data.(string); ok {
		message = err
	}

	r = ErrorResponse{
		Success:    false,
		Message:    message,
		ServerTime: time.Now().Format(time.RFC3339),
		Code:       code,
		Data:       data,
	}
	return r
}

func (r ErrorResponse) Throw(c *gin.Context) {
	c.JSON(r.Code, r)
	c.Abort()
}
