package authcontroller

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	dbmodels "github.com/niteshsiingh/doc-server/src/database/db-models"
	"github.com/niteshsiingh/doc-server/src/database/tables/databases"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
)

func (ac *AuthController) RefreshToken(c *gin.Context) {
	cxt := context.Background()
	var refreshTokenRequest dbmodels.RefreshToken

	err := c.ShouldBindJSON(&refreshTokenRequest)
	if err != nil {
		responses.NewResponse("Invalid request", http.StatusBadRequest).Send(c)
		return
	}

	refreshToken := refreshTokenRequest.Token
	isTokenActive, err := services.GetIsTokenActive(cxt, refreshToken, ac.DB)
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
	userRolesBytes, err := json.Marshal(user.Roles)
	if err != nil {
		responses.NewResponse("Internal server error", http.StatusInternalServerError).Send(c)
		return
	}
	userEntities := databases.User{
		Email:     pgtype.Text{String: user.EmailID, Valid: true},
		UserRoles: userRolesBytes,
	}
	authResponse, err := services.GenerateAuthResponse(cxt, &userEntities, ac.DB)
	if err != nil {
		responses.NewResponse("Internal server error", http.StatusInternalServerError).Send(c)
		return
	}

	responses.NewResponse(authResponse, http.StatusOK).Send(c)
}
