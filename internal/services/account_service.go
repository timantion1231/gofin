package services

import (
	"errors"
	"financeAppAPI/internal/models"
	"financeAppAPI/internal/repositories"
	"time"
)

type AccountService struct {
	AccountRepo *repositories.AccountRepository
}

func NewAccountService(accountRepo *repositories.AccountRepository) *AccountService {
	return &AccountService{AccountRepo: accountRepo}
}

func (s *AccountService) CreateAccount(userID int) (*models.Account, error) {
	account := &models.Account{
		UserID:        userID,
		AccountNumber: generateAccountNumber(),
		Balance:       0.0,
		CreatedAt:     time.Now(),
	}
	if err := s.AccountRepo.CreateAccount(account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountService) TransferFunds(fromAccountID, toAccountID int, amount float64) error {
	if fromAccountID == toAccountID {
		return errors.New("перевод на тот же счет невозможен")
	}
	fromAccount, err := s.AccountRepo.GetAccountByID(fromAccountID)
	if err != nil {
		return errors.New("исходный счет не найден")
	}
	if fromAccount.Balance < amount {
		return errors.New("недостаточно средств")
	}
	toAccount, err := s.AccountRepo.GetAccountByID(toAccountID)
	if err != nil {
		return errors.New("целевой счет не найден")
	}
	tx, err := s.AccountRepo.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`UPDATE accounts SET balance = $1 WHERE id = $2`, fromAccount.Balance-amount, fromAccountID); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE accounts SET balance = $1 WHERE id = $2`, toAccount.Balance+amount, toAccountID); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *AccountService) Deposit(accountID int, amount float64) error {
	if amount <= 0 {
		return errors.New("сумма должна быть положительной")
	}
	account, err := s.AccountRepo.GetAccountByID(accountID)
	if err != nil {
		return errors.New("счет не найден")
	}
	return s.AccountRepo.UpdateBalance(accountID, account.Balance+amount)
}

func (s *AccountService) Withdraw(accountID int, amount float64) error {
	if amount <= 0 {
		return errors.New("сумма должна быть положительной")
	}
	account, err := s.AccountRepo.GetAccountByID(accountID)
	if err != nil {
		return errors.New("счет не найден")
	}
	if account.Balance < amount {
		return errors.New("недостаточно средств")
	}
	return s.AccountRepo.UpdateBalance(accountID, account.Balance-amount)
}

func generateAccountNumber() string {
	return "ACC" + time.Now().Format("20060102150405")
}
