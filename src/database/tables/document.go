package tables

import "gorm.io/gorm"

func CreateDocumentTable(db *gorm.DB) {
	queryString := `
	CREATE TABLE IF NOT EXISTS documents (
		id INTEGER NOT NULL,
		title VARCHAR(255),
		content TEXT,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deletedAt TIMESTAMP DEFAULT NULL,
		isPublic BOOLEAN DEFAULT FALSE,
		userId INTEGER NOT NULL,
		users JSONB DEFAULT '[]'
	) PARTITION BY RANGE (document_id);
	`
	txDB := db.Exec(queryString)
	txDB.Commit()
}

func CreateDocumentPermissionTable(db *gorm.DB) {
	queryString := `
	CREATE TABLE IF NOT EXISTS document_users (
		id INTEGER NOT NULL,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		user_id INTEGER NOT NULL,
		document_id INTEGER NOT NULL,
		version INTEGER NOT NULL,
		read BOOLEAN DEFAULT FALSE,
		write BOOLEAN DEFAULT FALSE,
		download BOOLEAN DEFAULT FALSE,
		share BOOLEAN DEFAULT FALSE,
		admin BOOLEAN DEFAULT FALSE
	) PARTITION BY RANGE (id);
	`
	txDB := db.Exec(queryString)
	txDB.Commit()
}