package responses

import (
	"time"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success    bool        `json:"success"`
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	ServerTime string      `json:"serverTime"`
	Data       interface{} `json:"data,omitempty"`
}

type SuccessResponseInterface interface {
	Send(ctx *gin.Context)
}

func NewResponse(data interface{}, code int) (r SuccessResponse) {
	r = SuccessResponse{
		Success:    true,
		Message:    "Success",
		ServerTime: time.Now().Format(time.RFC3339),
		Code:       code,
		Data:       data,
	}
	return r
}

func (r SuccessResponse) Send(c *gin.Context) {
	c.JSON(r.Code, r)
}
