package models

import "time"

//go:generate goqueryset -in user_notif.go

// UserNotif model
// gen:qs
type UserNotif struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	NotifType int        `json:"notif_type"`
	Target    int        `json:"target"`
	CreatedAT *time.Time `json:"created_at"`
	Read      bool       `json:"read"`
}
