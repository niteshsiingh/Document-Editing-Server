package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/niteshsiingh/doc-server/src/entities"
)

type Auth struct {
	SecretKey     []byte
	Expiry        time.Duration
	RefreshExpiry time.Duration
}

var authInstance = &Auth{}

func Init(jwtKey string, expiry int, refreshExpiry int) {
	authInstance.Expiry = time.Duration(expiry) * time.Second
	authInstance.RefreshExpiry = time.Duration(refreshExpiry) * time.Second
	authInstance.SecretKey = []byte(jwtKey)
}

func GetAuth() *Auth {
	if len(authInstance.SecretKey) == 0 {
		panic("JWT not initialized")
	}
	return authInstance
}

type Claims struct {
	User    entities.JWTUser `json:"user,omitempty"`
	UserID  uint             `json:"userId,omitempty"`
	EmailID string           `json:"emailId,omitempty"`
	Roles   []string         `json:"roles,omitempty"`
	jwt.StandardClaims
}

func (a *Auth) ParseAuth(tokenString string, secretName string) (entities.JWTUser, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return a.SecretKey, nil
	})

	if err != nil {
		return entities.JWTUser{}, err
	}

	if !token.Valid {
		return entities.JWTUser{}, errors.New("invalid JWT token")
	}

	claims := token.Claims.(jwt.MapClaims)
	var user entities.JWTUser
	dbByte, err := json.Marshal(claims["user"])
	if err != nil {
		return entities.JWTUser{}, err
	}
	err = json.Unmarshal(dbByte, &user)
	if err != nil {
		return entities.JWTUser{}, err
	}

	return user, nil
}

func (a *Auth) ParseVerification(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return []byte(os.Getenv("VERIFY_EMAIL_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid JWT token")
	}

	claims := token.Claims.(jwt.MapClaims)
	emailID := claims["email"].(string)
	return emailID, nil
}

func (a *Auth) ParseToken(tokenString string) (*jwt.Token, *entities.ParsedToken, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return a.SecretKey, nil
	})
	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, errors.New("invalid JWT token")
	}
	claims := token.Claims.(jwt.MapClaims)
	jwtUserId := uint(claims["userId"].(float64))
	jwtUserEmailId := claims["emailId"].(string)
	parsedToken := entities.ParsedToken{
		UserID:  jwtUserId,
		EmailID: jwtUserEmailId,
		User:    entities.JWTUser{},
	}

	return token, &parsedToken, err
}
