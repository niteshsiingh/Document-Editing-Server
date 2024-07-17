package document

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/niteshsiingh/doc-server/src/entities"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
	"gorm.io/datatypes"
)

func (dc *DocumentController) UpdateDocument(ctx *gin.Context) {
	documentIDstr := ctx.Param("document_id")
	documentID, err := strconv.ParseUint(documentIDstr, 10, 64)
	if err != nil {
		responses.NewResponse("Invalid json request", 400).Send(ctx)
		return
	}
	var updateDocumentRequest entities.UpdateDocumentRequest
	err = ctx.ShouldBindJSON(&updateDocumentRequest)
	if err != nil {
		fmt.Println(err)
		responses.NewResponse("Invalid json request", 400).Send(ctx)
		return
	}
	document, err := services.FindDocumentByID(uint(documentID), uint(updateDocumentRequest.UserID), dc.DB)
	if err != nil {
		responses.NewResponse("Document not found", 404).Send(ctx)
		return
	}
	document.Title = updateDocumentRequest.Title

	document.Content = datatypes.JSON(updateDocumentRequest.Content)
	document.IsPublic = updateDocumentRequest.IsPublic

	err = dc.DB.Save(&document).Error
	if err != nil {
		responses.NewResponse("Failed to update document", 500).Send(ctx)
		return
	}
	responses.NewResponse("Document updated successfully", 200).Send(ctx)
}
