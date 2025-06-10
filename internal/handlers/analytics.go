package handlers

import (
	"encoding/json"
	"financeAppAPI/internal/middleware"
	"financeAppAPI/internal/services"
	"net/http"
	"strconv"
	"time"
)

func MonthlyStatsHandler(analyticsService *services.AnalyticsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)
		month := time.Now().Month()
		year := time.Now().Year()
		income, expense, err := analyticsService.GetMonthlyStats(userID, month, year)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]float64{
			"income":  income,
			"expense": expense,
		})
	}
}

func CreditLoadHandler(analyticsService *services.AnalyticsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)
		load, err := analyticsService.GetCreditLoad(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]float64{
			"credit_load": load,
		})
	}
}

func ForecastBalanceHandler(analyticsService *services.AnalyticsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountIDStr := r.URL.Query().Get("account_id")
		daysStr := r.URL.Query().Get("days")
		accountID, _ := strconv.Atoi(accountIDStr)
		days, _ := strconv.Atoi(daysStr)
		if days <= 0 || days > 365 {
			days = 30
		}
		balance, err := analyticsService.ForecastBalance(accountID, days)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]float64{
			"forecast_balance": balance,
		})
	}
}

func AnalyticsHandler(analyticsService *services.AnalyticsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)
		month := time.Now().Month()
		year := time.Now().Year()
		income, expense, err := analyticsService.GetMonthlyStats(userID, month, year)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		creditLoad, err := analyticsService.GetCreditLoad(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"income":      income,
			"expense":     expense,
			"credit_load": creditLoad,
		})
	}
}
