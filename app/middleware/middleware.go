package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/gin-gonic/gin"
)

// CurrentUser global variable for authenticated user
var CurrentUser models.User

// AccessToken global variable for authenticated user
var AccessToken models.AccessToken

// RequiresUserAuth middleware
func RequiresUserAuth(c *gin.Context) {
	apiResult := app.NewAPIResult()
	tokenString := c.Request.Header.Get("Authorization")

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
	} else if err := models.NewAccessTokenQuerySet(models.DB).TokenEq(tokenString).One(&AccessToken); err != nil {
		apiResult.Error(c, http.StatusUnauthorized, "Unauthorized")
		c.Abort()
	} else if AccessToken.IsExpired() {
		apiResult.Error(c, http.StatusUnauthorized, "Access Token Expired")
		c.Abort()
	}

	models.NewUserQuerySet(models.DB).IDEq(AccessToken.UserID).One(&CurrentUser)
}
