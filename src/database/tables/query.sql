-- name: CreateUser :exec
INSERT INTO users (email, password, verification_token, user_roles)
VALUES ($1, $2, $3, $4);

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: EditUser :exec
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
WHERE id = $10;

-- name: CreateDocument :one
INSERT INTO documents (title, user_id, is_public, body)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: DeleteDocument :exec
UPDATE documents
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: EditDocument :exec
UPDATE documents
SET title = $1,
    user_id = $2,
    is_public = $3,
    body = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $5;

-- name: CreateDocumentUser :exec
INSERT INTO document_users (user_id, document_id, version, read, write, download, share, admin)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: FindDocumentUser :one
SELECT * FROM document_users
WHERE document_id = $1 AND user_id = $2
LIMIT 1;

-- name: DeleteDocumentUser :exec
DELETE FROM document_users
WHERE id = $1;

-- name: EditDocumentUser :exec
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
WHERE id = $9;

-- name: CreateIdentifier :exec
INSERT INTO identifiers (elem, id, doc_id)
VALUES ($1, $2 , $3);

-- name: DeleteIdentifier :exec
DELETE FROM identifiers
WHERE doc_id = $1 AND id::jsonb @> $2::jsonb;

-- name: EditIdentifier :exec
UPDATE identifiers
SET elem = $1,
    id = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE identifier_id = $1;


-- name: CreateIdentifierId :exec
INSERT INTO identifier_id (_base, _c)
VALUES ($1, $2);

-- name: DeleteIdentifierId :exec
DELETE FROM identifier_id
WHERE id = $1;

-- name: EditIdentifierId :exec
UPDATE identifier_id
SET _base = $1,
    _c = $2
WHERE id = $3;

-- name: CreateRole :exec
INSERT INTO role (role)
VALUES ($1);

-- name: DeleteRole :exec
DELETE FROM role
WHERE id = $1;

-- name: EditRole :exec
UPDATE role
SET role = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: CreateUserRole :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2);

-- name: DeleteUserRole :exec
DELETE FROM user_roles
WHERE id = $1;

-- name: EditUserRole :exec
UPDATE user_roles
SET user_id = $1,
    role_id = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $3;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByVerificationToken :one
SELECT * FROM users
WHERE verification_token = $1 AND email = $2 LIMIT 1;

-- name: GetUserByPasswordResetToken :one
SELECT * FROM users
WHERE password_reset_token = $1 AND email = $2 LIMIT 1;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1 LIMIT 1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, user_id)
VALUES ($1, $2);

-- name: GetAllDocuments :many
SELECT * FROM documents
WHERE user_id = $1 OR is_public = true
ORDER BY updated_at DESC;

-- name: GetDocumentById :one
SELECT * FROM documents
WHERE id = $1 AND (is_public = true OR user_id = $2)
LIMIT 1;

-- name: FindDocument :one
SELECT * FROM documents
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: GetUserRole :one
SELECT * FROM role
WHERE id = $1 LIMIT 1;

-- name: GetAllIdentifiers :many
SELECT * FROM identifiers
WHERE doc_id = $1
ORDER BY created_at ASC;