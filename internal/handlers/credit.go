package handlers

import (
	"encoding/json"
	"financeAppAPI/internal/middleware"
	"financeAppAPI/internal/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateCreditHandler(creditService *services.CreditService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			AccountID int     `json:"account_id"`
			Amount    float64 `json:"amount"`
			Rate      float64 `json:"rate"`
			Months    int     `json:"months"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}
		userID := r.Context().Value(middleware.UserIDKey).(int)
		credit, err := creditService.CreateCredit(userID, req.AccountID, req.Amount, req.Rate, req.Months)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(credit)
	}
}

func GetCreditScheduleHandler(creditService *services.CreditService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		creditID, err := strconv.Atoi(vars["creditId"])
		if err != nil {
			http.Error(w, "Некорректный credit_id", http.StatusBadRequest)
			return
		}
		// TODO: проверить права пользователя на кредит
		payments, err := creditService.GetPayments(creditID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payments)
	}
}
