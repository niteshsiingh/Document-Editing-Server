package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	err_response "github.com/niteshsiingh/doc-server/src/responses"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		tokenHeader := c.GetHeader("Authorization")
		splitted := strings.Split(tokenHeader, " ")
		if tokenHeader == "" || len(splitted) != 2 {
			err_response.NewErrorResponse("Authentication token not found in the request header", http.StatusUnauthorized).Throw(c)
			return
		}

		user, err := GetAuth().ParseAuth(splitted[1], os.Getenv("JWT_APP_KEY"))
		if err != nil {
			err_response.NewErrorResponse("Invalid authentication token", http.StatusUnauthorized).Throw(c)
			return
		}

		c.Set("uid", user.ID)
		c.Set("user", user)
		c.Next()
	}
}
