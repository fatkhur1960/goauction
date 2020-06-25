package repository

import (
	"errors"
	"time"

	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/utils"
)

// AuthRepository init repo
type AuthRepository struct {
	UserQs     models.UserQuerySet
	TokenQs    models.AccessTokenQuerySet
	PasshashQs models.UserPasshashQuerySet
}

// AuthorizeUser method untuk mengotorisasi user
func (s *AuthRepository) AuthorizeUser(email string, passhash string) (interface{}, error) {
	var user models.User
	var userPasshash models.UserPasshash

	// check apakah user ada di db
	err := s.UserQs.EmailEq(email).ActiveEq(true).One(&user)
	if err != nil {
		return nil, errors.New("Email tidak ditemukan")
	}

	// check passhash
	s.PasshashQs.UserIDEq(user.ID).One(&userPasshash)
	if !utils.CheckPasshash(passhash, userPasshash.Passhash) {
		return nil, errors.New("Password tidak cocok")
	}

	// generate access token
	token, expireTime, _ := utils.GenerateToken(user.Email)
	accessToken := models.AccessToken{
		User:      user,
		Token:     token,
		Created:   time.Now().UTC(),
		ValidThru: expireTime,
	}

	return accessToken.CreateAccessToken()
}
