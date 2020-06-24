package router

import (
	mid "github.com/fatkhur1960/goauction/app/middleware"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/service"
	"github.com/fatkhur1960/goauction/docs"

	// import swagger doc
	_ "github.com/fatkhur1960/goauction/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// GetGeneratedRoutes list route pada app
func GetGeneratedRoutes(router *gin.Engine) *gin.Engine {
	apiGroup := router.Group(docs.SwaggerInfo.BasePath)
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

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.NoRoute(service.NoRouteHandler)

	return router
}
