package document

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	dbmodels "github.com/niteshsiingh/doc-server/src/database/db-models"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
	"gorm.io/gorm"
)

type DocumentController struct {
	DB *gorm.DB
	MS *services.MailService
}

func NewDocumentController(db *gorm.DB) *DocumentController {
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
		fmt.Println("err: ", err)
		responses.NewResponse("Invalid document ID", 400).Send(ctx)
		return
	}
	document, err := services.FindDocumentByID(uint(documentID), uint(userID), dc.DB)
	if err != nil {
		responses.NewResponse("Document not found", 404).Send(ctx)
		return
	}
	responses.NewResponse(document, 200).Send(ctx)
}

func (dc *DocumentController) GetAllDocuments(ctx *gin.Context) {
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
	var documents []dbmodels.Document
	dc.DB.Where("user_id = ?", userID).Find(&documents)

	var documentUsers []dbmodels.DocumentUser
	dc.DB.Where("user_id = ?", userID).Preload("Document").Find(&documentUsers)

	for _, documentUser := range documentUsers {
		documents = append(documents, documentUser.Document)
	}
	responses.NewResponse(documents, 200).Send(ctx)
}
