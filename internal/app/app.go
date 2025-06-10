package app

import (
	"database/sql"
	"financeAppAPI/internal/config"
	"financeAppAPI/internal/handlers"
	"financeAppAPI/internal/middleware"
	"financeAppAPI/internal/repositories"
	"financeAppAPI/internal/services"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func InitApp() error {
	logrus.SetLevel(logrus.InfoLevel)
	cfg := config.LoadConfig()
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	userRepo := repositories.NewUserRepository(db)
	accountRepo := repositories.NewAccountRepository(db)
	cardRepo := repositories.NewCardRepository(db)
	creditRepo := repositories.NewCreditRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	authService := services.NewAuthService(userRepo)
	accountService := services.NewAccountService(accountRepo)
	cardService := services.NewCardService(cardRepo)
	cardService.SetAccountRepo(accountRepo)
	creditService := services.NewCreditService(creditRepo, accountRepo)
	analyticsService := services.NewAnalyticsService(accountRepo, transactionRepo, creditRepo)

	r := mux.NewRouter()

	r.HandleFunc("/register", handlers.RegisterHandler(authService)).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler(authService)).Methods("POST")

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/accounts", handlers.CreateAccountHandler(accountService)).Methods("POST")
	protected.HandleFunc("/accounts/{accountId}/predict", handlers.PredictBalanceHandler(analyticsService)).Methods("GET")
	protected.HandleFunc("/cards", handlers.CreateCardHandler(cardService)).Methods("POST")
	protected.HandleFunc("/cards", handlers.GetCardsHandler(cardService)).Methods("GET")
	protected.HandleFunc("/transfers", handlers.TransferHandler(accountService)).Methods("POST")
	protected.HandleFunc("/accounts/deposit", handlers.DepositHandler(accountService)).Methods("POST")
	protected.HandleFunc("/accounts/withdraw", handlers.WithdrawHandler(accountService)).Methods("POST")
	protected.HandleFunc("/credits", handlers.CreateCreditHandler(creditService)).Methods("POST")
	protected.HandleFunc("/credits/{creditId}/schedule", handlers.GetCreditScheduleHandler(creditService)).Methods("GET")
	protected.HandleFunc("/analytics", handlers.AnalyticsHandler(analyticsService)).Methods("GET")

	go func() {
		for {
			creditService.ProcessPayments()
			time.Sleep(12 * time.Hour)
		}
	}()

	logrus.Info("Сервер запущен на :8080")
	return http.ListenAndServe(":8080", r)
}
