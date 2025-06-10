package repositories

import (
	"database/sql"
	"financeAppAPI/internal/models"
)

type AccountRepository struct {
	DB *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (r *AccountRepository) CreateAccount(account *models.Account) error {
	query := `INSERT INTO accounts (user_id, account_number, balance, created_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`
	return r.DB.QueryRow(query, account.UserID, account.AccountNumber, account.Balance, account.CreatedAt).Scan(&account.ID)
}

func (r *AccountRepository) GetAccountByID(id int) (*models.Account, error) {
	account := &models.Account{}
	query := `SELECT id, user_id, account_number, balance, created_at 
              FROM accounts WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(&account.ID, &account.UserID, &account.AccountNumber, &account.Balance, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *AccountRepository) UpdateBalance(accountID int, newBalance float64) error {
	query := `UPDATE accounts SET balance = $1 WHERE id = $2`
	_, err := r.DB.Exec(query, newBalance, accountID)
	return err
}

func (r *AccountRepository) GetUserByAccountID(accountID int) (*models.User, error) {
	var user models.User
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.created_at
		FROM users u
		JOIN accounts a ON a.user_id = u.id
		WHERE a.id = $1
	`
	err := r.DB.QueryRow(query, accountID).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
