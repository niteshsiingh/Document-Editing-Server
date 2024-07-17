package authcontroller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	dbmodels "github.com/niteshsiingh/doc-server/src/database/db-models"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
)

func (ac *AuthController) RefreshToken(c *gin.Context) {
	var refreshTokenRequest dbmodels.RefreshToken
	err := c.ShouldBindJSON(&refreshTokenRequest)
	if err != nil {
		responses.NewResponse("Invalid request", http.StatusBadRequest).Send(c)
		return
	}

	refreshToken := refreshTokenRequest.Token
	isTokenActive, err := services.GetIsTokenActive(refreshToken, ac.DB)
	if err != nil {
		responses.NewResponse("Internal server error", http.StatusInternalServerError).Send(c)
		return
	}
	if !isTokenActive {
		responses.NewResponse("Invalid token", http.StatusForbidden).Send(c)
		return
	}
	user, err := middleware.GetAuth().ParseAuth(refreshToken, os.Getenv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		responses.NewResponse("Internal server error", http.StatusInternalServerError).Send(c)
		return
	}

	userEntities := dbmodels.User{
		Email:     user.EmailID,
		UserRoles: user.Roles,
	}
	authResponse, err := services.GenerateAuthResponse(&userEntities, ac.DB)
	if err != nil {
		responses.NewResponse("Internal server error", http.StatusInternalServerError).Send(c)
		return
	}

	responses.NewResponse(authResponse, http.StatusOK).Send(c)
}
