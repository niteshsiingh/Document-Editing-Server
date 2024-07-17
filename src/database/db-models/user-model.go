package dbmodels

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

const (
	PERM_NONE              = 0
	PERM_DOCUMENT_ADMIN    = 1 << 0
	PERM_DOCUMENT_WRITE    = 1 << 1
	PERM_DOCUMENT_READ     = 1 << 2
	PERM_DOCUMENT_SHARE    = 1 << 3
	PERM_DOCUMENT_DOWNLOAD = 1 << 4
)

type User struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          *time.Time     `gorm:"index" json:"deletedAt,omitempty"`
	Email              string         `gorm:"type:varchar(255);not null" json:"email"`
	Password           string         `gorm:"type:varchar(255);not null" json:"password"`
	IsVerified         bool           `gorm:"not null" json:"is_verified"`
	VerificationToken  string         `gorm:"type:varchar(255);not null" json:"verification_token"`
	PasswordResetToken string         `gorm:"type:varchar(255);not null" json:"password_reset_token"`
	RefreshTokens      []RefreshToken `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"refresh_tokens"`
	Roles              []Role         `gorm:"many2many:user_roles;" json:"roles"`
	UserRoles          []UserRole     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user_roles"`
	SharedDocuments    []DocumentUser `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"shared_documents"`
	Documents          []Document     `gorm:"many2many:user_documents;" json:"documents"`
}

func GetUserWithRoles(db *gorm.DB) *gorm.DB {
	return db.Preload("UserRoles").Preload("UserRoles.Role")
}

func (u *User) GetDocumentPermissions() map[string]int64 {
	permissions := map[string]int64{}
	for _, document := range u.Documents {
		permission := int64(0)
		if document.Permissions.Permission.Admin {
			permission = permission | PERM_DOCUMENT_ADMIN |
				PERM_DOCUMENT_READ | PERM_DOCUMENT_WRITE | PERM_DOCUMENT_DOWNLOAD |
				PERM_DOCUMENT_SHARE
		}
		if document.Permissions.Permission.Read {
			permission = permission | PERM_DOCUMENT_READ
		}
		if document.Permissions.Permission.Write {
			permission = permission | PERM_DOCUMENT_READ |
				PERM_DOCUMENT_WRITE
		}
		if document.Permissions.Permission.Download {
			permission = permission | PERM_DOCUMENT_READ |
				PERM_DOCUMENT_DOWNLOAD
		}
		if document.Permissions.Permission.Share {
			permission = permission | PERM_DOCUMENT_READ |
				PERM_DOCUMENT_WRITE | PERM_DOCUMENT_DOWNLOAD
		}
		docHexID := strconv.Itoa(int(document.ID))
		permissions[docHexID] = permission
	}
	return permissions
}
