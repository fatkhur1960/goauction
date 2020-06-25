package main

import (
	"fmt"
	"log"

	"github.com/fatkhur1960/goauction/app/event"
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

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, _ int) {
		log.SetPrefix("[")
		log.Printf("%v] + endpoint %v %v\n", utils.ReplacePackages(handlerName), httpMethod, absolutePath)
	}

	// connect with database
	models.ConnectDatabase()

	// generating routes
	if err := utils.GenerateRoutes(); err == nil {
		// generate app route
		app := router.GetGeneratedRoutes(gin.New())

		// register event listener
		app.Use(event.RegisterEvents)

		defer models.DB.Close()
		// emmit startup event
		event.Listener.Emmit(&event.StartupEvent{})
		app.Run(docs.SwaggerInfo.Host)
	} else {
		fmt.Println("RouteGenerator error: " + err.Error())
	}
}
