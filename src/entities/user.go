package entities

import (
	"time"

	dbmodels "github.com/niteshsiingh/doc-server/src/database/db-models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PERM_NONE              = 0
	PERM_DOCUMENT_ADMIN    = 1 << 0
	PERM_DOCUMENT_WRITE    = 1 << 1
	PERM_DOCUMENT_READ     = 1 << 2
	PERM_DOCUMENT_SHARE    = 1 << 3
	PERM_DOCUMENT_DOWNLOAD = 1 << 4
)

type JWTUser struct {
	ID        uint                `json:"_id"`
	EmailID   string              `json:"emailId"`
	Documents map[string]int64    `json:"documents"`
	Roles     []dbmodels.UserRole `json:"roles"`
}

type AuthProvider struct {
	Type     string `bson:"type" json:"type"`
	Phone    string `bson:"phone,omitempty" json:"phone,omitempty"`
	Email    string `bson:"email,omitempty" json:"email,omitempty"`
	Password string `bson:"password,omitempty"`
}

type UserDocument struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID     primitive.ObjectID `bson:"_idUser" json:"_idUser"`
	DocumentID primitive.ObjectID `bson:"_idDocument" json:"_idDocument"`
	Status     int                `bson:"status" json:"status"`

	Permissions DocumentPermission `bson:"permissions" json:"permissions"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type DocumentPermission struct {
	Version    int            `json:"version" bson:"version"`
	Permission AllPermissions `json:"permission" bson:"permission"`
}

type AllPermissions struct {
	Read     bool `bson:"read" json:"read"`
	Write    bool `bson:"write" json:"write"`
	Download bool `bson:"download" json:"download"`
	Share    bool `bson:"share" json:"share"`
	Admin    bool `bson:"admin" json:"admin"`
}

type User struct {
	ID                       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	AuthProviders            []AuthProvider     `bson:"authProviders,omitempty" json:"-"`
	Documents                []UserDocument     `json:"permissions"`
	FirebaseUID              string             `bson:"uid" json:"uid,omitempty"`
	FirstName                string             `bson:"firstName" json:"firstName"`
	LastName                 string             `bson:"lastName" json:"lastName"`
	Email                    string             `bson:"email" json:"email,omitempty"`
	Phone                    string             `bson:"phone," json:"phone,omitempty"`
	CountryCode              string             `bson:"country" json:"country"`
	NotificationSubscription NotificationSubs   `bson:"notificationSubscription,omitempty" json:"notificationSubscription,omitempty"`
	IsDeleted                bool               `bson:"isDeleted" json:"isDeleted"`
	IsDisabled               bool               `bson:"isDisabled" json:"isDisabled"`
	IsActive                 bool               `bson:"isActive,omitempty" json:"isActive"`
	IsSubscriptionEnabled    bool               `bson:"isSubscriptionEnabled,omitempty" json:"isSubscriptionEnabled"`
	CreatedAt                time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt                time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type NotificationSubs struct {
	Day  string                       `json:"day,omitempty" bson:"day,omitempty"`
	Time NotificationSubscriptionTime `json:"time,omitempty" bson:"time,omitempty"`
}

type NotificationSubscriptionTime struct {
	Hour     string `json:"hour,omitempty" bson:"hour,omitempty"`
	Minute   string `json:"minute,omitempty" bson:"minute,omitempty"`
	Maredian string `json:"maredian,omitempty" bson:"maredian,omitempty"`
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
		permissions[document.DocumentID.Hex()] = permission
	}
	return permissions
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password1"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type ConfirmResetPasswordRequest struct {
	Password string `json:"password"`
}
