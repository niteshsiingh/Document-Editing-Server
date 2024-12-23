// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package databases

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createDocument = `-- name: CreateDocument :one
INSERT INTO documents (title, user_id, is_public, body)
VALUES ($1, $2, $3, $4)
RETURNING id
`

type CreateDocumentParams struct {
	Title    pgtype.Text
	UserID   pgtype.Int4
	IsPublic pgtype.Bool
	Body     []byte
}

func (q *Queries) CreateDocument(ctx context.Context, arg CreateDocumentParams) (int32, error) {
	row := q.db.QueryRow(ctx, createDocument,
		arg.Title,
		arg.UserID,
		arg.IsPublic,
		arg.Body,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const createDocumentUser = `-- name: CreateDocumentUser :exec
INSERT INTO document_users (user_id, document_id, version, read, write, download, share, admin)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

type CreateDocumentUserParams struct {
	UserID     pgtype.Int4
	DocumentID pgtype.Int4
	Version    pgtype.Int4
	Read       pgtype.Bool
	Write      pgtype.Bool
	Download   pgtype.Bool
	Share      pgtype.Bool
	Admin      pgtype.Bool
}

func (q *Queries) CreateDocumentUser(ctx context.Context, arg CreateDocumentUserParams) error {
	_, err := q.db.Exec(ctx, createDocumentUser,
		arg.UserID,
		arg.DocumentID,
		arg.Version,
		arg.Read,
		arg.Write,
		arg.Download,
		arg.Share,
		arg.Admin,
	)
	return err
}

const createIdentifier = `-- name: CreateIdentifier :exec
INSERT INTO identifiers (elem, id, doc_id)
VALUES ($1, $2 , $3)
`

type CreateIdentifierParams struct {
	Elem  string
	ID    []byte
	DocID pgtype.Int4
}

func (q *Queries) CreateIdentifier(ctx context.Context, arg CreateIdentifierParams) error {
	_, err := q.db.Exec(ctx, createIdentifier, arg.Elem, arg.ID, arg.DocID)
	return err
}

const createIdentifierId = `-- name: CreateIdentifierId :exec
INSERT INTO identifier_id (_base, _c)
VALUES ($1, $2)
`

type CreateIdentifierIdParams struct {
	Base []byte
	C    []byte
}

func (q *Queries) CreateIdentifierId(ctx context.Context, arg CreateIdentifierIdParams) error {
	_, err := q.db.Exec(ctx, createIdentifierId, arg.Base, arg.C)
	return err
}

const createRefreshToken = `-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, user_id)
VALUES ($1, $2)
`

type CreateRefreshTokenParams struct {
	Token  pgtype.Text
	UserID pgtype.Int4
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) error {
	_, err := q.db.Exec(ctx, createRefreshToken, arg.Token, arg.UserID)
	return err
}

const createRole = `-- name: CreateRole :exec
INSERT INTO role (role)
VALUES ($1)
`

func (q *Queries) CreateRole(ctx context.Context, role pgtype.Text) error {
	_, err := q.db.Exec(ctx, createRole, role)
	return err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (email, password, verification_token, user_roles)
VALUES ($1, $2, $3, $4)
`

type CreateUserParams struct {
	Email             pgtype.Text
	Password          pgtype.Text
	VerificationToken pgtype.Text
	UserRoles         []byte
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser,
		arg.Email,
		arg.Password,
		arg.VerificationToken,
		arg.UserRoles,
	)
	return err
}

const createUserRole = `-- name: CreateUserRole :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
`

type CreateUserRoleParams struct {
	UserID pgtype.Int4
	RoleID pgtype.Int4
}

func (q *Queries) CreateUserRole(ctx context.Context, arg CreateUserRoleParams) error {
	_, err := q.db.Exec(ctx, createUserRole, arg.UserID, arg.RoleID)
	return err
}

const deleteDocument = `-- name: DeleteDocument :exec
UPDATE documents
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
`

func (q *Queries) DeleteDocument(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteDocument, id)
	return err
}

const deleteDocumentUser = `-- name: DeleteDocumentUser :exec
DELETE FROM document_users
WHERE id = $1
`

