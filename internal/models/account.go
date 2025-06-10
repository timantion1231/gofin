package models

import "time"

type Account struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	AccountNumber string    `json:"account_number"`
	Balance       float64   `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}
