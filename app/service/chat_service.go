package service

import (
	"net/http"

	mid "github.com/fatkhur1960/goauction/app/middleware"
	"github.com/fatkhur1960/goauction/app/repository"
	"github.com/gin-gonic/gin"
)

type (
	// ChatService api implementation
	ChatService struct {
		chatRepo *repository.ChatRepository
	}

	// CreateChatQuery --
	CreateChatQuery struct {
		UserID int64 `json:"user_id" binding:"required"`
	}
)

// NewChatService instance
// @RouterGroup /chat/v1
func NewChatService() *ChatService {
	return &ChatService{chatRepo: repository.NewChatRepository()}
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
	// query := CreateChatQuery{}
	if err := validateRequest(c, query); err != nil {
		return
	}

	chat, err := s.chatRepo.CreateChat(mid.CurrentUser.ID, query.UserID)
	if err != nil {
		APIResult.Error(c, http.StatusBadRequest, "Tidak dapat membuat chat")
	}

	APIResult.Success(c, chat.ToAPI(mid.CurrentUser.ID))
}
