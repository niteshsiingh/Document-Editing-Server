package dbmodels

import "gorm.io/gorm"

type UserRole struct {
	gorm.Model
	UserID uint   `gorm:"primaryKey;autoIncrement:false"`
	RoleID uint   `gorm:"primaryKey;autoIncrement:false"`
	Role   string `gorm:"foreignKey:RoleID"`
}

func (UserRole) TableName() string {
	return "user_role"
}
