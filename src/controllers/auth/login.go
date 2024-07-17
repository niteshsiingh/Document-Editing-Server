package authcontroller

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/niteshsiingh/doc-server/src/config"
	"github.com/niteshsiingh/doc-server/src/entities"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
	"gorm.io/gorm"
)

type AuthController struct {
	DB   *gorm.DB
	SMTP *config.SMTP
}

func NewAuthController(db *gorm.DB, smtp *config.SMTP) *AuthController {
	return &AuthController{
		DB:   db,
		SMTP: smtp,
	}
}

func (ac *AuthController) Login(ctx *gin.Context) {

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

	user, err := services.FindUserByEmail(loginData.Email, ac.DB)
	if err != nil {
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

	if !user.IsVerified {
		responses.NewResponse("User is not verified", http.StatusUnauthorized).Send(ctx)
		return
	}

	msg, err := services.GenerateAuthResponse(&user, ac.DB)
	if err != nil {
		responses.NewResponse("Internal server error", http.StatusInternalServerError).Send(ctx)
		return
	}

	responses.NewResponse(msg, http.StatusOK).Send(ctx)
}
