package dbmodels

import "gorm.io/gorm"

type PermissionEnum string

const (
	VIEW PermissionEnum = "VIEW"
	EDIT PermissionEnum = "EDIT"
)

type DocumentUser struct {
	gorm.Model
	UserID     uint           `gorm:"primaryKey"`
	DocumentID uint           `gorm:"primaryKey"`
	Permission PermissionEnum `gorm:"type:varchar(255)"`
	User       User           `gorm:"foreignKey:UserID"`
	Document   Document       `gorm:"foreignKey:DocumentID;constraint:OnDelete:CASCADE;"`
}
