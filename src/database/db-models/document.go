package dbmodels

import (
	"time"
)

type Document struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	DeletedAt   *time.Time         `gorm:"index" json:"deletedAt,omitempty"`
	Title       string             `gorm:"type:string;not null" json:"title"`
	BodyID      uint               `gorm:"not null" json:"bodyId"`
	Body        Identifier         `gorm:"foreignKey:BodyID" json:"body"`
	UserID      uint               `gorm:"not null" json:"userId"`
	User        User               `gorm:"foreignKey:UserID" json:"user"`
	Users       []DocumentUser     `gorm:"constraint:OnDelete:CASCADE;" json:"users"`
	IsPublic    bool               `gorm:"default:false" json:"isPublic"`
	Permissions DocumentPermission `gorm:"foreignKey:DocumentID" json:"permissions"`
}

type DocumentPermission struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  *time.Time     `gorm:"index" json:"deletedAt,omitempty"`
	DocumentID uint           `gorm:"not null" json:"document_id"`
	Version    int            `gorm:"column:version" json:"version"`
	Permission AllPermissions `gorm:"embedded" json:"permission"`
}

type AllPermissions struct {
	Read     bool `gorm:"column:read" json:"read"`
	Write    bool `gorm:"column:write" json:"write"`
	Download bool `gorm:"column:download" json:"download"`
	Share    bool `gorm:"column:share" json:"share"`
	Admin    bool `gorm:"column:admin" json:"admin"`
}
