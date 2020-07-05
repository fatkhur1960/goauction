package models

import (
	"time"

	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/types"
)

//go:generate goqueryset -in chat.go

// Chat model
// gen:qs
type Chat struct {
	ID           int64      `json:"id"`
	InitiatorID  int64      `json:"initiator_id"`
	SubscriberID int64      `json:"subscriber_id"`
	LastUpdated  *time.Time `json:"last_updated"`
	TS           *time.Time `json:"ts"`
}

// Message model
// gen:qs
type Message struct {
	ID             int64      `json:"id"`
	ChatID         int64      `json:"chat_id"`
	SenderID       int64      `json:"sender_id"`
	ReceiverID     int64      `json:"receiver_id"`
	Text           string     `json:"text"`
	Deleted        bool       `json:"deleted"`
	AttachmentKind int        `json:"attachment_kind"`
	AttachmentData string     `json:"attachment_data"`
	TS             *time.Time `json:"ts"`
}

// ChatHistory model
// gen:qs
type ChatHistory struct {
	ID        int64 `json:"id"`
	ChatID    int64 `json:"chat_id"`
	OwnerID   int64 `json:"owner_id"`
	MessageID int64 `json:"message_id"`
}

// TableName override for model ChatHistory
func (ChatHistory) TableName() string {
	return "user_chat_histories"
}

// ToAPI implementation for chat
func (c *Chat) ToAPI(userID int64) types.Chat {
	display := UserSimple{}
	dao := app.DB.Select("id, full_name, avatar")
	if c.InitiatorID == userID {
		dao.Where("id = ?", c.SubscriberID).First(&display)
	} else {
		dao.Where("id = ?", c.InitiatorID).First(&display)
	}

	return types.Chat{
		ID:           c.ID,
		InitiatorID:  c.InitiatorID,
		SubscriberID: c.SubscriberID,
		LastUpdated:  c.LastUpdated,
		TS:           c.TS,
		Display:      display,
	}
}
