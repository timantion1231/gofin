package services

import (
	"financeAppAPI/internal/repositories"
	"time"
)

type AnalyticsService struct {
	accountRepo     *repositories.AccountRepository
	transactionRepo *repositories.TransactionRepository
	creditRepo      *repositories.CreditRepository
}

func NewAnalyticsService(accountRepo *repositories.AccountRepository, transactionRepo *repositories.TransactionRepository, creditRepo *repositories.CreditRepository) *AnalyticsService {
	return &AnalyticsService{accountRepo, transactionRepo, creditRepo}
}

func (s *AnalyticsService) GetMonthlyStats(userID int, month time.Month, year int) (income, expense float64, err error) {
	income, err = s.transactionRepo.SumByType(userID, "deposit", month, year)
	if err != nil {
		return
	}
	expense, err = s.transactionRepo.SumByType(userID, "withdraw", month, year)
	return
}

func (s *AnalyticsService) GetCreditLoad(userID int) (float64, error) {
	totalIncome, err := s.transactionRepo.SumByType(userID, "deposit", time.Now().Month(), time.Now().Year())
	if err != nil {
		return 0, err
	}
	monthlyCredits, err := s.creditRepo.SumMonthlyPayments(userID)
	if err != nil || totalIncome == 0 {
		return 0, err
	}
	return monthlyCredits / totalIncome, nil
}

func (s *AnalyticsService) ForecastBalance(accountID int, days int) (float64, error) {
	account, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return 0, err
	}
	plannedPayments, err := s.creditRepo.SumPlannedPayments(accountID, days)
	if err != nil {
		return 0, err
	}
	return account.Balance - plannedPayments, nil
}
