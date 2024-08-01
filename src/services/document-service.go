package services

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/niteshsiingh/doc-server/src/database/tables/databases"
	"gorm.io/gorm"
)

func FindDocumentByID(ctx context.Context, id uint, userID uint, db *databases.Queries) (*databases.Document, error) {
	var document databases.Document
	document, err := db.GetDocumentById(ctx, databases.GetDocumentByIdParams{
		ID:     int32(id),
		UserID: pgtype.Int4{Int32: int32(userID), Valid: true},
	})
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			// var sharedDocument dbmodels.DocumentUser
			// err = db.Where("user_id = ? AND id = ?", userID, id).First(&sharedDocument).Error
			// if err != nil {
			// 	if gorm.ErrRecordNotFound == err {
			// 		return nil, nil
			// 	}
			// 	return nil, err
			// }
			// document = sharedDocument.Document
			return nil, err
		} else {
			return nil, err
		}
	}

	return &document, nil
}

func FindIdentifiersWithDocID(ctx context.Context, id uint, db *databases.Queries) (*[]map[string]interface{}, error) {
	identifiers, err := db.GetAllIdentifiers(ctx, pgtype.Int4{Int32: int32(id), Valid: true})
	if err != nil {
		return nil, err
	}
	var newIdentifiers []map[string]interface{}
	for _, identifier := range identifiers {
		idf := make(map[string]interface{})
		var innerIdf databases.IdentifierID
		err := json.Unmarshal(identifier.ID, &innerIdf)
		if err != nil {
			return nil, err
		}
		iidf := make(map[string]interface{})
		var c interface{}
		err = json.Unmarshal(innerIdf.C, &c)
		if err != nil {
			return nil, err
		}
		var d interface{}
		err = json.Unmarshal(innerIdf.D, &d)
		if err != nil {
			return nil, err
		}
		var s interface{}
		err = json.Unmarshal(innerIdf.S, &s)
		if err != nil {
			return nil, err
		}
		// var b float64
		// err = json.Unmarshal(innerIdf.Base, &b)
		// if err != nil {
		// 	return nil, err
		// }
		var b float64
		err = json.Unmarshal(innerIdf.Base, &b)
		if err != nil {
			return nil, err
		}
		iidf["_base"] = b
		iidf["_base"] = map[string]interface{}{
			"_b": b,
		}
		iidf["_c"] = c
		iidf["_d"] = d
		iidf["_s"] = s
		idf["id"] = iidf
		idf["elem"] = identifier.Elem
		newIdentifiers = append(newIdentifiers, idf)
	}

	return &newIdentifiers, nil
}
