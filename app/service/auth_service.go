package service

import (
	"net/http"
	"time"

	mid "github.com/fatkhur1960/goauction/app/middleware"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type (
	// AuthService for authentication
	AuthService interface {
		AuthorizeUser(c *gin.Context)
		UnauthorizeUser(c *gin.Context)
	}

	// AuthServiceImpl for AuthService implementation
	AuthServiceImpl struct {
		userRepo     models.UserQuerySet
		tokenRepo    models.AccessTokenQuerySet
		passhashRepo models.UserPasshashQuerySet
	}

	// AuthQuery definisi query untuk login
	AuthQuery struct {
		Email    string `json:"email" binding:"required"`
		Passhash string `json:"passhash" binding:"required"`
	}
)

// NewAuthService create new instance
// @RouterGroup /auth/v1
func NewAuthService(db *gorm.DB) AuthService {
	return &AuthServiceImpl{
		userRepo:     models.NewUserQuerySet(db),
		tokenRepo:    models.NewAccessTokenQuerySet(db),
		passhashRepo: models.NewUserPasshashQuerySet(db),
	}
}

// AuthorizeUser docs
// @Summary Endpoint untuk melakukan otorisasi
// @Tags AuthService
// @Accept json
// @Produce  json
// @Param email body string true "Email"
// @Param passhash body string true "Passhash"
// @Success 200 {object} models.AccessToken
// @Failure 500 {object} app.Result
// @Router /authorize [post]
func (s *AuthServiceImpl) AuthorizeUser(c *gin.Context) {
	var query AuthQuery
	var user models.User
	var userPasshash models.UserPasshash

	if err := c.ShouldBindJSON(&query); err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// check apakah user ada di db
	err := s.userRepo.GetDB().First(&user, "email=? AND active=?", &query.Email, true).Related(&userPasshash).Error
	if err != nil {
		APIResult.Error(c, http.StatusUnauthorized, "Email tidak ditemukan")
		return
	}

	// check passhash
	if !utils.CheckPasshash(query.Passhash, userPasshash.Passhash) {
		APIResult.Error(c, http.StatusUnauthorized, "Password tidak cocok")
		return
	}

	// generate access token
	token, expireTime, _ := utils.GenerateToken(user.Email)
	accessToken := models.AccessToken{
		User:      user,
		Token:     token,
		Created:   time.Now().UTC(),
		ValidThru: expireTime,
	}
	resToken, _ := accessToken.CreateAccessToken()

	APIResult.Success(c, resToken)
}

// UnauthorizeUser docs
// @Summary Endpoint untuk menghapus otorisasi
// @Tags AuthService
// @Security bearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} app.Result
// @Failure 400 {object} app.Result
// @Router /unauthorize [post] [auth]
func (s *AuthServiceImpl) UnauthorizeUser(c *gin.Context) {
	err := s.tokenRepo.UserIDEq(mid.CurrentUser.ID).Delete()
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	APIResult.Success(c, nil)
}
