package utils

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// GeneratePasshash create hashed string
func GeneratePasshash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasshash check hashed string
func CheckPasshash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken method untuk generate jwt token
func GenerateToken(email string) (string, time.Time, error) {
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	expireTime := time.Now().Add(time.Hour * 24 * 7).UTC()
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_email"] = email
	atClaims["exp"] = expireTime.Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", expireTime, err
	}
	return token, expireTime, nil
}
