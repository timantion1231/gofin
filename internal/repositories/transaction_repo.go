package repositories

import (
	"database/sql"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db}
}

func (r *TransactionRepository) SumByType(userID int, txType string, month time.Month, year int) (float64, error) {
	var sum float64
	query := `
		SELECT COALESCE(SUM(t.amount),0)
		FROM transactions t
		JOIN accounts a ON t.account_id = a.id
		WHERE a.user_id = $1 AND t.type = $2 AND EXTRACT(MONTH FROM t.created_at) = $3 AND EXTRACT(YEAR FROM t.created_at) = $4
	`
	err := r.db.QueryRow(query, userID, txType, int(month), year).Scan(&sum)
	return sum, err
}
