package event

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/gin-gonic/gin"
)

// BaseEvent base for implement Handle
type BaseEvent interface {
	Handle(p interface{})
}

// Listener initialize eventListener
var Listener eventListener

type eventListener struct {
	handlers []BaseEvent
}

// Register adds an event handler for this event
func (l *eventListener) Register(handler BaseEvent) {
	l.handlers = append(l.handlers, handler)
}

// Emmit sends out an event with the payload
func (l eventListener) Emmit(payload BaseEvent) {
	if gin.IsDebugging() {
		data, _ := json.Marshal(&payload)
		log.Println("event.Listener] Got Event:", reflect.TypeOf(payload).String(), string(data))
	}

	for _, handler := range l.handlers {
		go handler.Handle(payload)
	}
}
