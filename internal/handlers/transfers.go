package handlers

import (
	"encoding/json"
	"financeAppAPI/internal/middleware"
	"financeAppAPI/internal/services"
	"net/http"
)

func TransferHandler(accountService *services.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			FromAccountID int     `json:"from_account_id"`
			ToAccountID   int     `json:"to_account_id"`
			Amount        float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}
		userID := r.Context().Value(middleware.UserIDKey).(int)
		account, err := accountService.AccountRepo.GetAccountByID(req.FromAccountID)
		if err != nil || account.UserID != userID {
			http.Error(w, "Нет доступа к счету", http.StatusForbidden)
			return
		}
		if req.Amount <= 0 {
			http.Error(w, "Сумма должна быть положительной", http.StatusBadRequest)
			return
		}
		if req.FromAccountID <= 0 || req.ToAccountID <= 0 {
			http.Error(w, "Некорректный account_id", http.StatusBadRequest)
			return
		}
		err = accountService.TransferFunds(req.FromAccountID, req.ToAccountID, req.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Перевод выполнен успешно"})
	}
}
