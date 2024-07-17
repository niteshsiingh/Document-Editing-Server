package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/niteshsiingh/doc-server/src/config"
	authcontroller "github.com/niteshsiingh/doc-server/src/controllers/auth"
	"github.com/niteshsiingh/doc-server/src/controllers/document"
	"github.com/niteshsiingh/doc-server/src/controllers/user"
	dbmodels "github.com/niteshsiingh/doc-server/src/database/db-models"
	"github.com/niteshsiingh/doc-server/src/middleware"
	"gorm.io/gorm"
)

func createRouter(db *gorm.DB, smtp *config.SMTP) *gin.Engine {
	db.AutoMigrate(&dbmodels.User{}, &dbmodels.Document{}, &dbmodels.DocumentUser{}, &dbmodels.RefreshToken{})
	ac := authcontroller.NewAuthController(db, smtp)
	uc := user.NewUserController(db)
	dc := document.NewDocumentController(db)
	fmt.Println("uc: ", uc)
	fmt.Println(ac, uc, dc, db)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())
	routes := router.Group("/v1")
	nonAuthRoutes := router.Group("/v1")
	routes.Use(middleware.AuthMiddleware())
	{
		nonAuthRoutes.POST("/auth/login", ac.Login)
		nonAuthRoutes.POST("/auth/refresh-token", ac.RefreshToken)
		nonAuthRoutes.DELETE("/auth/logout", ac.Logout)

		routes.GET("/user/:id", uc.GetUserByID)
		nonAuthRoutes.PUT("/user/verify-email/:token", uc.VerifyEmail)
		routes.PUT("/user/password/:token", uc.ConfirmResetPassword)
		routes.POST("/user/reset-password", uc.ResetPassword)
		nonAuthRoutes.POST("/user", uc.CreateUser)

		routes.GET("/document", dc.GetAllDocuments)
		routes.GET("/document/:document_id", dc.GetOneDocument)
		routes.PUT("/document/:document_id", dc.UpdateDocument)
		routes.POST("/document", dc.CreateDocument)
		routes.POST("/document/:document_id/share", dc.ShareDocument)
		routes.DELETE("/document/:document_id", dc.DeleteDocument)
		routes.DELETE("/document/:document_id/share", dc.RemoveSharedUser)
	}

	return router
}
