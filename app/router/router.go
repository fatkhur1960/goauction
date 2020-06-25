package router

import (
	mid "github.com/fatkhur1960/goauction/app/middleware"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/service"

	// import swagger doc
	_ "github.com/fatkhur1960/goauction/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
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

		// Generate route for ProductService
		productService := service.NewProductService(models.DB)
		productServiceGroup := apiGroup.Group("/product/v1")
		{
			productServiceGroup.POST("/add", mid.RequiresUserAuth, productService.AddProduct)
			productServiceGroup.GET("/list", mid.RequiresUserAuth, productService.ListProduct)
			productServiceGroup.GET("/detail/:id", mid.RequiresUserAuth, productService.DetailProduct)
			productServiceGroup.POST("/update/:id", mid.RequiresUserAuth, productService.UpdateProduct)
			productServiceGroup.POST("/delete/:id", mid.RequiresUserAuth, productService.DeleteProduct)
			productServiceGroup.POST("/bid/:id", mid.RequiresUserAuth, productService.BidProduct)
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
	router.Use(gin.Recovery())

	return router
}
