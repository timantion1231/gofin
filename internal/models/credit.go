package models

import "time"

type Credit struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	AccountID      int       `json:"account_id"`
	Amount         float64   `json:"amount"`
	Rate           float64   `json:"rate"`
	Months         int       `json:"months"`
	MonthlyPayment float64   `json:"monthly_payment"`
	CreatedAt      time.Time `json:"created_at"`
	Status         string    `json:"status"` // active, closed, overdue
}

type CreditPayment struct {
	ID       int        `json:"id"`
	CreditID int        `json:"credit_id"`
	DueDate  time.Time  `json:"due_date"`
	Amount   float64    `json:"amount"`
	Paid     bool       `json:"paid"`
	PaidAt   *time.Time `json:"paid_at,omitempty"`
	Fine     float64    `json:"fine"`
}
