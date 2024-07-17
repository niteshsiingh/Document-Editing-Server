package document

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	dbmodels "github.com/niteshsiingh/doc-server/src/database/db-models"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
)

func (dc *DocumentController) CreateDocument(ctx *gin.Context) {
	authTokenHeader := ctx.GetHeader("Authorization")
	splitted := strings.Split(authTokenHeader, " ")
	if authTokenHeader == "" || len(splitted) != 2 {
		responses.NewResponse("Authentication token not found in the request header", http.StatusUnauthorized).Send(ctx)
		return
	}
	authToken := splitted[1]
	_, parsedToken, err := middleware.GetAuth().ParseToken(authToken)
	if err != nil {
		responses.NewResponse("Invalid token", 403).Send(ctx)
		return
	}
	userID := parsedToken.UserID
	user, err := services.FindUserByID(userID, dc.DB)
	if err != nil {
		responses.NewResponse("User not found", http.StatusNotFound).Send(ctx)
		return
	}
	document := dbmodels.Document{
		UserID:   userID,
		IsPublic: false,
		User:     user,
	}

	if err := dc.DB.Create(&document).Error; err != nil {
		responses.NewResponse("Internal server error", http.StatusInternalServerError).Send(ctx)
		return
	}
	type res struct {
		ID uint `json:"id"`
	}
	respon := res{ID: document.ID}
	responses.NewResponse(respon, 200).Send(ctx)
}
