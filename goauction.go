package main

import (
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/router"
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/fatkhur1960/goauction/docs"
	"github.com/gin-gonic/gin"
)

// @title GoAuction API
// @version 1.0
// @description Backend lelah online
// @termsOfService https://github.com/fatkhur1960/goauction
// @license.name GNU
// @license.url https://github.com/fatkhur1960/goauction/blob/master/LICENSE
// @securityDefinitions.apiKey bearerAuth
// @in header
// @name Authorization
func main() {
	docs.SwaggerInfo.Title = "GoAuction - API"
	docs.SwaggerInfo.Description = "GoAuction API Documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = "localhost:8081"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// generating routes
	utils.GenerateRoutes()
	// connect with database
	models.ConnectDatabase()
	// generate app route
	app := router.GetGeneratedRoutes(gin.Default())
	defer models.DB.Close()
	app.Run(docs.SwaggerInfo.Host)
}
