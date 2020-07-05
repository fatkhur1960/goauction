package middleware

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/repository"
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// CurrentUser global variable for authenticated user
var (
	CurrentUser models.User
	apiResult   = app.NewAPIResult()
)

type authHeader struct {
	sync.Mutex
	Authorization string `binding:"required"`
}

// ReqValidate for handle request with params/json validator
func ReqValidate(c *gin.Context, query interface{}, bindType binding.Binding) (interface{}, error) {
	if err := c.ShouldBindWith(query, bindType); err != nil {
		apiResult.Error(c, http.StatusBadRequest, utils.ParseError(err.Error()))
		c.Abort()
		return nil, err
	}

	return query, nil
}

// MethodValidator for handling request method
func MethodValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		re := regexp.MustCompile("GET|POST")
		if !re.MatchString(c.Request.Method) {
			apiResult.Error(c, http.StatusMethodNotAllowed, "Request method not allowed")
			return
		}

		c.Next()
	}
}

// RequiresUserAuth middleware
func RequiresUserAuth(c *gin.Context) {
	auth := authHeader{}
	authRepo := repository.NewAuthRepository()
	userRepo := repository.NewUserRepository()
	const bearerScheme = "Bearer "
	if err := c.ShouldBindHeader(&auth); err != nil {
		apiResult.Error(c, http.StatusUnauthorized, "Header `Authorization` is not set")
		c.Abort()
		return
	}

	tokenString := strings.ReplaceAll(auth.Authorization, bearerScheme, "")
	accessToken, atErr := authRepo.GetAccessToken(tokenString)

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
	} else if atErr != nil {
		apiResult.Error(c, http.StatusUnauthorized, "Unauthorized")
		c.Abort()
	} else if accessToken.IsExpired() {
		apiResult.Error(c, http.StatusUnauthorized, "Access Token Expired")
		c.Abort()
	}

	CurrentUser, _ = userRepo.GetByID(accessToken.UserID)
	c.Next()
}
