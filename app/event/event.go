package event

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// RegisterEvents list of available event register
func RegisterEvents(c *gin.Context) {
	Listener.Register(&StartupEvent{})
	Listener.Register(&UserRegisteredPayload{})
	c.Next()
}

// StartupEvent --
type StartupEvent struct{}

// Handle startup event
func (e *StartupEvent) Handle(p interface{}) {
	// TODO: Create something here where server started
}

// UserRegisteredPayload is the data for when a user is created
type UserRegisteredPayload struct {
	FullName string
	Email    string
	PhoneNum string
	Token    string
}

// Handle event user registered
func (e *UserRegisteredPayload) Handle(p interface{}) {
	fmt.Println(p.(*UserRegisteredPayload).Email, " Registered")
}
