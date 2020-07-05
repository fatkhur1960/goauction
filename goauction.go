package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatkhur1960/goauction/app"
	mid "github.com/fatkhur1960/goauction/app/middleware"
	"github.com/fatkhur1960/goauction/app/router"
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/fatkhur1960/goauction/docs"
	"github.com/fatkhur1960/goauction/system/monitor"
	"github.com/fatkhur1960/goauction/system/queue"
	"github.com/fatkhur1960/goauction/system/socket"
	"github.com/gin-contrib/cors"
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
	log.SetPrefix("[")
	// generating routes
	log.Println("RouteGenerator] generating routes...")
	go utils.GenerateRoutes()

	docs.SwaggerInfo.Title = "GoAuction - API"
	docs.SwaggerInfo.Description = "GoAuction API Documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf(":%v", os.Getenv("PORT"))
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// connect with database
	app.ConnectDatabase()
	defer app.CloseDatabase()

	QueueDispatcher := queue.NewDispatcher(4)
	QueueDispatcher.Run()
	go monitor.StartMonitors()

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, _ int) {
		log.Printf("%v] + endpoint %v %v\n", utils.ReplacePackages(handlerName), httpMethod, absolutePath)
	}

	wsHandler := socket.Handler()
	go wsHandler.Serve()
	defer wsHandler.Close()

	goauction := router.GetGeneratedRoutes(gin.Default())
	goauction.Use(mid.MethodValidator())
	goauction.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	goauction.GET("/socket.io/*any", mid.RequiresUserAuth, gin.WrapH(wsHandler))
	goauction.POST("/socket.io/*any", mid.RequiresUserAuth, gin.WrapH(wsHandler))
	goauction.Handle("WS", "/socket.io/*any", mid.RequiresUserAuth, gin.WrapH(wsHandler))
	goauction.Handle("WSS", "/socket.io/*any", mid.RequiresUserAuth, gin.WrapH(wsHandler))
	// goauction.Run(docs.SwaggerInfo.Host)

	srv := &http.Server{
		Handler:      goauction,
		Addr:         docs.SwaggerInfo.Host,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("\nListening on", docs.SwaggerInfo.Host)
	log.Fatal(srv.ListenAndServe())
}
