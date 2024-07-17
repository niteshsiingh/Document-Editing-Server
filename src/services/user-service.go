package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	dbmodels "github.com/niteshsiingh/doc-server/src/database/db-models"
	"github.com/niteshsiingh/doc-server/src/entities"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func FindUserByEmail(email string, db *gorm.DB) (dbmodels.User, error) {
	var user dbmodels.User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return dbmodels.User{}, err
	}
	return user, nil
}

func FindUserByID(id uint, db *gorm.DB) (dbmodels.User, error) {
	var user dbmodels.User
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return dbmodels.User{}, err
	}
	return user, nil
}

func FindUserByVerificationToken(email string, verificationToken string, db *gorm.DB) (dbmodels.User, error) {
	var user dbmodels.User
	err := db.Where(&dbmodels.User{VerificationToken: verificationToken, Email: email}).First(&user).Error
	if err != nil {
		return dbmodels.User{}, err
	}
	return user, nil
}

func FindUserByPasswordResetToken(email string, passwordResetToken string, db *gorm.DB) (dbmodels.User, error) {
	var user dbmodels.User
	err := db.Where(&dbmodels.User{PasswordResetToken: passwordResetToken, Email: email}).First(&user).Error
	if err != nil {
		return dbmodels.User{}, err
	}
	return user, nil
}

func (ms *MailService) CreateUser(email string, password string, db *gorm.DB) error {
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
	user := dbmodels.User{
		Email:             email,
		Password:          string(hashedPassword),
		VerificationToken: verificationTokenString,
	}
	err = db.Create(&user).Error
	if err != nil {
		return err
	}
	err = ms.sendVerificationEmail(&user, db)
	if err != nil {
		return err
	}
	return nil
}

func CheckPassword(user *dbmodels.User, password string, db *gorm.DB) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func GenerateAuthResponse(user *dbmodels.User, db *gorm.DB) (entities.TokenPair, error) {
	requestAccessTokenClaim, err := getUserClaims(*user, "ACCESS_TOKEN_EXPIRATION")
	if err != nil {
		return entities.TokenPair{}, err
	}
	requestRefreshTokenClaim, err := getUserClaims(*user, "REFRESH_TOKEN_EXPIRATION")
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
	err = db.Where("user_id = ?", requestRefreshTokenClaim.UserID).Delete(&dbmodels.RefreshToken{}).Error
	if err != nil {
		return entities.TokenPair{}, err
	}

	refreshTokenDB := dbmodels.RefreshToken{
		Token:  refreshTokenString,
		UserID: requestRefreshTokenClaim.UserID,
		User:   *user,
	}
	err = db.Create(&refreshTokenDB).Error
	if err != nil {
		return entities.TokenPair{}, err
	}

	return entities.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func GetIsTokenActive(token string, db *gorm.DB) (bool, error) {
	var refreshToken dbmodels.RefreshToken
	err := db.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func LogoutUser(userID uint, db *gorm.DB) error {
	err := db.Where("user_id = ?", userID).Delete(&dbmodels.RefreshToken{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (ms *MailService) ResetPassword(user *dbmodels.User, db *gorm.DB) error {
	refreshExpiry, err := strconv.Atoi(os.Getenv("REFRESH_VALIDITY"))
	if err != nil {
		return err
	}
	expiryTime := time.Duration(refreshExpiry) * time.Second
	expirationTime := time.Now().Add(expiryTime)
	claims := &middleware.Claims{
		UserID:  user.ID,
		EmailID: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	passwordResetToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := passwordResetToken.SignedString([]byte(os.Getenv("PASSWORD_RESET_SECRET")))
	if err != nil {
		return err
	}

	user.PasswordResetToken = tokenString
	err = db.Save(&user).Error
	if err != nil {
		return err
	}

	err = ms.sendPasswordResetEmail(user, db)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePassword(user dbmodels.User, password string, db *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	err = db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateIsVerified(user *dbmodels.User, isVerified bool, db *gorm.DB) error {
	user.IsVerified = isVerified
	err := db.Save(&user).Error
	if err != nil {
		return err
	}
	fmt.Println("user ", user)
	var savedUser dbmodels.User
	err = db.Where("id = ?", user.ID).First(&savedUser).Error
	if err != nil {
		fmt.Println("not found")
	}
	fmt.Println(savedUser)
	return nil
}

func (ms *MailService) sendPasswordResetEmail(user *dbmodels.User, db *gorm.DB) error {
	err := ms.SendMail(MailOptions{
		From:    os.Getenv("SMTP_USER"),
		To:      []string{user.Email},
		Subject: "Reset your password!",
		Body:    os.Getenv("FRONT_END_URL") + "/user/reset-email/" + user.PasswordResetToken,
	})
	if err != nil {
		return err
	}
	return nil
}

func (ms *MailService) sendVerificationEmail(user *dbmodels.User, db *gorm.DB) error {
	err := ms.SendMail(MailOptions{
		From:    os.Getenv("SMTP_USER"),
		To:      []string{user.Email},
		Subject: "Welcome to Docs!",
		Body:    "Click the following link to verify your email: " + os.Getenv("FRONT_END_URL") + "/user/verify-email/" + user.VerificationToken,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func getUserClaims(user dbmodels.User, expiryPeriod string) (*middleware.Claims, error) {
	var roles []string
	for _, userRole := range user.UserRoles {
		roles = append(roles, string(userRole.Role.Name))
	}
	expiryStr, err := strconv.Atoi(os.Getenv(expiryPeriod))
	if err != nil {
		return nil, err
	}
	expiryTime := time.Duration(expiryStr) * time.Second
	expirationTime := time.Now().Add(expiryTime)
	return &middleware.Claims{
		User: entities.JWTUser{
			ID:        user.ID,
			EmailID:   user.Email,
			Documents: user.GetDocumentPermissions(),
			Roles:     user.UserRoles,
		},
		UserID:  user.ID,
		EmailID: user.Email,
		Roles:   roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}, nil
}