func (q *Queries) DeleteDocumentUser(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteDocumentUser, id)
	return err
}

const deleteIdentifier = `-- name: DeleteIdentifier :exec
DELETE FROM identifiers
WHERE doc_id = $1 AND id::jsonb @> $2::jsonb
`

type DeleteIdentifierParams struct {
	DocID   pgtype.Int4
	Column2 []byte
}

func (q *Queries) DeleteIdentifier(ctx context.Context, arg DeleteIdentifierParams) error {
	_, err := q.db.Exec(ctx, deleteIdentifier, arg.DocID, arg.Column2)
	return err
}

const deleteIdentifierId = `-- name: DeleteIdentifierId :exec
DELETE FROM identifier_id
WHERE id = $1
`

func (q *Queries) DeleteIdentifierId(ctx context.Context, id pgtype.Int4) error {
	_, err := q.db.Exec(ctx, deleteIdentifierId, id)
	return err
}

const deleteRefreshToken = `-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE user_id = $1
`

func (q *Queries) DeleteRefreshToken(ctx context.Context, userID pgtype.Int4) error {
	_, err := q.db.Exec(ctx, deleteRefreshToken, userID)
	return err
}

const deleteRole = `-- name: DeleteRole :exec
DELETE FROM role
WHERE id = $1
`

func (q *Queries) DeleteRole(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteRole, id)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const deleteUserRole = `-- name: DeleteUserRole :exec
DELETE FROM user_roles
WHERE id = $1
`

func (q *Queries) DeleteUserRole(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUserRole, id)
	return err
}

const editDocument = `-- name: EditDocument :exec
UPDATE documents
SET title = $1,
    user_id = $2,
    is_public = $3,
    body = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $5
`

type EditDocumentParams struct {
	Title    pgtype.Text
	UserID   pgtype.Int4
	IsPublic pgtype.Bool
	Body     []byte
	ID       int32
}

func (q *Queries) EditDocument(ctx context.Context, arg EditDocumentParams) error {
	_, err := q.db.Exec(ctx, editDocument,
		arg.Title,
		arg.UserID,
		arg.IsPublic,
		arg.Body,
		arg.ID,
	)
	return err
}

const editDocumentUser = `-- name: EditDocumentUser :exec
UPDATE document_users
SET user_id = $1,
    document_id = $2,
    version = $3,
    read = $4,
    write = $5,
    download = $6,
    share = $7,
    admin = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $9
`

type EditDocumentUserParams struct {
	UserID     pgtype.Int4
	DocumentID pgtype.Int4
	Version    pgtype.Int4
	Read       pgtype.Bool
	Write      pgtype.Bool
	Download   pgtype.Bool
	Share      pgtype.Bool
	Admin      pgtype.Bool
	ID         int32
}

func (q *Queries) EditDocumentUser(ctx context.Context, arg EditDocumentUserParams) error {
	_, err := q.db.Exec(ctx, editDocumentUser,
		arg.UserID,
		arg.DocumentID,
		arg.Version,
		arg.Read,
		arg.Write,
		arg.Download,
		arg.Share,
		arg.Admin,
		arg.ID,
	)
	return err
}

const editIdentifier = `-- name: EditIdentifier :exec
UPDATE identifiers
SET elem = $1,
    id = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE identifier_id = $1
`

type EditIdentifierParams struct {
	Elem string
	ID   []byte
}

func (q *Queries) EditIdentifier(ctx context.Context, arg EditIdentifierParams) error {
	_, err := q.db.Exec(ctx, editIdentifier, arg.Elem, arg.ID)
	return err
}

const editIdentifierId = `-- name: EditIdentifierId :exec
UPDATE identifier_id
SET _base = $1,
    _c = $2
WHERE id = $3
`

type EditIdentifierIdParams struct {
	Base []byte
	C    []byte
	ID   pgtype.Int4
}

func (q *Queries) EditIdentifierId(ctx context.Context, arg EditIdentifierIdParams) error {
	_, err := q.db.Exec(ctx, editIdentifierId, arg.Base, arg.C, arg.ID)
	return err
}

const editRole = `-- name: EditRole :exec
UPDATE role
SET role = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2
`

type EditRoleParams struct {
	Role pgtype.Text
	ID   int32
}

