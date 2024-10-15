package services

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgtype"
	dbmodels "github.com/niteshsiingh/doc-server/src/database/db-models"
	"github.com/niteshsiingh/doc-server/src/database/tables/databases"
	"github.com/niteshsiingh/doc-server/src/entities"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"golang.org/x/crypto/bcrypt"
)

func FindUserByEmail(ctx context.Context, email string, db *databases.Queries) (databases.User, error) {
	var user databases.User
	emailpg := pgtype.Text{
		String: email,
		Valid:  true,
	}
	user, err := db.GetUserByEmail(ctx, emailpg)
	if err != nil {
		return databases.User{}, err
	}
	return user, nil
}

func FindUserByID(ctx context.Context, id uint, db *databases.Queries) (databases.User, error) {
	var user databases.User
	user, err := db.GetUserById(ctx, int32(id))
	if err != nil {
		return databases.User{}, err
	}
	return user, nil
}

func FindUserByVerificationToken(ctx context.Context, email string, verificationToken string, db *databases.Queries) (databases.User, error) {
	var user databases.User
	emailpg := pgtype.Text{
		String: email,
		Valid:  true,
	}
	verifypg := pgtype.Text{
		String: verificationToken,
		Valid:  true,
	}
	user, err := db.GetUserByVerificationToken(ctx, databases.GetUserByVerificationTokenParams{
		Email:             emailpg,
		VerificationToken: verifypg,
	})
	if err != nil {
		return databases.User{}, err
	}
	return user, nil
}

func FindUserByPasswordResetToken(ctx context.Context, email string, passwordResetToken string, db *databases.Queries) (databases.User, error) {
	var user databases.User
	emailpg := pgtype.Text{
		String: email,
		Valid:  true,
	}
	passwordpg := pgtype.Text{
		String: passwordResetToken,
		Valid:  true,
	}
	user, err := db.GetUserByPasswordResetToken(ctx, databases.GetUserByPasswordResetTokenParams{
		Email:              emailpg,
		PasswordResetToken: passwordpg,
	})
	if err != nil {
		return databases.User{}, err
	}
	return user, nil
}

func (ms *MailService) CreateUser(ctx context.Context, email string, password string, db *databases.Queries) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	verificationToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})
	verificationTokenString, err := verificationToken.SignedString([]byte(os.Getenv("VERIFY_EMAIL_SECRET")))
	if err != nil {
		return err
	}
	user := databases.User{
		Email:             pgtype.Text{String: email, Valid: true},
		Password:          pgtype.Text{String: string(hashedPassword), Valid: true},
		VerificationToken: pgtype.Text{String: verificationTokenString, Valid: true},
	}
	err = db.CreateUser(ctx, databases.CreateUserParams{
		Email:             pgtype.Text{String: email, Valid: true},
		Password:          pgtype.Text{String: string(hashedPassword), Valid: true},
		VerificationToken: pgtype.Text{String: verificationTokenString, Valid: true},
	})
	if err != nil {
		return err
	}
	err = ms.sendVerificationEmail(&user, db)
	if err != nil {
		return err
	}
	return nil
}

func CheckPassword(user *databases.User, password string, db *databases.Queries) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func GenerateAuthResponse(ctx context.Context, user *databases.User, db *databases.Queries) (entities.TokenPair, error) {
	requestAccessTokenClaim, err := getUserClaims(ctx, user, "ACCESS_TOKEN_EXPIRATION", db)
	if err != nil {
		return entities.TokenPair{}, err
	}
	requestRefreshTokenClaim, err := getUserClaims(ctx, user, "REFRESH_TOKEN_EXPIRATION", db)
	if err != nil {
		return entities.TokenPair{}, err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, requestAccessTokenClaim)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return entities.TokenPair{}, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, requestRefreshTokenClaim)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return entities.TokenPair{}, err
	}
	userID := pgtype.Int4{
		Int32: int32(requestRefreshTokenClaim.UserID),
		Valid: true,
	}
	err = db.DeleteRefreshToken(ctx, userID)
	if err != nil {
		return entities.TokenPair{}, err
	}
	err = db.CreateRefreshToken(ctx, databases.CreateRefreshTokenParams{
		Token: pgtype.Text{String: refreshTokenString, Valid: true},
		UserID: pgtype.Int4{
			Int32: int32(requestRefreshTokenClaim.UserID),
			Valid: true,
		},
	})
	if err != nil {
		return entities.TokenPair{}, err
	}

	return entities.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func GetIsTokenActive(ctx context.Context, token string, db *databases.Queries) (bool, error) {
	_, err := db.GetRefreshToken(ctx, pgtype.Text{String: token, Valid: true})
	if err != nil {
		return false, err
	}
	return true, nil
}

