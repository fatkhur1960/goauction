package event

import (
	"fmt"

	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/repository"
	"github.com/fatkhur1960/goauction/system/core"
)

// StartupEvent --
type StartupEvent struct{}

// Handle startup event
func (e StartupEvent) Handle() error {
	// TODO: Create something here where server started
	return nil
}

// UserRegisteredEvent is the data for when a user is created
type UserRegisteredEvent struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	PhoneNum string `json:"phone_num"`
	Token    string `json:"token"`
}

// Handle event user registered
func (e UserRegisteredEvent) Handle() error {
	fmt.Println(e.Email, " Registered")
	return nil
}

// UserBidProductEvent is the data when user bid a product
type UserBidProductEvent struct {
	User    models.User
	Product models.Product
	BidData models.ProductBidder
}

// Handle event bid product
func (e UserBidProductEvent) Handle() error {
	notifRepo := repository.NewNotifRepository()

	title := fmt.Sprintf("Hai, %s ngebid produk Anda", e.User.FullName)
	content := fmt.Sprintf("%s ngebid produk `%s` dengan harga %v", e.User.FullName, e.Product.ProductName, e.BidData.BidPrice)

	_, err := notifRepo.CreateNotif(e.Product.UserID, title, content, core.GotBidder, e.Product.ID)
	if err != nil {
		return err
	}

	return nil
}
