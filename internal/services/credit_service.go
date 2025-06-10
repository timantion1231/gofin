package services

import (
	"database/sql"
	"financeAppAPI/internal/models"
	"financeAppAPI/internal/repositories"
	"math"
	"time"
)

type CreditService struct {
	creditRepo  *repositories.CreditRepository
	accountRepo *repositories.AccountRepository
}

func NewCreditService(creditRepo *repositories.CreditRepository, accountRepo *repositories.AccountRepository) *CreditService {
	return &CreditService{creditRepo, accountRepo}
}

func annuityPayment(amount, rate float64, months int) float64 {
	r := rate / 12 / 100
	return amount * r * math.Pow(1+r, float64(months)) / (math.Pow(1+r, float64(months)) - 1)
}

func (s *CreditService) CreateCredit(userID, accountID int, amount, rate float64, months int) (*models.Credit, error) {
	monthly := annuityPayment(amount, rate, months)
	credit := &models.Credit{
		UserID:         userID,
		AccountID:      accountID,
		Amount:         amount,
		Rate:           rate,
		Months:         months,
		MonthlyPayment: monthly,
		CreatedAt:      time.Now(),
		Status:         "active",
	}
	if err := s.creditRepo.CreateCredit(credit); err != nil {
		return nil, err
	}
	for i := 1; i <= months; i++ {
		payment := &models.CreditPayment{
			CreditID: credit.ID,
			DueDate:  time.Now().AddDate(0, i, 0),
			Amount:   monthly,
			Paid:     false,
			Fine:     0,
		}
		if err := s.creditRepo.CreatePayment(payment); err != nil {
			return nil, err
		}
	}
	return credit, nil
}

func (s *CreditService) GetPayments(creditID int) ([]*models.CreditPayment, error) {
	return s.creditRepo.GetPaymentsByCreditID(creditID)
}

func (s *CreditService) ProcessPayments() error {
	credits, err := s.creditRepo.GetActiveCredits()
	if err != nil {
		return err
	}
	emailService := NewEmailService()
	for _, credit := range credits {
		payments, err := s.creditRepo.GetPaymentsByCreditID(credit.ID)
		if err != nil {
			continue
		}
		for _, p := range payments {
			if !p.Paid && time.Now().After(p.DueDate) {
				account, err := s.accountRepo.GetAccountByID(credit.AccountID)
				if err != nil {
					continue
				}
				if account.Balance >= p.Amount {
					s.accountRepo.UpdateBalance(account.ID, account.Balance-p.Amount)
					now := time.Now()
					nullTime := sql.NullTime{Time: now, Valid: true}
					s.creditRepo.MarkPaymentPaid(p.ID, &nullTime)
					user, err := s.accountRepo.GetUserByAccountID(account.ID)
					if err == nil {
						emailService.Send(user.Email, "Платеж по кредиту проведен", "Ваш платеж по кредиту успешно проведен.")
					}
				} else {
					fine := p.Amount * 0.10
					s.creditRepo.AddFine(p.ID, fine)
				}
			}
		}
	}
	return nil
}
