// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package databases

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Base struct {
	B pgtype.Int4
}

type Document struct {
	ID        int32
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	DeletedAt pgtype.Timestamp
	Title     pgtype.Text
	UserID    pgtype.Int4
	IsPublic  pgtype.Bool
	Body      []byte
}

type DocumentUser struct {
	ID         int32
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
	UserID     pgtype.Int4
	DocumentID pgtype.Int4
	Version    pgtype.Int4
	Read       pgtype.Bool
	Write      pgtype.Bool
	Download   pgtype.Bool
	Share      pgtype.Bool
	Admin      pgtype.Bool
}

type Identifier struct {
	IdentifierID int32
	CreatedAt    pgtype.Timestamp
	DocID        pgtype.Int4
	Elem         string
	ID           []byte
}

type IdentifierID struct {
	ID   pgtype.Int4
	Base []byte
	C    []byte
	D    []byte
	S    []byte
}

type RefreshToken struct {
	ID        int32
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	Token     pgtype.Text
	UserID    pgtype.Int4
}

type Role struct {
	ID        int32
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	Role      pgtype.Text
}

type User struct {
	ID                 int32
	CreatedAt          pgtype.Timestamp
	UpdatedAt          pgtype.Timestamp
	DeletedAt          pgtype.Timestamp
	Email              pgtype.Text
	Password           pgtype.Text
	IsVerified         pgtype.Bool
	VerificationToken  pgtype.Text
	PasswordResetToken pgtype.Text
	UserRoles          []byte
	SharedDocuments    []byte
	Documents          []byte
	RefreshTokens      []byte
}

type UserRole struct {
	ID        int32
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	UserID    pgtype.Int4
	RoleID    pgtype.Int4
}
