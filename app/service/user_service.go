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
	// UserService service for user
	UserService interface {
		RegisterUser(c *gin.Context)
		ActivateUser(c *gin.Context)
		MeInfo(c *gin.Context)
		UpdateUserInfo(c *gin.Context)
	}

	// UserServiceImpl struct implementation for user service
	UserServiceImpl struct {
		userRepo     models.UserQuerySet
		registerRepo models.RegisterUserQuerySet
		passhashRepo models.UserPasshashQuerySet
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

	// UpdateUserQuery definisi query untuk update user
	UpdateUserQuery struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		PhoneNum string `json:"phone_num" binding:"required"`
		Address  string `json:"address"`
		Avatar   string `json:"avatar"`
	}
)

// NewUserService instance for UserService
// @RouterGroup /user/v1
func NewUserService(db *gorm.DB) UserService {
	return &UserServiceImpl{
		userRepo:     models.NewUserQuerySet(db),
		registerRepo: models.NewRegisterUserQuerySet(db),
		passhashRepo: models.NewUserPasshashQuerySet(db),
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
// @Success 200 {object} app.Result{result=service.RegisterToken}
// @Failure 400 {object} app.Result
// @Router /register [post]
func (s *UserServiceImpl) RegisterUser(c *gin.Context) {
	query := RegisterUserQuery{}
	registerModel := models.RegisterUser{}

	if err := c.ShouldBindJSON(&query); err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// cek apakah email sudah ada sebagai user
	count, _ := s.userRepo.EmailEq(query.Email).Count()
	if count > 0 {
		APIResult.Error(c, http.StatusBadRequest, "Email sudah terdaftar")
		return
	}

	// cek apakah email sudah terdaftar
	if err := s.registerRepo.EmailEq(query.Email).One(&registerModel); err != nil {
		token, _, _ := utils.GenerateToken(query.Email)

		// simpan user apabila belum terdaftar
		registerModel.FullName = query.FullName
		registerModel.Email = query.Email
		registerModel.PhoneNum = query.PhoneNum
		registerModel.Token = token
		registerModel.Create(models.DB)

		APIResult.Success(c, &token)
		return
	}

	// return token-nya apabila email sudah terdaftar
	APIResult.Success(c, &RegisterToken{Token: registerModel.Token})
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
func (s *UserServiceImpl) ActivateUser(c *gin.Context) {
	query := ActivateUserQuery{}
	registerModel := models.RegisterUser{}

	if err := c.ShouldBindJSON(&query); err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// cek apakah token tersedia
	if err := s.registerRepo.TokenEq(query.Token).One(&registerModel); err != nil {
		// kembalikan pesan apiError apabila
		// token tidak ada
		APIResult.Error(c, http.StatusBadRequest, "Token invalid")
		return
	}

	// cek apakah email sudah ada sebagai user
	count, _ := s.userRepo.EmailEq(registerModel.Email).Count()
	if count > 0 {
		APIResult.Error(c, http.StatusBadRequest, "Email sudah terdaftar")
		return
	}

	// simpan user
	user := models.User{
		FullName:     registerModel.FullName,
		Email:        registerModel.Email,
		PhoneNum:     registerModel.PhoneNum,
		Active:       true,
		Type:         1,
		RegisteredAt: time.Now().UTC(),
	}
	resUser, err := user.CreateUser()
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, "Tidak dapat mengaktifkan user")
		return
	}

	// buat passhash untuk user
	passhash, _ := utils.GeneratePasshash(query.Passhash)
	userPasshash := models.UserPasshash{
		User:       resUser,
		Passhash:   passhash,
		Deprecated: false,
	}
	// aktifkan user
	userPasshash.ActivateUser()
	// hapus dari register user
	registerModel.Delete(models.DB)

	APIResult.Success(c, &resUser)
}

// MeInfo docs
// @Tags UserService
// @Summary Endpoint untuk informasi user
// @Security bearerAuth
// @Produce json
// @Success 200 {object} app.Result{result=models.User}
// @Failure 401 {object} app.Result
// @Router /me/info [get] [auth]
func (s *UserServiceImpl) MeInfo(c *gin.Context) {
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
func (s *UserServiceImpl) UpdateUserInfo(c *gin.Context) {
	query := UpdateUserQuery{}

	if err := c.ShouldBindJSON(&query); err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	mid.CurrentUser.FullName = query.FullName
	mid.CurrentUser.Email = query.Email
	mid.CurrentUser.PhoneNum = query.PhoneNum
	mid.CurrentUser.Address = query.Address
	mid.CurrentUser.Avatar = query.Avatar
	mid.CurrentUser.Update(
		models.DB,
		models.UserDBSchema.FullName,
		models.UserDBSchema.Email,
		models.UserDBSchema.PhoneNum,
		models.UserDBSchema.Address,
		models.UserDBSchema.Avatar,
	)

	APIResult.Success(c, &mid.CurrentUser)
}