func LogoutUser(ctx context.Context, userID uint, db *databases.Queries) error {
	userId := pgtype.Int4{
		Int32: int32(userID),
		Valid: true,
	}
	err := db.DeleteRefreshToken(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MailService) ResetPassword(ctx context.Context, user *databases.User, db *databases.Queries) error {
	refreshExpiry, err := strconv.Atoi(os.Getenv("REFRESH_VALIDITY"))
	if err != nil {
		return err
	}
	expiryTime := time.Duration(refreshExpiry) * time.Second
	expirationTime := time.Now().Add(expiryTime)
	claims := &middleware.Claims{
		UserID:  uint(user.ID),
		EmailID: user.Email.String,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	passwordResetToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := passwordResetToken.SignedString([]byte(os.Getenv("PASSWORD_RESET_SECRET")))
	if err != nil {
		return err
	}

	user.PasswordResetToken = pgtype.Text{
		String: tokenString,
		Valid:  true,
	}
	err = db.EditUser(ctx, databases.EditUserParams{
		Email:              user.Email,
		Password:           user.Password,
		IsVerified:         user.IsVerified,
		VerificationToken:  user.VerificationToken,
		PasswordResetToken: user.PasswordResetToken,
		UserRoles:          user.UserRoles,
		SharedDocuments:    user.SharedDocuments,
		Documents:          user.Documents,
		RefreshTokens:      user.RefreshTokens,
		ID:                 user.ID,
	})
	if err != nil {
		return err
	}

	err = ms.sendPasswordResetEmail(user, db)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePassword(ctx context.Context, user databases.User, password string, db *databases.Queries) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = pgtype.Text{
		String: string(hashedPassword),
		Valid:  true,
	}
	err = db.EditUser(ctx, databases.EditUserParams{
		Email:              user.Email,
		Password:           user.Password,
		IsVerified:         user.IsVerified,
		VerificationToken:  user.VerificationToken,
		PasswordResetToken: user.PasswordResetToken,
		UserRoles:          user.UserRoles,
		SharedDocuments:    user.SharedDocuments,
		Documents:          user.Documents,
		RefreshTokens:      user.RefreshTokens,
		ID:                 user.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func UpdateIsVerified(ctx context.Context, user *databases.User, isVerified bool, db *databases.Queries) error {
	user.IsVerified = pgtype.Bool{
		Bool:  isVerified,
		Valid: true,
	}
	err := db.EditUser(ctx, databases.EditUserParams{
		Email:              user.Email,
		Password:           user.Password,
		IsVerified:         user.IsVerified,
		VerificationToken:  user.VerificationToken,
		PasswordResetToken: user.PasswordResetToken,
		UserRoles:          user.UserRoles,
		SharedDocuments:    user.SharedDocuments,
		Documents:          user.Documents,
		RefreshTokens:      user.RefreshTokens,
		ID:                 user.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (ms *MailService) sendPasswordResetEmail(user *databases.User, db *databases.Queries) error {
	err := ms.SendMail(MailOptions{
		From:    os.Getenv("SMTP_USER"),
		To:      []string{user.Email.String},
		Subject: "Reset your password!",
		Body:    os.Getenv("FRONT_END_URL") + "/user/reset-email/" + user.PasswordResetToken.String,
	})
	if err != nil {
		return err
	}
	return nil
}

func (ms *MailService) sendVerificationEmail(user *databases.User, db *databases.Queries) error {
	err := ms.SendMail(MailOptions{
		From:    os.Getenv("SMTP_USER"),
		To:      []string{user.Email.String},
		Subject: "Welcome to Docs!",
		Body:    "Click the following link to verify your email: " + os.Getenv("FRONT_END_URL") + "/user/verify-email/" + user.VerificationToken.String,
	})
	if err != nil {
		return err
	}
	return nil
}

func getUserClaims(ctx context.Context, user *databases.User, expiryPeriod string, db *databases.Queries) (*middleware.Claims, error) {
	var roles []string
	var userRoles []dbmodels.UserRole
	var userRolesMap map[int]int32
	if len(user.UserRoles) != 0 {
		err := json.Unmarshal(user.UserRoles, &userRolesMap)
		if err != nil {
			return nil, err
		}
	}
	for _, userRoleId := range userRolesMap {
		role, err := db.GetUserRole(ctx, userRoleId)
		if err != nil {
			continue
		}
		userRoles = append(userRoles, dbmodels.UserRole{
			UserID: uint(userRoleId),
			RoleID: uint(role.ID),
			Role:   role.Role.String,
		})
		roles = append(roles, role.Role.String)
	}
	expiryStr, err := strconv.Atoi(os.Getenv(expiryPeriod))
	if err != nil {
		return nil, err
	}
	expiryTime := time.Duration(expiryStr) * time.Second
	expirationTime := time.Now().Add(expiryTime)
	return &middleware.Claims{
		User: entities.JWTUser{
			ID:      uint(user.ID),
			EmailID: user.Email.String,
			// Documents: user.GetDocumentPermissions(),
			Roles: userRoles,
		},
		UserID:  uint(user.ID),
		EmailID: user.Email.String,
		Roles:   roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}, nil
}
