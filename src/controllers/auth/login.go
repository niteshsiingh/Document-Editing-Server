package authcontroller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/niteshsiingh/doc-server/src/config"
	"github.com/niteshsiingh/doc-server/src/database/tables/databases"
	"github.com/niteshsiingh/doc-server/src/entities"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
)

type AuthController struct {
	DB   *databases.Queries
	SMTP *config.SMTP
}

func NewAuthController(db *databases.Queries, smtp *config.SMTP) *AuthController {
	return &AuthController{
		DB:   db,
		SMTP: smtp,
	}
}

func (ac *AuthController) Login(ctx *gin.Context) {
	cxt := context.Background()
	var loginData entities.LoginData
	err := ctx.ShouldBindJSON(&loginData)
	if err != nil {
		responses.NewResponse("Invalid request", http.StatusBadRequest).Send(ctx)
		return
	}

	if !govalidator.IsEmail(loginData.Email) {
		responses.NewResponse("Must provide a valid email address.", http.StatusBadRequest).Send(ctx)
		return
	}

	if loginData.Password == "" {
		responses.NewResponse("Must provide a password.", http.StatusBadRequest).Send(ctx)
		return
	}

	user, err := services.FindUserByEmail(cxt, loginData.Email, ac.DB)
	if err != nil {
		fmt.Println(err)
		responses.NewResponse("Internal server error", http.StatusInternalServerError).Send(ctx)
		return
	}

	validPassword, err := services.CheckPassword(&user, loginData.Password, ac.DB)
	if err != nil {
		responses.NewResponse("Internal server error", http.StatusInternalServerError).Send(ctx)
		return
	}

	if !validPassword {
		responses.NewResponse("Invalid credentials", http.StatusUnauthorized).Send(ctx)
		return
	}

	if !user.IsVerified.Bool {
		responses.NewResponse("User is not verified", http.StatusUnauthorized).Send(ctx)
		return
	}

	msg, err := services.GenerateAuthResponse(cxt, &user, ac.DB)
	if err != nil {
		responses.NewResponse("Internal server error", http.StatusInternalServerError).Send(ctx)
		return
	}

	responses.NewResponse(msg, http.StatusOK).Send(ctx)
}
