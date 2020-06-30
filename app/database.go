package app

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	// import postgres drive
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB connection
var DB *gorm.DB

// ConnectDatabase method to connect with db
func ConnectDatabase() {
	dbConf := fmt.Sprintf(
		"host=%s sslmode=%s port=%s user=%s dbname=%s password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("SSL_MODE"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)
	database, err := gorm.Open("postgres", dbConf)

	if err != nil {
		panic("DB Error: " + err.Error())
	}

	DB = database
}

// ConnectDatabaseTest method for testing
func ConnectDatabaseTest() {
	dbConf := fmt.Sprintf(
		"host=%s sslmode=%s port=%s user=%s dbname=%s password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("SSL_MODE"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME_TEST"),
		os.Getenv("DB_PASSWORD"),
	)
	database, err := gorm.Open("postgres", dbConf)

	if err != nil {
		panic("DB Error: " + err.Error())
	}

	DB = database
}
