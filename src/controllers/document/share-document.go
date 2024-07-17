package document

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	dbmodels "github.com/niteshsiingh/doc-server/src/database/db-models"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
	"gorm.io/gorm"
)

func (dc *DocumentController) ShareDocument(ctx *gin.Context) {
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
	var document dbmodels.Document
	err = dc.DB.First(&document, documentID).Error
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

	var sharedUser dbmodels.User
	if err := dc.DB.Where("email = ?", requestBody.Email).First(&sharedUser).Error; err != nil {
		responses.NewResponse("User not found", http.StatusNotFound).Send(ctx)
		return
	}

	documentUser := dbmodels.DocumentUser{
		DocumentID: uint(documentID),
		UserID:     sharedUser.ID,
		Permission: requestBody.Permission,
	}
	if err := dc.DB.Create(&documentUser).Error; err != nil {
		responses.NewResponse("Failed to share document", 500).Send(ctx)
		return
	}
	user, err := services.FindUserByID(uint(userID), dc.DB)
	if err != nil {
		responses.NewResponse("User not found", 404).Send(ctx)
		return
	}
	mail := services.MailOptions{
		From:    "8826ns@gmail.com",
		To:      []string{sharedUser.Email},
		Subject: user.Email + " shared a document with you!",
		Body:    "Click the following link to view and edit the document: " + os.Getenv("FRONT_END_URL") + "/document/" + strconv.Itoa(int(documentID)),
	}
	if err := dc.MS.SendMail(mail); err != nil {
		responses.NewResponse("Failed to send email", 500).Send(ctx)
		return
	}
	responses.NewResponse("Document shared successfully", 200).Send(ctx)
}

func (dc *DocumentController) RemoveSharedUser(ctx *gin.Context) {
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
	var document dbmodels.Document
	if err := dc.DB.Where("id = ? AND user_id = ?", documentID, userID).First(&document).Error; err != nil {
		responses.NewResponse("Document not found or unauthorized", http.StatusNotFound).Send(ctx)
		return
	}
	var documentUser dbmodels.DocumentUser
	if err := dc.DB.Where("document_id = ? AND user_id = ?", documentID, userID).First(&documentUser).Error; err != nil {
		responses.NewResponse("DocumentUser association not found", http.StatusNotFound).Send(ctx)
		return
	}
	if err := dc.DB.Delete(&documentUser).Error; err != nil {
		responses.NewResponse("Failed to delete document user association", http.StatusInternalServerError).Send(ctx)
		return
	}
	responses.NewResponse("User removed successfully", 200).Send(ctx)
}
