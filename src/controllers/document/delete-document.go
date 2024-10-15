package document

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"github.com/niteshsiingh/doc-server/src/responses"
)

func (dc *DocumentController) DeleteDocument(ctx *gin.Context) {
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
	_, _, err = middleware.GetAuth().ParseToken(authToken)
	if err != nil {
		responses.NewResponse("Invalid token", 403).Send(ctx)
		return
	}
	// userID := parsedToken.UserID
	// var document dbmodels.Document
	// err = dc.DB.Where("id = ? and user_id = ?", documentID, userID).First(&document).Error
	// if err != nil {
	// 	if err == gorm.ErrRecordNotFound {
	// 		responses.NewResponse("Document not found", 404).Send(ctx)
	// 		return
	// 	}
	// 	responses.NewResponse("Failed to retrieve document", 500).Send(ctx)
	// 	return
	// }
	err = dc.DB.DeleteDocument(cxt, int32(documentID))
	if err != nil {
		responses.NewResponse("Failed to delete document", 500).Send(ctx)
		return
	}
	responses.NewResponse("Document deleted successfully", 200).Send(ctx)
}
