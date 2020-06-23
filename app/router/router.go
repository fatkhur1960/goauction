package router

import (
	mid "github.com/fatkhur1960/goauction/app/middleware"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/service"
	"github.com/gin-gonic/gin"
)

// GetGeneratedRoutes list route pada app
func GetGeneratedRoutes(router *gin.Engine) *gin.Engine {
	apiGroup := router.Group("/api")
	{
		// Do not edit this code by your hand
		// this code generate automatically when program running
		// @StartCodeBlocks

		// Generate route for AuthService
		authService := service.NewAuthService(models.DB)
		authServiceGroup := apiGroup.Group("/auth/v1")
		{
			authServiceGroup.POST("/authorize", authService.AuthorizeUser)
			authServiceGroup.POST("/unauthorize", mid.RequiresUserAuth, authService.UnauthorizeUser)
		}

		// Generate route for UserService
		userService := service.NewUserService(models.DB)
		userServiceGroup := apiGroup.Group("/user/v1")
		{
			userServiceGroup.POST("/register", userService.RegisterUser)
			userServiceGroup.POST("/activate", userService.ActivateUser)
			userServiceGroup.GET("/me/info", mid.RequiresUserAuth, userService.MeInfo)
			userServiceGroup.POST("/me/info", mid.RequiresUserAuth, userService.UpdateUserInfo)
		}

		// @EndCodeBlocks
	}

	router.NoRoute(service.NoRouteHandler)

	return router
}
