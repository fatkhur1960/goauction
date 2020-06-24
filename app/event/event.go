package event

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type (
	// Event base event
	Event interface {
		Send(pid int, e interface{})
	}

	// EventImpl struct for event implementation
	EventImpl struct {}

	// UserRegistered event when user register
	UserRegistered struct {
		FullName string
		Email    string
		Token    string
	}

	// SetBehaviorActor base behavior
	SetBehaviorActor struct{}
)

// Receive events
func (state *SetBehaviorActor) Receive(context actor.Context) {
	switch e := context.Message().(type) {
	case UserRegistered:
		fmt.Printf("Hello %v\n", e.FullName)
	}
}

// NewSetBehaviorActor create instance
func NewSetBehaviorActor() actor.Actor {
    return &SetBehaviorActor{}
}

// NewEvent create instance
func NewEvent() Event {
	return &EventImpl{}
}
