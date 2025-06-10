package handlers

import (
	"encoding/json"
	"financeAppAPI/internal/middleware"
	"financeAppAPI/internal/services"
	"net/http"
	"strconv"
)

func CreateCardHandler(cardService *services.CardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			AccountID int `json:"account_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}
		userID := r.Context().Value(middleware.UserIDKey).(int)
		if !cardService.ValidateAccountOwnership(userID, req.AccountID) {
			http.Error(w, "Account does not belong to the user", http.StatusForbidden)
			return
		}
		card, err := cardService.CreateCard(req.AccountID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(card)
	}
}

func GetCardsHandler(cardService *services.CardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountIDStr := r.URL.Query().Get("account_id")
		accountID, err := strconv.Atoi(accountIDStr)
		if err != nil || accountID <= 0 {
			http.Error(w, "Некорректный account_id", http.StatusBadRequest)
			return
		}
		cards, err := cardService.GetCards(accountID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cards)
	}
}
