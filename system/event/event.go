package event

import (
	"fmt"
	"log"

	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/repository"
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/fatkhur1960/goauction/system/core"
	"github.com/fatkhur1960/goauction/system/notificator"
)

var notif = notificator.NewNotifHandler()

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

// Handle event for UserRegisteredEvent
func (e UserRegisteredEvent) Handle() error {
	log.Println("Event]", e.Email, "Registered")
	return nil
}

// UserBidProductEvent is the data when user bid a product
type UserBidProductEvent struct {
	User    *models.User
	Product models.Product
	BidData models.ProductBidder
}

// Handle event for UserBidProductEvent
func (e *UserBidProductEvent) Handle() error {
	notifRepo := repository.NewNotifRepository()
	storeRepo := repository.NewStoreRepository()
	store, _ := storeRepo.GetByID(e.Product.StoreID)

	title := fmt.Sprintf("Hai, %s ngebid produk Anda", e.User.FullName)
	content := fmt.Sprintf("%s ngebid produk `%s` dengan harga %v", e.User.FullName, e.Product.ProductName, e.BidData.BidPrice)

	userNotif, err := notifRepo.CreateNotif(store.OwnerID, title, content, core.GotBidder, e.Product.ID)
	if err != nil {
		return err
	}

	payload := &notificator.Payload{
		NotifID:    userNotif.ID,
		ReceiverID: userNotif.UserID,
		TargetID:   e.Product.ID,
		NotifKind:  core.GotBidder,
		Item:       &e.Product,
		Title:      title,
		Message:    content,
		Created:    &utils.NOW,
	}

	notif.Send(payload)

	return nil
}
