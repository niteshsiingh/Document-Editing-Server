package authcontroller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	dbmodels "github.com/niteshsiingh/doc-server/src/database/db-models"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
)

func (ac *AuthController) Logout(ctx *gin.Context) {
	var user dbmodels.User
	cxt := context.Background()
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		responses.NewResponse("Invalid request", http.StatusBadRequest).Send(ctx)
		return
	}
	err = services.LogoutUser(cxt, user.ID, ac.DB)
	if err != nil {
		responses.NewResponse("Internal server error", http.StatusInternalServerError).Send(ctx)
		return
	}
	responses.NewResponse("Successfully logged out", http.StatusOK).Send(ctx)
}
