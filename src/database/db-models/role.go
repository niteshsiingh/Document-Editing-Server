package dbmodels

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name      string     `gorm:"type:varchar(20);not null;check:name IN ('ADMIN', 'SUPERADMIN')"`
	Users     []*User    `gorm:"many2many:user_roles;"`
	RoleUsers []UserRole `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE;"`
}
