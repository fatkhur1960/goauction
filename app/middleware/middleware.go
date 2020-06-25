package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/gin-gonic/gin"
)

// CurrentUser global variable for authenticated user
var (
	CurrentUser models.User
	apiResult   = app.NewAPIResult()
)

type authHeader struct {
	Authorization string `binding:"required"`
}

// RequiresUserAuth middleware
func RequiresUserAuth(c *gin.Context) {
	auth := authHeader{}
	const bearerScheme = "Bearer "
	if err := c.ShouldBindHeader(&auth); err != nil {
		apiResult.Error(c, http.StatusUnauthorized, "Header `Authorization` is not set")
		c.Abort()
		return
	}

	tokenString := strings.ReplaceAll(auth.Authorization, bearerScheme, "")
	accessToken := models.AccessToken{}

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		apiResult.Error(c, http.StatusUnauthorized, "Invalid Access Token")
		c.Abort()
	} else if err := models.NewAccessTokenQuerySet(models.DB).TokenEq(tokenString).One(&accessToken); err != nil {
		apiResult.Error(c, http.StatusUnauthorized, "Unauthorized")
		c.Abort()
	} else if accessToken.IsExpired() {
		apiResult.Error(c, http.StatusUnauthorized, "Access Token Expired")
		c.Abort()
	}

	models.NewUserQuerySet(models.DB).IDEq(accessToken.UserID).One(&CurrentUser)
	c.Next()
}
