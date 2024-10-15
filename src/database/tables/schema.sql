CREATE TABLE IF NOT EXISTS document_users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER,
    document_id INTEGER,
    version INTEGER,
    read BOOLEAN DEFAULT FALSE,
    write BOOLEAN DEFAULT FALSE,
    download BOOLEAN DEFAULT FALSE,
    share BOOLEAN DEFAULT FALSE,
    admin BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    title TEXT ,
    user_id INTEGER ,
    is_public BOOLEAN DEFAULT FALSE,
    body  JSON 
);

CREATE TABLE users (
    id SERIAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    email TEXT ,
    password TEXT ,
    is_verified BOOLEAN DEFAULT FALSE,
    verification_token TEXT ,
    password_reset_token TEXT ,
    user_roles JSONB ,
    shared_documents JSON ,
    documents JSON ,
    refresh_tokens JSON ,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS identifiers (
    identifier_id SERIAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    doc_id INTEGER,
    elem TEXT NOT NULL,
    id JSON ,
    PRIMARY KEY (identifier_id)
);

CREATE TABLE IF NOT EXISTS base (
    _b INTEGER
);

CREATE TABLE IF NOT EXISTS identifier_id (
    id SERIAL,
    _base JSON,
    _c JSON,
    _d JSON,
    _s JSON
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    token TEXT ,
    user_id INTEGER ,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS role (
    id SERIAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role TEXT ,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS user_roles (
    id SERIAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER ,
    role_id INTEGER ,
    PRIMARY KEY (id)
);
