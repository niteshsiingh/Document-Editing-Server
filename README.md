# doc-server

## Overview

`doc-server` is a web-based document management server that allows users to upload, manage, and share documents securely. It provides a RESTful API for interacting with documents and supports real-time collaboration through WebSocket connections.

## Features

- **Document Management**: CRUD operations for documents.
- **User Authentication**: Secure user authentication and authorization.
- **Real-time Collaboration**: Real-time document editing and updates using WebSockets.
- **Search Functionality**: Full-text search for documents.
- **Version Control**: Track changes and maintain versions of documents.
- **CRDT-based Editing**: Ensures consistent and concurrent document editing.

## Installation

### Prerequisites

- Go (version 1.16 or higher)
- Postgres
- Redis (for session management)

### Steps

1. **Clone the repository**:
    ```sh
    git clone https://github.com/niteshsiingh/doc-server.git
    cd doc-server
    ```

2. **Install dependencies**:
    ```sh
    go mod tidy
    ```

3. **Set up environment variables**:
    Create a `.env` file in the root directory and add the following variables:
    ```env
    NODE_ENV="dev"
    HOST="localhost"
    PORT="3000"
    DATABASE_URL="postgres://localhost"
    USER="username_postgres"
    PASSWORD="postgres_pass"
    DB_HOST="localhost"
    DB_PORT="5432"
    DATABASE="name"
    SMTP_HOST="smtp.google.com"
    SMTP_USER="example@gmail.com"
    SMTP_PASSWORD="password"
    SMTP_SECURE="false"
    SMTP_PORT="587"
    ACCESS_TOKEN_SECRET="secret"
    ACCESS_TOKEN_EXPIRATION="3600000"
    REFRESH_TOKEN_SECRET="secret"
    REFRESH_TOKEN_EXPIRATION="8640000"
    VERIFY_EMAIL_SECRET="secret"
    PASSWORD_RESET_SECRET="secret"
    PASSWORD_RESET_EXPIRATION="5"
    FRONT_END_URL="http://localhost:3001"
    VALIDITY="3600000"
    REFRESH_VALIDITY="3000000"
    JWT_APP_KEY="secret"
    ```

4. **Run the server**:
    ```sh
    go build -o app ./src/*.go
    ./app
    ```

## Folder Structure


### `/controllers`
Contains the controller logic for handling HTTP requests. Each controller corresponds to a specific resource (e.g., documents, users, auth).

### `/database`
Defines the data models and schemas used in the application. This includes MongoDB models and any other data structures.

### `/services`
Contains the business logic of the application. Services interact with the models and perform operations such as data validation, processing, and interaction with external APIs.


### `/middlewares`
Contains middleware functions that are used to process requests before they reach the controllers. Examples include authentication, logging, and error handling.

### `/utils`
Utility functions and helpers that are used throughout the application. This can include functions for formatting, parsing, and other common tasks.

### `/config`
Configuration files and logic for setting up the application. This includes loading environment variables and setting up database connections.

## Middleware
- **gin.Logger()**: Logs HTTP requests.
- **middleware.CORSMiddleware()**: Handles Cross-Origin Resource Sharing (CORS).
- **gin.Recovery()**: Recovers from any panics and writes a 500 if there was one.

## API Endpoints

### Authentication

- **POST /v1/auth/login**: Login a user.
  - Handler: `ac.Login`
- **POST /v1/auth/refresh-token**: Refresh authentication token.
  - Handler: `ac.RefreshToken`
- **DELETE /v1/auth/logout**: Logout a user.
  - Handler: `ac.Logout`

### User Management

- **PUT /v1/user/verify-email/:token**: Verify user email.
  - Handler: `uc.VerifyEmail`
- **POST /v1/user**: Create a new user.
  - Handler: `uc.CreateUser`
- **GET /v1/user/:id**: Get user by ID.
  - Handler: `uc.GetUserByID`
- **PUT /v1/user/password/:token**: Confirm password reset.
  - Handler: `uc.ConfirmResetPassword`
- **POST /v1/user/reset-password**: Reset user password.
  - Handler: `uc.ResetPassword`

### Documents

- **GET /v1/document**: Get all documents.
  - Handler: `dc.GetAllDocuments`
- **GET /v1/document/:document_id**: Get a document by ID.
  - Handler: `dc.GetOneDocument`
- **GET /v1/document/:document_id/identifiers**: Get document identifiers.
  - Handler: `dc.GetDocumentIdentifiers`
- **PUT /v1/document/:document_id**: Update a document by ID.
  - Handler: `dc.UpdateDocument`
- **POST /v1/document**: Create a new document.
  - Handler: `dc.CreateDocument`
- **POST /v1/document/:document_id/share**: Share a document.
  - Handler: `dc.ShareDocument`
- **DELETE /v1/document/:document_id**: Delete a document by ID.
  - Handler: `dc.DeleteDocument`
- **DELETE /v1/document/:document_id/share**: Remove a shared user from a document.
  - Handler: `dc.RemoveSharedUser`

### WebSocket

- **/ws**: WebSocket endpoint for real-time collaboration.

## Usage

### Real-time Collaboration

Connect to the WebSocket endpoint `/ws` to receive real-time updates and collaborate on documents.

## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a new Pull Request.


## Contact

For any questions or suggestions, please open an issue or contact the maintainer at [nitesh28iitdmaths@gmail.com].