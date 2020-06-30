package event

import (
	"github.com/fatkhur1960/goauction/system/queue"
)

// Listener ...
type Listener struct {
	Queue chan queue.Queuable
}

// NewListener instance
func NewListener(q chan queue.Queuable) *Listener {
	event := &Listener{
		Queue: q,
	}

	return event
}

// Emmit sends out an event with the payload
func (l Listener) Emmit(payload interface{ Handle() error }) {
	l.Queue <- payload
}
