package models

import (
	"sync"
	"time"

	"github.com/fatkhur1960/goauction/app"
)

//go:generate goqueryset -in user.go

// User definisi model untuk user
// gen:qs
type User struct {
	ID           int64      `json:"id"`
	FullName     string     `json:"full_name"`
	Email        string     `json:"email"`
	PhoneNum     string     `json:"phone_num,omitempty"`
	Address      string     `json:"address"`
	Avatar       string     `json:"avatar"`
	Type         int        `json:"type,omitempty"`
	Active       bool       `json:"active,omitempty"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	RegisteredAt time.Time  `json:"registered_at,omitempty"`
}

// UserSimple ...
type UserSimple struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
	Address  string `json:"address"`
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
	User       *User     `json:"-"`
	Passhash   string    `json:"-"`
	Deprecated bool      `json:"deprecated"`
	Ver        int       `json:"ver"`
	Created    time.Time `json:"created"`
}

// AccessToken definisi model untuk access token user
// gen:qs
type AccessToken struct {
	sync.Mutex
	UserID    int64     `json:"user_id"`
	User      *User     `json:"-"`
	Token     string    `json:"token"`
	Created   time.Time `json:"created"`
	ValidThru time.Time `json:"valid_thru"`
}

// Store definisi model untuk store
// gen:qs
type Store struct {
	ID           int64       `json:"id"`
	Name         string      `json:"name"`
	Info         string      `json:"info"`
	Owner        *UserSimple `json:"owner"`
	OwnerID      int64       `json:"owner_id"`
	Announcement string      `json:"announcement"`
	ProductCount int         `json:"product_count"`
	Province     string      `json:"province"`
	Regency      string      `json:"regency"`
	SUBDistrict  string      `json:"sub_district"`
	Village      string      `json:"village"`
	Address      string      `json:"address"`
	LastUpdated  *time.Time  `json:"last_updated"`
	TS           *time.Time  `json:"ts"`
}

// UserConnect model
// gen:qs
type UserConnect struct {
	UserID       int64  `json:"user_id"`
	ProviderName string `json:"provider_name"`
	AppID        string `json:"app_id"`
}

// TableName for UserSimple model
func (UserSimple) TableName() string {
	return "users"
}

// CreateUser dao untuk menambahkan user
func (user *User) CreateUser() (*User, error) {
	err := app.DB.FirstOrCreate(&user, User{Email: user.Email}).Error
	return user, err
}

// CreateAccessToken dao untuk menambahkan user
func (token *AccessToken) CreateAccessToken() (*AccessToken, error) {
	err := app.DB.FirstOrCreate(&token, AccessToken{UserID: token.User.ID}).Error
	return token, err
}

// RemoveAccessToken dao untuk menghapus token
func (token *AccessToken) RemoveAccessToken() error {
	return app.DB.Delete(token).Error
}

// IsExpired dao untuk check apakah token sudah expired
func (token *AccessToken) IsExpired() bool {
	diff := token.ValidThru.Sub(time.Now().UTC())
	return diff.Hours() < 0
}

// ActivateUser dao untuk mengaktifkan user
func (userPasshash *UserPasshash) ActivateUser() error {
	return app.DB.Create(&userPasshash).Error
}

// ClearRegisteredUser dao untuk menghapus user yang sudah diaktifkan
func (user *RegisterUser) ClearRegisteredUser() error {
	return app.DB.Delete(&user).Error
}
