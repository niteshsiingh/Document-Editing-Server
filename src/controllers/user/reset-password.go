package user

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/niteshsiingh/doc-server/src/entities"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"github.com/niteshsiingh/doc-server/src/responses"
	"github.com/niteshsiingh/doc-server/src/services"
)

func (uc *UserController) ResetPassword(ctx *gin.Context) {
	cxt := context.Background()
	var resetPassword entities.ResetPasswordRequest
	err := ctx.ShouldBindJSON(&resetPassword)
	if err != nil {
		responses.NewResponse("Internal server error", 500).Send(ctx)
		return
	}
	user, err := services.FindUserByEmail(cxt, resetPassword.Email, uc.DB)
	if err != nil {
		responses.NewResponse("Invalid Request", 400).Send(ctx)
	}
	uc.MS.ResetPassword(cxt, &user, uc.DB)
	responses.NewResponse(user, 200).Send(ctx)
}

func (uc *UserController) ConfirmResetPassword(ctx *gin.Context) {
	cxt := context.Background()
	var confirmResetPassword entities.ConfirmResetPasswordRequest
	passwordResetToken := ctx.Param("token")
	err := ctx.ShouldBindJSON(&confirmResetPassword)
	if err != nil {
		responses.NewResponse("Internal server error", 500).Send(ctx)
		return
	}
	email, err := middleware.GetAuth().ParseVerification(passwordResetToken)
	if err != nil {
		responses.NewResponse("Invalid token", 403).Send(ctx)
		return
	}
	user, err := services.FindUserByPasswordResetToken(cxt, email, passwordResetToken, uc.DB)
	if err != nil {
		responses.NewResponse("Internal server error", 500).Send(ctx)
		return
	}
	err = services.UpdatePassword(cxt, user, confirmResetPassword.Password, uc.DB)
	if err != nil {
		responses.NewResponse("Internal server error", 500).Send(ctx)
		return
	}
	responses.NewResponse("Password updated successfully", 200).Send(ctx)

}
