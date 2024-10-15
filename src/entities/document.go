package entities

import (
	"time"

	dbmodels "github.com/niteshsiingh/doc-server/src/database/db-models"
	"gorm.io/datatypes"
)

type UpdateDocumentRequest struct {
	Title       string         `json:"title"`
	Content     datatypes.JSON `json:"content"`
	IsPublic    bool           `json:"isPublic"`
	User        dbmodels.User  `json:"user"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	ID          int            `json:"id"`
	Permissions struct {
		ID         int       `json:"id"`
		CreatedAt  time.Time `json:"createdAt"`
		UpdatedAt  time.Time `json:"updatedAt"`
		DocumentID int       `json:"document_id"`
	} `json:"permissions"`
	UserID int         `json:"userId"`
	Users  interface{} `json:"users"`
}

type UpdateDocumentRequest2 struct {
	Title     string         `json:"title"`
	Content   datatypes.JSON `json:"content"`
	IsPublic  bool           `json:"isPublic"`
	User      dbmodels.User  `json:"user"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}
