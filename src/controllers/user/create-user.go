package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/niteshsiingh/doc-server/src/entities"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
	MS *services.MailService
}

func NewUserController(db *gorm.DB) *UserController {
	ms, err := services.NewMailService()
	if err != nil {
		return nil
	}
	t := &UserController{
		DB: db,
		MS: ms,
	}
	return t
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var createUserRequest entities.CreateUserRequest

	err := ctx.ShouldBindJSON(&createUserRequest)
	if err != nil {
		responses.NewResponse("Invalid request", 400).Send(ctx)
		return
	}
	err = uc.MS.CreateUser(createUserRequest.Email, createUserRequest.Password, uc.DB)
	if err != nil {
		responses.NewResponse("Internal server error", 500).Send(ctx)
		return
	}
	responses.NewResponse("User created successfully", 200).Send(ctx)
}

func (uc *UserController) VerifyEmail(ctx *gin.Context) {
	verificationToken := ctx.Param("token")
	email, err := middleware.GetAuth().ParseVerification(verificationToken)
	if err != nil {
		responses.NewResponse("Invalid token", 403).Send(ctx)
		return
	}
	user, err := services.FindUserByVerificationToken(email, verificationToken, uc.DB)
	if err != nil {
		responses.NewResponse("Internal server error", 500).Send(ctx)
		return
	}
	if user.IsVerified {
		responses.NewResponse("Invalid Request", 400).Send(ctx)
		return
	}
	err = services.UpdateIsVerified(&user, true, uc.DB)
	if err != nil {
		responses.NewResponse("Internal server error", 500).Send(ctx)
		return
	}
	responses.NewResponse("User Verified Successfully", 200).Send(ctx)
}

func (uc *UserController) GetUserByID(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		responses.NewResponse("Internal server error", 500).Send(ctx)
		return
	}
	user, err := services.FindUserByID(uint(userID), uc.DB)
	if err != nil {
		responses.NewResponse("Invalid Request", 400).Send(ctx)
	}
	responses.NewResponse(user, 200).Send(ctx)
}
