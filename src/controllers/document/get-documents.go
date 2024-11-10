package document

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/niteshsiingh/doc-server/src/database/tables/databases"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
)

type DocumentController struct {
	DB *databases.Queries
	MS *services.MailService
}

func NewDocumentController(db *databases.Queries) *DocumentController {
	ms, err := services.NewMailService()
	if err != nil {
		return nil
	}
	return &DocumentController{
		DB: db,
		MS: ms,
	}
}

func (dc *DocumentController) GetOneDocument(ctx *gin.Context) {
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
	documentIDStr := ctx.Param("document_id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		responses.NewResponse("Invalid document ID", 400).Send(ctx)
		return
	}
	document, err := services.FindDocumentByID(cxt, uint(documentID), uint(userID), dc.DB)
	if err != nil {
		responses.NewResponse("Document not found", 404).Send(ctx)
		return
	}
	responses.NewResponse(document, 200).Send(ctx)
}

func (dc *DocumentController) GetDocumentIdentifiers(ctx *gin.Context) {
	cxt := context.Background()
	authTokenHeader := ctx.GetHeader("Authorization")
	splitted := strings.Split(authTokenHeader, " ")
	if authTokenHeader == "" || len(splitted) != 2 {
		responses.NewResponse("Authentication token not found in the request header", http.StatusUnauthorized).Send(ctx)
		return
	}
	authToken := splitted[1]
	_, _, err := middleware.GetAuth().ParseToken(authToken)
	if err != nil {
		responses.NewResponse("Invalid token", 403).Send(ctx)
		return
	}
	// userID := parsedToken.UserID
	documentIDStr := ctx.Param("document_id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		responses.NewResponse("Invalid document ID", 400).Send(ctx)
		return
	}

	document, err := services.FindIdentifiersWithDocID(cxt, uint(documentID), dc.DB)
	if err != nil {
		responses.NewResponse("Document not found", 404).Send(ctx)
		return
	}

	responses.NewResponse(document, 200).Send(ctx)
}

func (dc *DocumentController) GetAllDocuments(ctx *gin.Context) {
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
	// userID, err := strconv.ParseUint(jwtUser.ID, 10, 64)
	// if err != nil {
	// 	responses.NewResponse("Invalid user ID", 400).Send(ctx)
	// 	return
	// }
	// userID := jwtUser.ID
	var documents []databases.Document
	documents, err = dc.DB.GetAllDocuments(cxt, pgtype.Int4{Int32: int32(userID), Valid: true})
	if err != nil {
		fmt.Println(err)
		responses.NewResponse("Failed to get documents", 500).Send(ctx)
		return
	}
	for _, document := range documents {
		fmt.Println(document)
	}
	responses.NewResponse(documents, 200).Send(ctx)
}
