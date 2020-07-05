package service

import (
	"net/http"
	"sync"

	mid "github.com/fatkhur1960/goauction/app/middleware"
	repo "github.com/fatkhur1960/goauction/app/repository"
	"github.com/gin-gonic/gin"
)

type (
	// AuthService for Authentication implementation
	AuthService struct {
		sync.Mutex
		authRepo *repo.AuthRepository
		userRepo *repo.UserRepository
	}

	// AuthQuery definisi query untuk login
	AuthQuery struct {
		Email    string `json:"email" binding:"required"`
		Passhash string `json:"passhash" binding:"required"`
	}
)

// NewAuthService create new instance
// @RouterGroup /auth/v1
func NewAuthService() *AuthService {
	return &AuthService{
		authRepo: repo.NewAuthRepository(),
		userRepo: repo.NewUserRepository(),
	}
}

// AuthorizeUser docs
// @Summary Endpoint untuk melakukan otorisasi
// @Tags AuthService
// @Accept json
// @Produce  json
// @Param email body string true "Email"
// @Param passhash body string true "Passhash"
// @Success 200 {object} app.Result{result=models.AccessToken}
// @Failure 500 {object} app.Result
// @Router /authorize [post]
func (s *AuthService) AuthorizeUser(c *gin.Context, query *AuthQuery) {
	token, err := s.authRepo.AuthorizeUser(query.Email, query.Passhash)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	APIResult.Success(c, token)
}

// UnauthorizeUser docs
// @Summary Endpoint untuk menghapus otorisasi
// @Tags AuthService
// @Security bearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} app.Result
// @Failure 400 {object} app.Result
// @Failure 401 {object} app.Result
// @Router /unauthorize [post] [auth]
func (s *AuthService) UnauthorizeUser(c *gin.Context) {
	err := s.authRepo.UnauthorizeUser(mid.CurrentUser.ID)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	APIResult.Success(c, nil)
}
