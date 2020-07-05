package service

import (
	"net/http"
	"sync"

	mid "github.com/fatkhur1960/goauction/app/middleware"
	"github.com/fatkhur1960/goauction/app/models"
	repo "github.com/fatkhur1960/goauction/app/repository"
	"github.com/fatkhur1960/goauction/app/types"
	"github.com/gin-gonic/gin"
)

type (
	// ChatService api implementation
	ChatService struct {
		sync.Mutex
		chatRepo *repo.ChatRepository
	}

	// CreateChatQuery --
	CreateChatQuery struct {
		UserID int64 `json:"user_id" binding:"required"`
	}
)

// NewChatService instance
// @RouterGroup /chat/v1
func NewChatService() *ChatService {
	return &ChatService{chatRepo: repo.NewChatRepository()}
}

// CreateChatRoom docs
// @Summary Endpoint untuk membuat chat room
// @Tags ChatService
// @Accept json
// @Produce json
// @Param user_id body int true "UserID"
// @Success 200 {object} app.Result{result=models.Chat}
// @Failure 400 {object} app.Result
// @Router /new-room [post] [auth]
func (s *ChatService) CreateChatRoom(c *gin.Context, query *CreateChatQuery) {
	chat, err := s.chatRepo.CreateChat(mid.CurrentUser.ID, query.UserID)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, "Tidak dapat membuat chat")
		return
	}

	APIResult.Success(c, chat.ToAPI(mid.CurrentUser.ID))
}

// ListChatRooms docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk menampilkan list chat room
// @Accept json
// @Produce json
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Param query query string false "Query"
// @Param filter query string false "Filter"
// @Success 200 {object} app.Result{result=EntriesResult{entries=[]types.Chat}}
// @Failure 400 {object} app.Result
// @Router /list [get] [auth]
func (s *ChatService) ListChatRooms(c *gin.Context, query *QueryEntries) {
	entries := []types.Chat{}
	chats, count, _ := s.chatRepo.GetUserChatRooms(mid.CurrentUser.ID, query.Offset, query.Limit)

	for _, chat := range chats {
		entries = append(entries, chat.ToAPI(mid.CurrentUser.ID))
	}

	APIResult.Success(c, EntriesResult{entries, count})
}

// SendMessage docs
// @Tags ChatService
// @Security bearerAuth
// @Summary Endpoint untuk menambahkan product
// @Accept json
// @Param chat_id body int true "ChatID"
// @Param receiver_id body int true "ReceiverID"
// @Param text body string true "Text"
// @Param attachment_kind body int true "AttachmentKind"
// @Param attachment_data body string true "AttachmentData"
// @Produce json
// @Success 200 {object} app.Result{result=models.Message}
// @Failure 400 {object} app.Result
// @Router /send-message [post] [auth]
func (s *ChatService) SendMessage(c *gin.Context, query *repo.ChatMessageQuery) {
	message, err := s.chatRepo.CreateChatMessage(mid.CurrentUser.ID, *query)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	APIResult.Success(c, message)
}

// ListChatMessages docs
// @Tags ProductService
// @Security bearerAuth
// @Summary Endpoint untuk menampilkan list chat room
// @Accept json
// @Produce json
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Param query query string false "Query"
// @Param filter query string false "Filter"
// @Success 200 {object} app.Result{result=EntriesResult{entries=[]models.Message}}
// @Failure 400 {object} app.Result
// @Router /list-messages [get] [auth]
func (s *ChatService) ListChatMessages(c *gin.Context, query *QueryMessages) {
	entries := []models.Message{}
	messages, count, _ := s.chatRepo.GetChatMessages(query.ChatID, mid.CurrentUser.ID, query.Offset, query.Limit)

	for _, message := range messages {
		entries = append(entries, message)
	}

	APIResult.Success(c, EntriesResult{entries, count})
}
