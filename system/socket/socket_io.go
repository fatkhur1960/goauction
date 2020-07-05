package socket

import (
	"fmt"
	"log"
	"time"

	"github.com/fatkhur1960/goauction/app/middleware"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/repository"
	socketio "github.com/googollee/go-socket.io"
)

type join struct {
	FullName string `json:"full_name"`
	RoomName string `json:"room_name"`
}

type message struct {
	ID         int64     `json:"id"`
	Room       string    `json:"room"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Text       string    `json:"text"`
	Ts         time.Time `json:"ts"`
}

func (msg *message) toReplyMsg() *messageReply {
	userRepo := repository.NewUserRepository()
	return &messageReply{
		ID:       msg.ID,
		Room:     msg.Room,
		Sender:   userRepo.UserSimple(msg.SenderID),
		Receiver: userRepo.UserSimple(msg.ReceiverID),
		Ts:       msg.Ts,
	}
}

type messageReply struct {
	ID       int64              `json:"id"`
	Room     string             `json:"room"`
	Sender   *models.UserSimple `json:"sender"`
	Receiver *models.UserSimple `json:"receiver"`
	Text     string             `json:"text"`
	Ts       time.Time          `json:"ts"`
}

// Handler websocket function
func Handler() *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Printf("WS] Connected: %s\n", middleware.CurrentUser.FullName)
		return nil
	})
	server.OnEvent("/chat", "join", func(s socketio.Conn, join join) {
		log.Printf("WS] %s joined on chat id %s\n", join.FullName, join.RoomName)
		s.Join(join.RoomName)
	})
	server.OnEvent("/chat", "send", func(s socketio.Conn, msg message) {
		log.Printf("WS] got messsage: %s\n", msg.Text)
		server.BroadcastToRoom("/", msg.Room, "reply", msg.toReplyMsg())
	})
	server.OnEvent("/chat", "leave", func(s socketio.Conn, join join) {
		log.Printf("WS] %s leave from %s\n", join.FullName, join.RoomName)
		s.Leave(join.RoomName)
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("WS] meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("WS] closed", reason)
		server.LeaveAllRooms("/chat", s)
	})

	return server
}
