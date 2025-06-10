package repositories

import (
	"database/sql"
	"financeAppAPI/internal/models"
)

type CardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db}
}

func (r *CardRepository) CreateCard(card *models.Card) error {
	query := `INSERT INTO cards (account_id, card_number_encrypted, expiry_date_encrypted, cvv_hash, hmac, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.db.QueryRow(query, card.AccountID, card.CardNumberEncrypted, card.ExpiryDateEncrypted, card.CVVHash, card.HMAC, card.CreatedAt).Scan(&card.ID)
}

func (r *CardRepository) GetCardsByAccountID(accountID int) ([]*models.Card, error) {
	rows, err := r.db.Query(`SELECT id, account_id, card_number_encrypted, expiry_date_encrypted, cvv_hash, hmac, created_at FROM cards WHERE account_id = $1`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cards []*models.Card
	for rows.Next() {
		card := &models.Card{}
		if err := rows.Scan(&card.ID, &card.AccountID, &card.CardNumberEncrypted, &card.ExpiryDateEncrypted, &card.CVVHash, &card.HMAC, &card.CreatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}
