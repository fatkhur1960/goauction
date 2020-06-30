package service

import (
	"log"
	"net/http"

	mid "github.com/fatkhur1960/goauction/app/middleware"
	repo "github.com/fatkhur1960/goauction/app/repository"
	"github.com/fatkhur1960/goauction/app/types"
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/fatkhur1960/goauction/system/event"
	"github.com/fatkhur1960/goauction/system/queue"
	"github.com/gin-gonic/gin"
)

type (
	// UserService implementation for users
	UserService struct {
		userRepo      *repo.UserRepository
		productRepo   *repo.ProductRepository
		notifRepo     *repo.NotifRepository
		eventListener *event.Listener
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

	// ReadNotifQuery definisi query untuk menandai notif sudah dibaca
	ReadNotifQuery struct {
		NotifIds []int64 `json:"notif_ids" binding:"required"`
	}
)

// NewUserService instance for UserService
// @RouterGroup /user/v1
func NewUserService() *UserService {
	return &UserService{
		userRepo:      repo.NewUserRepository(),
		notifRepo:     repo.NewNotifRepository(),
		productRepo:   repo.NewProductRepository(),
		eventListener: event.NewListener(queue.JobQueue),
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
func (s *UserService) RegisterUser(c *gin.Context, query *RegisterUserQuery) {
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
	s.eventListener.Emmit(event.UserRegisteredEvent{
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
func (s *UserService) ActivateUser(c *gin.Context, query *ActivateUserQuery) {
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
	if validateRequest(c, &query) != nil {
		return
	}

	user, err := s.userRepo.UpdateUser(mid.CurrentUser.ID, query)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	APIResult.Success(c, user)
}

// ListUserBids docs
// @Tags UserService
// @Security bearerAuth
// @Summary Endpoint untuk mendapatkan bid history
// @Product json
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Param query query string false "Query"
// @Param filter query string false "Filter"
// @Success 200 {object} app.Result{result=EntriesResult{entries=[]service.Product}}
// @Failure 400 {object} app.Result
// @Router /bids [get] [auth]
func (s *UserService) ListUserBids(c *gin.Context) {
	query := QueryEntries{}
	if err := query.validate(c, types.ValidateQuery); err != nil {
		return
	}

	rawEntries, count, err := s.productRepo.GetBidProductList(mid.CurrentUser.ID, query.Offset, query.Limit)
	if err != nil {
		log.Fatal("UserService]", err)
	}

	entries := []types.Product{}
	for _, product := range rawEntries {
		entries = append(entries, product.ToAPI(mid.CurrentUser.ID))
	}

	APIResult.Success(c, EntriesResult{entries, count})
}

// ListUserNotifs docs
// @Tags UserService
// @Security bearerAuth
// @Summary Endpoint untuk mendapatkan list notif untuk current user
// @Produce json
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Param query query string false "Query"
// @Param filter query string false "Filter"
// @Success 200 {object} app.Result{result=EntriesResult{entries=[]models.UserNotif}}
// @Failure 400 {object} app.Result
// @Router /notifs [get] [auth]
func (s *UserService) ListUserNotifs(c *gin.Context) {
	query := QueryEntries{}
	if query.validate(c, types.ValidateQuery) != nil {
		return
	}

	entries, count, err := s.notifRepo.GetUserNotif(mid.CurrentUser.ID, query.Offset, query.Limit)
	if err != nil {
		log.Fatal(err)
	}

	APIResult.Success(c, EntriesResult{entries, count})
}

// MarkAsReadNotif docs
// @Tags UserService
// @Security bearerAuth
// @Summary endpoint untuk menandai notif sudah terbaca
// @Produce json
// @Param notif_ids body []int true "NotifIds"
// @Success 200 {object} app.Result
// @Failure 400 {object} app.Result
// @Router /notifs/read [post] [auth]
func (s *UserService) MarkAsReadNotif(c *gin.Context) {
	query := ReadNotifQuery{}
	if validateRequest(c, &query) != nil {
		return
	}

	err := s.notifRepo.MarkAsRead(query.NotifIds, mid.CurrentUser.ID)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	APIResult.Success(c, nil)
}
