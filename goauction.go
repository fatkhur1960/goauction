package main

import (
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/router"
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// generating routes
	utils.GenerateRoutes()
	// connect with database
	models.ConnectDatabase()
	// generate app route
	app := router.GetGeneratedRoutes(gin.Default())
	app.Run("0.0.0.0:8081")
}
