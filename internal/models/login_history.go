package models

import "time"

type LoginHistory struct {
	ID        uint      `json:"-"`
	UserID    uint      `json:"-"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}
