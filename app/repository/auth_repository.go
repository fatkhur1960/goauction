package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/utils"
)

// AuthRepository init repo
type AuthRepository struct {
	sync.Mutex
	userQs     models.UserQuerySet
	tokenQs    models.AccessTokenQuerySet
	passhashQs models.UserPasshashQuerySet
}

// NewAuthRepository new instance
func NewAuthRepository() *AuthRepository {
	return &AuthRepository{
		userQs:     models.NewUserQuerySet(app.DB),
		tokenQs:    models.NewAccessTokenQuerySet(app.DB),
		passhashQs: models.NewUserPasshashQuerySet(app.DB),
	}
}

// GetAccessToken dao
func (s *AuthRepository) GetAccessToken(token string) (*models.AccessToken, error) {
	s.Lock()
	defer s.Unlock()
	at := models.AccessToken{}
	if err := s.tokenQs.TokenEq(token).One(&at); err != nil {
		return &at, err
	}

	return &at, nil
}

// AuthorizeUser method untuk mengotorisasi user
func (s *AuthRepository) AuthorizeUser(email string, passhash string) (interface{}, error) {
	var user models.User
	var userPasshash models.UserPasshash

	// check apakah user ada di db
	err := s.userQs.EmailEq(email).ActiveEq(true).One(&user)
	if err != nil {
		return nil, errors.New("Email tidak ditemukan")
	}

	// check passhash
	s.passhashQs.UserIDEq(user.ID).One(&userPasshash)
	if !utils.CheckPasshash(passhash, userPasshash.Passhash) {
		return nil, errors.New("Password tidak cocok")
	}

	// generate access token
	token, expireTime, _ := utils.GenerateToken(user.Email)
	accessToken := models.AccessToken{
		User:      &user,
		Token:     token,
		Created:   time.Now().UTC(),
		ValidThru: expireTime,
	}
	s.userQs.IDEq(user.ID).GetUpdater().SetLastLogin(&utils.NOW).Update()

	return accessToken.CreateAccessToken()
}

// UnauthorizeUser method untuk menghapus otorisasi user
func (s *AuthRepository) UnauthorizeUser(userID int64) error {
	return s.tokenQs.UserIDEq(userID).Delete()
}