func (q *Queries) EditRole(ctx context.Context, arg EditRoleParams) error {
	_, err := q.db.Exec(ctx, editRole, arg.Role, arg.ID)
	return err
}

const editUser = `-- name: EditUser :exec
UPDATE users
SET email = $1,
    password = $2,
    is_verified = $3,
    verification_token = $4,
    password_reset_token = $5,
    user_roles = $6,
    shared_documents = $7,
    documents = $8,
    refresh_tokens = $9,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $10
`

type EditUserParams struct {
	Email              pgtype.Text
	Password           pgtype.Text
	IsVerified         pgtype.Bool
	VerificationToken  pgtype.Text
	PasswordResetToken pgtype.Text
	UserRoles          []byte
	SharedDocuments    []byte
	Documents          []byte
	RefreshTokens      []byte
	ID                 int32
}

func (q *Queries) EditUser(ctx context.Context, arg EditUserParams) error {
	_, err := q.db.Exec(ctx, editUser,
		arg.Email,
		arg.Password,
		arg.IsVerified,
		arg.VerificationToken,
		arg.PasswordResetToken,
		arg.UserRoles,
		arg.SharedDocuments,
		arg.Documents,
		arg.RefreshTokens,
		arg.ID,
	)
	return err
}

const editUserRole = `-- name: EditUserRole :exec
UPDATE user_roles
SET user_id = $1,
    role_id = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $3
`

type EditUserRoleParams struct {
	UserID pgtype.Int4
	RoleID pgtype.Int4
	ID     int32
}

func (q *Queries) EditUserRole(ctx context.Context, arg EditUserRoleParams) error {
	_, err := q.db.Exec(ctx, editUserRole, arg.UserID, arg.RoleID, arg.ID)
	return err
}

const findDocument = `-- name: FindDocument :one
SELECT id, created_at, updated_at, deleted_at, title, user_id, is_public, body FROM documents
WHERE id = $1 AND user_id = $2
LIMIT 1
`

type FindDocumentParams struct {
	ID     int32
	UserID pgtype.Int4
}

func (q *Queries) FindDocument(ctx context.Context, arg FindDocumentParams) (Document, error) {
	row := q.db.QueryRow(ctx, findDocument, arg.ID, arg.UserID)
	var i Document
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Title,
		&i.UserID,
		&i.IsPublic,
		&i.Body,
	)
	return i, err
}

const findDocumentUser = `-- name: FindDocumentUser :one
SELECT id, created_at, updated_at, user_id, document_id, version, read, write, download, share, admin FROM document_users
WHERE document_id = $1 AND user_id = $2
LIMIT 1
`

type FindDocumentUserParams struct {
	DocumentID pgtype.Int4
	UserID     pgtype.Int4
}

func (q *Queries) FindDocumentUser(ctx context.Context, arg FindDocumentUserParams) (DocumentUser, error) {
	row := q.db.QueryRow(ctx, findDocumentUser, arg.DocumentID, arg.UserID)
	var i DocumentUser
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.DocumentID,
		&i.Version,
		&i.Read,
		&i.Write,
		&i.Download,
		&i.Share,
		&i.Admin,
	)
	return i, err
}

const getAllDocuments = `-- name: GetAllDocuments :many
SELECT id, created_at, updated_at, deleted_at, title, user_id, is_public, body FROM documents
WHERE user_id = $1 OR is_public = true
ORDER BY updated_at DESC
`

