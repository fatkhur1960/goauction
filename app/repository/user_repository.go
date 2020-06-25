package repository

import (
	"errors"
	"log"
	"time"

	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/utils"
)

// UserRepository init user repo
type (
	// UserRepository init implementation
	UserRepository struct {
		UserQs     models.UserQuerySet
		RegisterQs models.RegisterUserQuerySet
		PasshashQs models.UserPasshashQuerySet
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

// RegisterUser dao
func (s *UserRepository) RegisterUser(f string, e string, p string, t string) (models.RegisterUser, error) {
	registerModel := models.RegisterUser{}

	// cek apakah email sudah ada sebagai user
	count, _ := s.UserQs.EmailEq(e).Count()
	if count > 0 {
		return registerModel, errors.New("Email sudah terdaftar")
	}

	// cek apakah email sudah terdaftar
	if err := s.RegisterQs.EmailEq(e).One(&registerModel); err != nil {

		// simpan user apabila belum terdaftar
		registerModel.FullName = f
		registerModel.Email = e
		registerModel.PhoneNum = p
		registerModel.Token = t
		registerModel.RegisteredAt = time.Now().UTC()
		registerModel.Create(models.DB)

		// return user yang baru mendaftar
		return registerModel, nil
	}

	// return user yang sudah terdaftar
	return registerModel, nil
}

// ActivateUser dao
func (s *UserRepository) ActivateUser(token string, passhash string) (models.User, error) {
	registerModel := models.RegisterUser{}

	// cek apakah token tersedia
	if err := s.RegisterQs.TokenEq(token).One(&registerModel); err != nil {
		// kembalikan pesan apiError apabila token tidak ada
		return models.User{}, errors.New("Token invalid")
	}

	// cek apakah email sudah ada sebagai user
	count, _ := s.UserQs.EmailEq(registerModel.Email).Count()
	if count > 0 {
		return models.User{}, errors.New("Email sudah terdaftar")
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
		return models.User{}, errors.New("Tidak dapat mengaktifkan user")
	}

	// buat passhash untuk user
	passhash, hashErr := utils.GeneratePasshash(passhash)
	if hashErr != nil {
		return models.User{}, hashErr
	}

	userPasshash := models.UserPasshash{
		User:       resUser,
		Passhash:   passhash,
		Deprecated: false,
	}
	// aktifkan user
	userPasshash.Create(models.DB)
	// hapus dari register user
	registerModel.Delete(models.DB)

	return resUser, nil
}

// UpdateUser dao
func (s *UserRepository) UpdateUser(userID int64, query UpdateUserQuery) (models.User, error) {
	user := models.User{}
	conn := s.UserQs.IDEq(userID)

	err := conn.GetUpdater().SetFullName(
		query.FullName,
	).SetEmail(query.Email).SetAddress(
		query.Address,
	).SetPhoneNum(query.PhoneNum).SetAvatar(
		query.Avatar,
	).Update()

	if err != nil {
		return models.User{}, errors.New("Tidak dapat mengupdate user")
	}

	conn.One(&user)

	return user, nil
}

// CleanUpUser dao clean all users after testing
// NOTE: using this for testing only
func (s *UserRepository) CleanUpUser() {
	users := []models.User{}
	s.UserQs.All(&users)
	for _, user := range users {
		if err := user.Delete(models.DB); err != nil {
			log.Fatal(err.Error())
		}
	}
}
