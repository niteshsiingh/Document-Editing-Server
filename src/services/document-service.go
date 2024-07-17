package services

import (
	dbmodels "github.com/niteshsiingh/doc-server/src/database/db-models"
	"gorm.io/gorm"
)

func FindDocumentByID(id uint, userID uint, db *gorm.DB) (*dbmodels.Document, error) {
	var document dbmodels.Document

	err := db.Where("id = ? AND (user_id = ? OR is_public = ?)", id, userID, true).First(&document).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			var sharedDocument dbmodels.DocumentUser
			err = db.Where("user_id = ? AND id = ?", userID, id).First(&sharedDocument).Error
			if err != nil {
				if gorm.ErrRecordNotFound == err {
					return nil, nil
				}
				return nil, err
			}
			document = sharedDocument.Document
		} else {
			return nil, err
		}
	}

	return &document, nil
}
