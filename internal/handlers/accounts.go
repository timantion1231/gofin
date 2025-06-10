package handlers

import (
	"encoding/json"
	"financeAppAPI/internal/middleware"
	"financeAppAPI/internal/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateAccountHandler(accountService *services.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)
		account, err := accountService.CreateAccount(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(account)
	}
}

func DepositHandler(accountService *services.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			AccountID int     `json:"account_id"`
			Amount    float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}
		userID := r.Context().Value(middleware.UserIDKey).(int)
		account, err := accountService.AccountRepo.GetAccountByID(req.AccountID)
		if err != nil || account.UserID != userID {
			http.Error(w, "Нет доступа к счету", http.StatusForbidden)
			return
		}
		if err := accountService.Deposit(req.AccountID, req.Amount); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Баланс пополнен"})
	}
}

func WithdrawHandler(accountService *services.AccountService) http.HandlerFunc {
	fmt.Println("WithdrawHandler called")
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			AccountID int     `json:"account_id"`
			Amount    float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}
		if req.AccountID <= 0 || req.Amount <= 0 {
			http.Error(w, "Некорректные данные", http.StatusBadRequest)
			return
		}
		userID := r.Context().Value(middleware.UserIDKey).(int)
		account, err := accountService.AccountRepo.GetAccountByID(req.AccountID)
		if err != nil || account.UserID != userID {
			http.Error(w, "Нет доступа к счету", http.StatusForbidden)
			return
		}
		if err := accountService.Withdraw(req.AccountID, req.Amount); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Средства списаны"})
	}
}

func PredictBalanceHandler(analyticsService *services.AnalyticsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID, err := strconv.Atoi(vars["accountId"])
		if err != nil {
			http.Error(w, "Некорректный account_id", http.StatusBadRequest)
			return
		}
		daysStr := r.URL.Query().Get("days")
		days := 30
		if daysStr != "" {
			if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
				days = d
			}
		}
		// TODO: проверить права пользователя на счет
		balance, err := analyticsService.ForecastBalance(accountID, days)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]float64{"forecast_balance": balance})
	}
}
