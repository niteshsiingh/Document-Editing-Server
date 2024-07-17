package dbmodels

import "gorm.io/gorm"

type RefreshToken struct {
	gorm.Model
	Token  string `gorm:"type:text;not null"`
	UserID uint   `gorm:"not null"`
	User   User   `gorm:"foreignKey:UserID"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
