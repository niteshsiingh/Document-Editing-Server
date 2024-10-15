package document

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/niteshsiingh/doc-server/src/database/tables/databases"
	"github.com/niteshsiingh/doc-server/src/entities"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
)

func (dc *DocumentController) UpdateDocument(ctx *gin.Context) {
	cxt := context.Background()
	documentIDstr := ctx.Param("document_id")
	documentID, err := strconv.ParseUint(documentIDstr, 10, 64)
	if err != nil {
		responses.NewResponse("Invalid json request", 400).Send(ctx)
		return
	}
	var updateDocumentRequest entities.UpdateDocumentRequest
	err = ctx.ShouldBindJSON(&updateDocumentRequest)
	if err != nil {
		responses.NewResponse("Invalid json request", 400).Send(ctx)
		return
	}
	document, err := services.FindDocumentByID(cxt, uint(documentID), uint(updateDocumentRequest.UserID), dc.DB)
	if err != nil {
		responses.NewResponse("Document not found", 404).Send(ctx)
		return
	}
	document.Title = pgtype.Text{String: updateDocumentRequest.Title, Valid: true}
	document.IsPublic = pgtype.Bool{Bool: updateDocumentRequest.IsPublic, Valid: true}
	err = dc.DB.EditDocument(cxt, databases.EditDocumentParams{
		Title:    document.Title,
		UserID:   document.UserID,
		IsPublic: document.IsPublic,
		Body:     document.Body,
		ID:       document.ID,
	})
	if err != nil {
		responses.NewResponse("Failed to update document", 500).Send(ctx)
		return
	}
	responses.NewResponse("Document updated successfully", 200).Send(ctx)
}
