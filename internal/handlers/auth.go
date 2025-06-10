package handlers

import (
    "encoding/json"
    "net/http"
    "financeAppAPI/internal/services"
    "financeAppAPI/internal/utils"
)

func RegisterHandler(authService *services.AuthService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username string `json:"username"`
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
            return
        }
        if !utils.IsValidEmail(req.Email) {
            http.Error(w, "Неверный email", http.StatusBadRequest)
            return
        }
        if !utils.IsValidPassword(req.Password) {
            http.Error(w, "Пароль слишком простой", http.StatusBadRequest)
            return
        }
        user, err := authService.Register(req.Username, req.Email, req.Password)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    }
}

func LoginHandler(authService *services.AuthService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
            return
        }
        token, err := authService.Login(req.Email, req.Password)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"token": token})
    }
}