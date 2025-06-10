package repositories

import (
	"database/sql"
	"time"

	"financeAppAPI/internal/models"
)

type CreditRepository struct {
	db *sql.DB
}

func NewCreditRepository(db *sql.DB) *CreditRepository {
	return &CreditRepository{db}
}

func (r *CreditRepository) CreateCredit(credit *models.Credit) error {
	query := `INSERT INTO credits (user_id, account_id, amount, rate, months, monthly_payment, created_at, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`
	return r.db.QueryRow(query, credit.UserID, credit.AccountID, credit.Amount, credit.Rate, credit.Months, credit.MonthlyPayment, credit.CreatedAt, credit.Status).Scan(&credit.ID)
}

func (r *CreditRepository) CreatePayment(payment *models.CreditPayment) error {
	query := `INSERT INTO credit_payments (credit_id, due_date, amount, paid, fine)
		VALUES ($1,$2,$3,$4,$5) RETURNING id`
	return r.db.QueryRow(query, payment.CreditID, payment.DueDate, payment.Amount, payment.Paid, payment.Fine).Scan(&payment.ID)
}

func (r *CreditRepository) GetActiveCredits() ([]*models.Credit, error) {
	rows, err := r.db.Query(`SELECT id, user_id, account_id, amount, rate, months, monthly_payment, created_at, status FROM credits WHERE status = 'active'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var credits []*models.Credit
	for rows.Next() {
		c := &models.Credit{}
		if err := rows.Scan(&c.ID, &c.UserID, &c.AccountID, &c.Amount, &c.Rate, &c.Months, &c.MonthlyPayment, &c.CreatedAt, &c.Status); err != nil {
			return nil, err
		}
		credits = append(credits, c)
	}
	return credits, nil
}

func (r *CreditRepository) GetPaymentsByCreditID(creditID int) ([]*models.CreditPayment, error) {
	rows, err := r.db.Query(`SELECT id, credit_id, due_date, amount, paid, paid_at, fine FROM credit_payments WHERE credit_id = $1`, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []*models.CreditPayment
	for rows.Next() {
		p := &models.CreditPayment{}
		var paidAt sql.NullTime
		if err := rows.Scan(&p.ID, &p.CreditID, &p.DueDate, &p.Amount, &p.Paid, &paidAt, &p.Fine); err != nil {
			return nil, err
		}
		if paidAt.Valid {
			p.PaidAt = &paidAt.Time
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *CreditRepository) MarkPaymentPaid(paymentID int, paidAt *sql.NullTime) error {
	query := `UPDATE credit_payments SET paid = true, paid_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, paidAt, paymentID)
	return err
}

func (r *CreditRepository) AddFine(paymentID int, fine float64) error {
	query := `UPDATE credit_payments SET fine = fine + $1 WHERE id = $2`
	_, err := r.db.Exec(query, fine, paymentID)
	return err
}

func (r *CreditRepository) SumMonthlyPayments(userID int) (float64, error) {
	var sum float64
	query := `
		SELECT COALESCE(SUM(monthly_payment),0)
		FROM credits
		WHERE user_id = $1 AND status = 'active'
	`
	err := r.db.QueryRow(query, userID).Scan(&sum)
	return sum, err
}

func (r *CreditRepository) SumPlannedPayments(accountID int, days int) (float64, error) {
	var sum float64
	query := `
		SELECT COALESCE(SUM(amount),0)
		FROM credit_payments cp
		JOIN credits c ON cp.credit_id = c.id
		WHERE c.account_id = $1 AND cp.paid = false AND cp.due_date <= $2
	`
	until := time.Now().AddDate(0, 0, days)
	err := r.db.QueryRow(query, accountID, until).Scan(&sum)
	return sum, err
}
