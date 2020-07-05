package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/jinzhu/gorm"
)

// UserRepository init user repo
type (
	// UserRepository init implementation
	UserRepository struct {
		sync.Mutex
		userQs     models.UserQuerySet
		connQs     models.UserConnectQuerySet
		registerQs models.RegisterUserQuerySet
		passhashQs models.UserPasshashQuerySet
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

// NewUserRepository intance
func NewUserRepository() *UserRepository {
	return &UserRepository{
		userQs:     models.NewUserQuerySet(app.DB),
		connQs:     models.NewUserConnectQuerySet(app.DB),
		registerQs: models.NewRegisterUserQuerySet(app.DB),
		passhashQs: models.NewUserPasshashQuerySet(app.DB),
	}
}

// GetByID get user by id
func (s *UserRepository) GetByID(userID int64) (models.User, error) {
	s.Lock()
	defer s.Unlock()
	user := models.User{}
	if err := s.userQs.IDEq(userID).One(&user); err != nil {
		return user, err
	}

	return user, nil
}

// UserSimple digunakan untuk medapatkan simple user
func (s *UserRepository) UserSimple(userID int64) *models.UserSimple {
	simple := models.User{}
	s.userQs.IDEq(userID).Select(
		models.UserDBSchema.ID,
		models.UserDBSchema.FullName,
		models.UserDBSchema.Avatar,
		models.UserDBSchema.Address,
	).One(&simple)

	return &models.UserSimple{
		ID:       simple.ID,
		FullName: simple.FullName,
		Avatar:   simple.Avatar,
		Address:  simple.Address,
	}
}

// RegisterUser dao
func (s *UserRepository) RegisterUser(f string, e string, p string, t string) (models.RegisterUser, error) {
	registerModel := models.RegisterUser{}

	// cek apakah email sudah ada sebagai user
	count, _ := s.userQs.EmailEq(e).Count()
	if count > 0 {
		return registerModel, errors.New("Email sudah terdaftar")
	}

	// cek apakah email sudah terdaftar
	if err := s.registerQs.EmailEq(e).One(&registerModel); err != nil {

		// simpan user apabila belum terdaftar
		registerModel.FullName = f
		registerModel.Email = e
		registerModel.PhoneNum = p
		registerModel.Token = t
		registerModel.RegisteredAt = time.Now().UTC()
		registerModel.Create(app.DB)

		// return user yang baru mendaftar
		return registerModel, nil
	}

	// return user yang sudah terdaftar
	return registerModel, nil
}

// ActivateUser dao
func (s *UserRepository) ActivateUser(token string, passhash string) (*models.User, error) {
	registerModel := models.RegisterUser{}

	// cek apakah token tersedia
	if err := s.registerQs.TokenEq(token).One(&registerModel); err != nil {
		// kembalikan pesan apiError apabila token tidak ada
		return &models.User{}, errors.New("Token invalid")
	}

	// cek apakah email sudah ada sebagai user
	count, _ := s.userQs.EmailEq(registerModel.Email).Count()
	if count > 0 {
		return &models.User{}, errors.New("Email sudah terdaftar")
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
		return &models.User{}, errors.New("Tidak dapat mengaktifkan user")
	}

	// buat passhash untuk user
	passhash, hashErr := utils.GeneratePasshash(passhash)
	if hashErr != nil {
		return &models.User{}, hashErr
	}

	userPasshash := models.UserPasshash{
		User:       resUser,
		Passhash:   passhash,
		Deprecated: false,
	}
	// aktifkan user
	userPasshash.Create(app.DB)
	// hapus dari register user
	registerModel.Delete(app.DB)

	return resUser, nil
}

// UpdateUser dao
func (s *UserRepository) UpdateUser(userID int64, query UpdateUserQuery) (*models.User, error) {
	user := models.User{}
	conn := s.userQs.IDEq(userID)

	err := conn.GetUpdater().SetFullName(
		query.FullName,
	).SetEmail(query.Email).SetAddress(
		query.Address,
	).SetPhoneNum(query.PhoneNum).SetAvatar(
		query.Avatar,
	).Update()

	if err != nil {
		return &models.User{}, errors.New("Tidak dapat mengupdate user")
	}

	conn.One(&user)

	return &user, nil
}

// CreateUserConnect app id untuk spesifik user, digunakan untuk event push notif.
func (s *UserRepository) CreateUserConnect(userID int64, appID string, providerName string) error {
	conn := models.UserConnect{
		UserID:       userID,
		AppID:        appID,
		ProviderName: providerName,
	}

	err := s.connQs.UserIDEq(userID).GetUpdater().SetAppID(appID).SetProviderName(providerName).Update()

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return conn.Create(app.DB)
		}
		return err
	}

	return nil
}

// GetAppID berdasarkan user id
func (s *UserRepository) GetAppID(userID int64) (models.UserConnect, error) {
	conn := models.UserConnect{}
	if err := s.connQs.UserIDEq(userID).One(&conn); err != nil {
		return conn, err
	}

	return conn, nil
}

// RemoveUserConnect app id berdasarkan user id
func (s *UserRepository) RemoveUserConnect(userID int64) error {
	return s.connQs.UserIDEq(userID).Delete()
}
