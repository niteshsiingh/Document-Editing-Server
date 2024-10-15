package document

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/niteshsiingh/doc-server/src/database/tables/databases"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
)

func (dc *DocumentController) CreateDocument(ctx *gin.Context) {
	cxt := context.Background()
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
	_, err = services.FindUserByID(cxt, userID, dc.DB)
	if err != nil {
		responses.NewResponse("User not found", http.StatusNotFound).Send(ctx)
		return
	}
	document := databases.Document{
		UserID:   pgtype.Int4{Int32: int32(userID), Valid: true},
		IsPublic: pgtype.Bool{Bool: false, Valid: true},
		// User:     user,
	}
	newDocID, err := dc.DB.CreateDocument(cxt, databases.CreateDocumentParams{
		UserID:   document.UserID,
		IsPublic: document.IsPublic,
	})
	if err != nil {
		responses.NewResponse("Internal server error", http.StatusInternalServerError).Send(ctx)
		return
	}
	type res struct {
		ID int32 `json:"id"`
	}
	respon := res{ID: newDocID}
	responses.NewResponse(respon, 200).Send(ctx)
}
