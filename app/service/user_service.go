package service

import (
	"net/http"

	"github.com/fatkhur1960/goauction/app/event"
	mid "github.com/fatkhur1960/goauction/app/middleware"
	"github.com/fatkhur1960/goauction/app/models"
	repo "github.com/fatkhur1960/goauction/app/repository"
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type (
	// UserService implementation for users
	UserService struct {
		userRepo repo.UserRepository
	}

	// RegisterUserQuery definisi query untuk register user
	RegisterUserQuery struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		PhoneNum string `json:"phone_num" binding:"required"`
	}

	// ActivateUserQuery definisi query untuk mengaktifkan registered user
	ActivateUserQuery struct {
		Token    string `json:"token" binding:"required"`
		Passhash string `json:"passhash" binding:"required"`
	}
)

// NewUserService instance for UserService
// @RouterGroup /user/v1
func NewUserService(db *gorm.DB) UserService {
	userRepository := repo.UserRepository{
		UserQs:     models.NewUserQuerySet(db),
		RegisterQs: models.NewRegisterUserQuerySet(db),
		PasshashQs: models.NewUserPasshashQuerySet(db),
	}

	return UserService{
		userRepo: userRepository,
	}
}

// RegisterUser docs
// @Tags UserService
// @Summary Endpoint untuk register user
// @Accept json
// @Produce json
// @Param full_name body string true "FullName"
// @Param email body string true "Email"
// @Param phone_num body string true "PhoneNum"
// @Success 200 {object} app.Result{result=models.RegisterUser}
// @Failure 400 {object} app.Result
// @Router /register [post]
func (s *UserService) RegisterUser(c *gin.Context) {
	query := RegisterUserQuery{}
	validateRequest(c, &query)

	token, _, _ := utils.GenerateToken(query.Email)
	user, err := s.userRepo.RegisterUser(
		query.FullName,
		query.Email,
		query.PhoneNum,
		token,
	)

	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Emmit register event
	event.Listener.Emmit(&event.UserRegisteredPayload{
		FullName: user.FullName,
		Email:    user.Email,
		PhoneNum: user.PhoneNum,
		Token:    user.Token,
	})

	// return token-nya apabila email sudah terdaftar
	APIResult.Success(c, &user)
}

// ActivateUser docs
// @Tags UserService
// @Summary Endpoint untuk mengaktifkan user
// @Accept json
// @Produce json
// @Param token body string true "Token"
// @Param passhash body string true "Passhash"
// @Success 200 {object} app.Result{result=models.User}
// @Failure 400 {object} app.Result
// @Router /activate [post]
func (s *UserService) ActivateUser(c *gin.Context) {
	query := ActivateUserQuery{}
	validateRequest(c, &query)

	user, err := s.userRepo.ActivateUser(query.Token, query.Passhash)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	APIResult.Success(c, &user)
}

// MeInfo docs
// @Tags UserService
// @Summary Endpoint untuk informasi user
// @Security bearerAuth
// @Produce json
// @Success 200 {object} app.Result{result=models.User}
// @Failure 401 {object} app.Result
// @Router /me/info [get] [auth]
func (s *UserService) MeInfo(c *gin.Context) {
	APIResult.Success(c, mid.CurrentUser)
}

// UpdateUserInfo docs
// @Tags UserService
// @Summary Endpoint untuk mengupdate informasi user
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param full_name body string true "FullName"
// @Param email body string true "Email"
// @Param phone_num body string true "PhoneNum"
// @Param address body string false "Address"
// @Param avatar body string false "Avatar"
// @Success 200 {object} app.Result{result=models.User}
// @Failure 400 {object} app.Result
// @Failure 401 {object} app.Result
// @Router /me/info [post] [auth]
func (s *UserService) UpdateUserInfo(c *gin.Context) {
	query := repo.UpdateUserQuery{}
	validateRequest(c, &query)

	user, err := s.userRepo.UpdateUser(query)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	APIResult.Success(c, &user)
}
