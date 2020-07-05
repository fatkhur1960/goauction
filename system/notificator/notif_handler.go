package notificator

import (
	"log"
	"os"
	"time"

	"github.com/appleboy/go-fcm"
	"github.com/fatkhur1960/goauction/app/repository"
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/fatkhur1960/goauction/system/core"
)

// Payload for notificator
type Payload struct {
	NotifID     int64          `json:"notif_id"`
	ReceiverID  int64          `json:"receiver_id"`
	TargetID    int64          `json:"target_id"`
	NotifKind   core.NotifType `json:"notif_kind"`
	Item        interface{}    `json:"item"`
	Title       string         `json:"title"`
	Message     string         `json:"message"`
	Created     *time.Time     `json:"created"`
	ClickAction string         `json:"click_action"`
}

// NotifHandler holder
type NotifHandler struct {
	ServerKey string
}

// NewNotifHandler instance
func NewNotifHandler() *NotifHandler {
	return &NotifHandler{
		ServerKey: os.Getenv("FCM_SERVER_KEY"),
	}
}

// Send notif with payload
func (h *NotifHandler) Send(payload *Payload) {
	repo := repository.NewUserRepository()
	userConn, err := repo.GetAppID(payload.ReceiverID)
	if err != nil {
		log.Printf("NotifHandler] AppID Error: %s", err.Error())
		return
	}

	data := utils.StructToMap(payload)
	msg := &fcm.Message{
		To:   userConn.AppID,
		Data: data,
	}

	// Create a FCM client to send the message.
	client, clientErr := fcm.NewClient(h.ServerKey)
	if err != nil {
		log.Printf("NotifHandler] Client Error: %s", clientErr.Error())
		return
	}

	// Send the message and receive the response without retries.
	response, err := client.Send(msg)
	if err != nil {
		log.Printf("NotifHandler] Send Error: %s", err.Error())
		return
	}

	log.Printf("NotifHandler] Response %#v\n", response)
}
