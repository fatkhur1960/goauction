package models

import (
	"time"
)

//go:generate goqueryset -in user.go

// User definisi model untuk user
// gen:qs
type User struct {
	ID           int64     `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PhoneNum     string    `json:"phone_num"`
	Address      string    `json:"address"`
	Avatar       string    `json:"avatar"`
	Type         int       `json:"type"`
	Active       bool      `json:"active"`
	LastLogin    time.Time `json:"last_login"`
	RegisteredAt time.Time `json:"registered_at"`
}

// RegisterUser definisi model untuk register user
// gen:qs
type RegisterUser struct {
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PhoneNum     string    `json:"phone_num"`
	Token        string    `json:"token"`
	Code         string    `json:"code"`
	RegisteredAt time.Time `json:"registered_at"`
}

// UserPasshash definisi model untuk mengaktifkan user
// gen:qs
type UserPasshash struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	User       User      `json:"-"`
	Passhash   string    `json:"-"`
	Deprecated bool      `json:"deprecated"`
	Ver        int       `json:"ver"`
	Created    time.Time `json:"created"`
}

// AccessToken definisi model untuk access token user
// gen:qs
type AccessToken struct {
	UserID    int64     `json:"user_id"`
	User      User      `json:"-"`
	Token     string    `json:"token"`
	Created   time.Time `json:"created"`
	ValidThru time.Time `json:"valid_thru"`
}

// CreateUser dao untuk menambahkan user
func (user User) CreateUser() (User, error) {
	err := DB.FirstOrCreate(&user, User{Email: user.Email}).Error
	return user, err
}

// CreateAccessToken dao untuk menambahkan user
func (token AccessToken) CreateAccessToken() (AccessToken, error) {
	err := DB.FirstOrCreate(&token, AccessToken{UserID: token.User.ID}).Error
	return token, err
}

// RemoveAccessToken dao untuk menghapus token
func (token AccessToken) RemoveAccessToken() error {
	return DB.Delete(token).Error
}

// IsExpired dao untuk check apakah token sudah expired
func (token AccessToken) IsExpired() bool {
	diff := token.ValidThru.Sub(time.Now().UTC())
	return diff.Hours() < 0
}

// ActivateUser dao untuk mengaktifkan user
func (userPasshash UserPasshash) ActivateUser() error {
	return DB.Create(&userPasshash).Error
}

// ClearRegisteredUser dao untuk menghapus user yang sudah diaktifkan
func (user RegisterUser) ClearRegisteredUser() error {
	return DB.Delete(&user).Error
}
