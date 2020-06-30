package repository

import (
	"time"

	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/utils"
)

// ChatRepository init
type ChatRepository struct {
	chatQs models.ChatQuerySet
	msgQs  models.MessageQuerySet
	chQs   models.ChatHistoryQuerySet
}

// ChatMessageQuery --
type ChatMessageQuery struct {
	ChatID         int64  `json:"chat_id" binding:"required"`
	SenderID       int64  `json:"sender_id" binding:"required"`
	ReceiverID     int64  `json:"receiver_id" binding:"required"`
	Text           string `json:"text" binding:"required"`
	AttachmentKind int    `json:"attachment_kind"`
	AttachmentData string `json:"attachment_data"`
}

// NewChatRepository instance
func NewChatRepository() *ChatRepository {
	return &ChatRepository{
		chatQs: models.NewChatQuerySet(app.DB),
		msgQs:  models.NewMessageQuerySet(app.DB),
		chQs:   models.NewChatHistoryQuerySet(app.DB),
	}
}

// CreateChat create new chat room
func (r *ChatRepository) CreateChat(initiatorID int64, subscriberID int64) (models.Chat, error) {
	chat := models.Chat{
		InitiatorID:  initiatorID,
		SubscriberID: subscriberID,
		LastUpdated:  &utils.NOW,
		TS:           &utils.NOW,
	}

	if err := chat.Create(app.DB); err != nil {
		return models.Chat{}, err
	}

	return chat, nil
}

// CreateChatMessage create new chat message
func (r *ChatRepository) CreateChatMessage(query ChatMessageQuery) (models.Message, error) {
	message := models.Message{
		ChatID:         query.ChatID,
		SenderID:       query.SenderID,
		ReceiverID:     query.ReceiverID,
		Text:           query.Text,
		AttachmentKind: query.AttachmentKind,
		AttachmentData: query.AttachmentData,
	}
	message.Create(app.DB)

	// Create chat history
	{
		u1History := models.ChatHistory{
			ChatID:    query.ChatID,
			OwnerID:   query.SenderID,
			MessageID: message.ID,
		}
		u1History.Create(app.DB)

		u2History := models.ChatHistory{
			ChatID:    query.ChatID,
			OwnerID:   query.ReceiverID,
			MessageID: message.ID,
		}
		u2History.Create(app.DB)
	}

	// Update last updated time on chat room
	{
		updated := time.Now().UTC()
		r.chatQs.IDEq(query.ChatID).GetUpdater().SetLastUpdated(&updated).Update()
	}

	return message, nil
}

// GetUserChats listing user chats
func (r *ChatRepository) GetUserChats(userID int64, offset int, limit int) ([]models.Chat, int, error) {
	chats := []models.Chat{}
	count := 0
	dao := r.chatQs.GetDB().Where("initiator_id = ?", userID).Or("subscriber_id = ?", userID)
	dao.Order("id DESC").Offset(offset).Limit(limit).Find(&chats)
	dao.Count(&count)

	return chats, count, nil
}

// GetChatMessages list chat message
func (r *ChatRepository) GetChatMessages(chatID int64, userID int64, offset int, limit int) ([]models.Message, int, error) {
	messages := []models.Message{}
	count := 0
	dao := r.chQs.GetDB().Select("messages.*").Joins("JOIN messages ON messages.id = user_chat_histories.message_id").Where("user_chat_histories.chat_id = ?", chatID).Where("user_chat_histories.owner_id = ?", userID)
	dao.Count(&count)
	res := dao.Offset(offset).Limit(limit).Find(&messages)
	if res.Error != nil {
		return messages, count, res.Error
	}

	return messages, count, nil
}