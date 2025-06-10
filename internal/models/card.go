package models

import "time"

type Card struct {
	ID                  int       `json:"id"`
	AccountID           int       `json:"account_id"`
	CardNumber          string    `json:"card_number"`
	ExpiryDate          string    `json:"expiry_date"`
	CardNumberEncrypted string    `json:"-"`
	ExpiryDateEncrypted string    `json:"-"`
	CVVHash             string    `json:"-"`
	HMAC                string    `json:"-"`
	CreatedAt           time.Time `json:"created_at"`
}
