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

// RequestValidator for params/uri/json validator
func RequestValidator(c *gin.Context, query interface{}) {
	if err := c.ShouldBindJSON(query); err != nil {
		apiResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	c.Set("validated", query)
	c.Next()
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
	} else if err := models.NewAccessTokenQuerySet(app.DB).TokenEq(tokenString).One(&accessToken); err != nil {
		apiResult.Error(c, http.StatusUnauthorized, "Unauthorized")
		c.Abort()
	} else if accessToken.IsExpired() {
		apiResult.Error(c, http.StatusUnauthorized, "Access Token Expired")
		c.Abort()
	}

	CurrentUser = models.User{}
	models.NewUserQuerySet(app.DB).IDEq(accessToken.UserID).One(&CurrentUser)
	c.Next()
}
