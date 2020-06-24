package service

import (
	"net/http"

	"github.com/fatkhur1960/goauction/app"
	"github.com/gin-gonic/gin"
)

type (
	// RegisterToken result after user register
	RegisterToken struct {
		Token string `json:"token"`
	}
)

// APIResult for json output
var APIResult = app.NewAPIResult()

// NoRouteHandler handle 404
func NoRouteHandler(c *gin.Context) {
	APIResult.Error(c, http.StatusNotFound, "Endpoint not found")
}
