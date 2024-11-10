package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/niteshsiingh/doc-server/src/database/tables/databases"
	"github.com/niteshsiingh/doc-server/src/entities"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
)

type UserController struct {
	DB *databases.Queries
	MS *services.MailService
}

func NewUserController(db *databases.Queries) *UserController {
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
	cxt := context.Background()
	var createUserRequest entities.CreateUserRequest

	err := ctx.ShouldBindJSON(&createUserRequest)
	if err != nil {
		responses.NewResponse("Invalid request", 400).Send(ctx)
		return
	}
	err = uc.MS.CreateUser(cxt, createUserRequest.Email, createUserRequest.Password, uc.DB)
	if err != nil {
		fmt.Println(err)
		responses.NewResponse("Internal server error", 500).Send(ctx)
		return
	}
	responses.NewResponse("User created successfully", 200).Send(ctx)
}

func (uc *UserController) VerifyEmail(ctx *gin.Context) {
	cxt := context.Background()
	verificationToken := ctx.Param("token")
	email, err := middleware.GetAuth().ParseVerification(verificationToken)
	if err != nil {
		responses.NewResponse("Invalid token", 403).Send(ctx)
		return
	}
	user, err := services.FindUserByVerificationToken(cxt, email, verificationToken, uc.DB)
	if err != nil {
		responses.NewResponse("Internal server error", 500).Send(ctx)
		return
	}
	if user.IsVerified.Bool {
		responses.NewResponse("Invalid Request", 400).Send(ctx)
		return
	}
	err = services.UpdateIsVerified(cxt, &user, true, uc.DB)
	if err != nil {
		responses.NewResponse("Internal server error", 500).Send(ctx)
		return
	}
	responses.NewResponse("User Verified Successfully", 200).Send(ctx)
}

func (uc *UserController) GetUserByID(ctx *gin.Context) {
	cxt := context.Background()
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		responses.NewResponse("Internal server error", 500).Send(ctx)
		return
	}
	user, err := services.FindUserByID(cxt, uint(userID), uc.DB)
	if err != nil {
		responses.NewResponse("Invalid Request", 400).Send(ctx)
	}
	responses.NewResponse(user, 200).Send(ctx)
}