func (q *Queries) GetAllDocuments(ctx context.Context, userID pgtype.Int4) ([]Document, error) {
	rows, err := q.db.Query(ctx, getAllDocuments, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Document
	for rows.Next() {
		var i Document
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Title,
			&i.UserID,
			&i.IsPublic,
			&i.Body,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllIdentifiers = `-- name: GetAllIdentifiers :many
SELECT identifier_id, created_at, doc_id, elem, id FROM identifiers
WHERE doc_id = $1
ORDER BY created_at ASC
`

func (q *Queries) GetAllIdentifiers(ctx context.Context, docID pgtype.Int4) ([]Identifier, error) {
	rows, err := q.db.Query(ctx, getAllIdentifiers, docID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Identifier
	for rows.Next() {
		var i Identifier
		if err := rows.Scan(
			&i.IdentifierID,
			&i.CreatedAt,
			&i.DocID,
			&i.Elem,
			&i.ID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDocumentById = `-- name: GetDocumentById :one
SELECT id, created_at, updated_at, deleted_at, title, user_id, is_public, body FROM documents
WHERE id = $1 AND (is_public = true OR user_id = $2)
LIMIT 1
`

type GetDocumentByIdParams struct {
	ID     int32
	UserID pgtype.Int4
}

func (q *Queries) GetDocumentById(ctx context.Context, arg GetDocumentByIdParams) (Document, error) {
	row := q.db.QueryRow(ctx, getDocumentById, arg.ID, arg.UserID)
	var i Document
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Title,
		&i.UserID,
		&i.IsPublic,
		&i.Body,
	)
	return i, err
}

const getRefreshToken = `-- name: GetRefreshToken :one
SELECT id, created_at, updated_at, token, user_id FROM refresh_tokens
WHERE token = $1 LIMIT 1
`

func (q *Queries) GetRefreshToken(ctx context.Context, token pgtype.Text) (RefreshToken, error) {
	row := q.db.QueryRow(ctx, getRefreshToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Token,
		&i.UserID,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, deleted_at, email, password, is_verified, verification_token, password_reset_token, user_roles, shared_documents, documents, refresh_tokens FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email pgtype.Text) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Email,
		&i.Password,
		&i.IsVerified,
		&i.VerificationToken,
		&i.PasswordResetToken,
		&i.UserRoles,
		&i.SharedDocuments,
		&i.Documents,
		&i.RefreshTokens,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, created_at, updated_at, deleted_at, email, password, is_verified, verification_token, password_reset_token, user_roles, shared_documents, documents, refresh_tokens FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Email,
		&i.Password,
		&i.IsVerified,
		&i.VerificationToken,
		&i.PasswordResetToken,
		&i.UserRoles,
		&i.SharedDocuments,
		&i.Documents,
		&i.RefreshTokens,
	)
	return i, err
}

const getUserByPasswordResetToken = `-- name: GetUserByPasswordResetToken :one
SELECT id, created_at, updated_at, deleted_at, email, password, is_verified, verification_token, password_reset_token, user_roles, shared_documents, documents, refresh_tokens FROM users
WHERE password_reset_token = $1 AND email = $2 LIMIT 1
`

type GetUserByPasswordResetTokenParams struct {
	PasswordResetToken pgtype.Text
	Email              pgtype.Text
}

func (q *Queries) GetUserByPasswordResetToken(ctx context.Context, arg GetUserByPasswordResetTokenParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserByPasswordResetToken, arg.PasswordResetToken, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Email,
		&i.Password,
		&i.IsVerified,
		&i.VerificationToken,
		&i.PasswordResetToken,
		&i.UserRoles,
		&i.SharedDocuments,
		&i.Documents,
		&i.RefreshTokens,
	)
	return i, err
}

const getUserByVerificationToken = `-- name: GetUserByVerificationToken :one
SELECT id, created_at, updated_at, deleted_at, email, password, is_verified, verification_token, password_reset_token, user_roles, shared_documents, documents, refresh_tokens FROM users
WHERE verification_token = $1 AND email = $2 LIMIT 1
`

type GetUserByVerificationTokenParams struct {
	VerificationToken pgtype.Text
	Email             pgtype.Text
}

func (q *Queries) GetUserByVerificationToken(ctx context.Context, arg GetUserByVerificationTokenParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserByVerificationToken, arg.VerificationToken, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Email,
		&i.Password,
		&i.IsVerified,
		&i.VerificationToken,
		&i.PasswordResetToken,
		&i.UserRoles,
		&i.SharedDocuments,
		&i.Documents,
		&i.RefreshTokens,
	)
	return i, err
}

const getUserRole = `-- name: GetUserRole :one
SELECT id, created_at, updated_at, role FROM role
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserRole(ctx context.Context, id int32) (Role, error) {
	row := q.db.QueryRow(ctx, getUserRole, id)
	var i Role
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Role,
	)
	return i, err
}
