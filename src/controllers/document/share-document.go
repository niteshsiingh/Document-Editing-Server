package document

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	dbmodels "github.com/niteshsiingh/doc-server/src/database/db-models"
	"github.com/niteshsiingh/doc-server/src/database/tables/databases"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
	"gorm.io/gorm"
)

func (dc *DocumentController) ShareDocument(ctx *gin.Context) {
	cxt := context.Background()
	documentIDstr := ctx.Param("document_id")
	documentID, err := strconv.ParseUint(documentIDstr, 10, 64)
	if err != nil {
		responses.NewResponse("Invalid document ID", 400).Send(ctx)
		return
	}
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
	// var document databases.Document
	_, err = dc.DB.GetDocumentById(cxt, databases.GetDocumentByIdParams{
		ID:     int32(documentID),
		UserID: pgtype.Int4{Int32: int32(userID), Valid: true},
	})
	// err = dc.DB.First(&document, documentID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.NewResponse("Document not found", 404).Send(ctx)
			return
		}
		responses.NewResponse("Failed to retrieve document", 500).Send(ctx)
		return
	}
	var requestBody struct {
		Email      string                  `json:"email"`
		Permission dbmodels.PermissionEnum `json:"permission"`
	}
	err = ctx.ShouldBindJSON(&requestBody)
	if err != nil {
		responses.NewResponse("Invalid json request", 400).Send(ctx)
		return
	}

	var sharedUser databases.User
	sharedUser, err = dc.DB.GetUserByEmail(cxt, pgtype.Text{String: requestBody.Email, Valid: true})
	if err != nil {
		responses.NewResponse("User not found", http.StatusNotFound).Send(ctx)
		return
	}

	// documentUser := dbmodels.DocumentUser{
	// 	DocumentID: uint(documentID),
	// 	UserID:     sharedUser.ID,
	// 	Permission: requestBody.Permission,
	// }
	err = dc.DB.CreateDocumentUser(cxt, databases.CreateDocumentUserParams{
		UserID:     pgtype.Int4{Int32: int32(sharedUser.ID), Valid: true},
		DocumentID: pgtype.Int4{Int32: int32(documentID), Valid: true},
		Read:       pgtype.Bool{Bool: true, Valid: true},
		Write:      pgtype.Bool{Bool: true, Valid: true},
		Share:      pgtype.Bool{Bool: true, Valid: true},
		Download:   pgtype.Bool{Bool: true, Valid: true},
		Admin:      pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		fmt.Println(err)
		responses.NewResponse("Failed to share document", 500).Send(ctx)
		return
	}
	user, err := services.FindUserByID(cxt, uint(userID), dc.DB)
	if err != nil {
		responses.NewResponse("User not found", 404).Send(ctx)
		return
	}
	mail := services.MailOptions{
		From:    os.Getenv("BREVO_SENDER_EMAIL"),
		To:      []string{sharedUser.Email.String},
		Subject: user.Email.String + " shared a document with you!",
		Body:    "Click the following link to view and edit the document: " + os.Getenv("FRONT_END_URL") + "/document/" + strconv.Itoa(int(documentID)),
	}
	if err := dc.MS.SendMail(mail); err != nil {
		fmt.Println(err, 2)
		responses.NewResponse("Failed to send email", 500).Send(ctx)
		return
	}
	responses.NewResponse("Document shared successfully", 200).Send(ctx)
}

func (dc *DocumentController) RemoveSharedUser(ctx *gin.Context) {
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
	documentIDstr := ctx.Param("document_id")
	documentID, err := strconv.ParseUint(documentIDstr, 10, 64)
	if err != nil {
		responses.NewResponse("Invalid document ID", 400).Send(ctx)
		return
	}
	// var document databases.Document
	_, err = dc.DB.FindDocument(cxt, databases.FindDocumentParams{
		ID:     int32(documentID),
		UserID: pgtype.Int4{Int32: int32(userID), Valid: true},
	})
	if err != nil {
		responses.NewResponse("Document not found or unauthorized", http.StatusNotFound).Send(ctx)
		return
	}
	// var documentUser dbmodels.DocumentUser
	documentUser, err := dc.DB.FindDocumentUser(cxt, databases.FindDocumentUserParams{
		DocumentID: pgtype.Int4{Int32: int32(documentID), Valid: true},
		UserID:     pgtype.Int4{Int32: int32(userID), Valid: true},
	})
	if err != nil {
		responses.NewResponse("DocumentUser association not found", http.StatusNotFound).Send(ctx)
		return
	}
	err = dc.DB.DeleteDocumentUser(cxt, documentUser.ID)
	if err != nil {
		responses.NewResponse("Failed to delete document user association", http.StatusInternalServerError).Send(ctx)
		return
	}
	responses.NewResponse("User removed successfully", 200).Send(ctx)
}
