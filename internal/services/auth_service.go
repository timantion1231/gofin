package services

import (
	"errors"
	"financeAppAPI/internal/models"
	"financeAppAPI/internal/repositories"
	"financeAppAPI/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo}
}

func (s *AuthService) Register(username, email, password string) (*models.User, error) {
	if _, err := s.userRepo.GetUserByEmail(email); err == nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}
	if _, err := s.userRepo.GetUserByUsername(username); err == nil {
		return nil, errors.New("пользователь с таким username уже существует")
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("пользователь не найден")
	}
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("неверный пароль")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(utils.JWTSecret()))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
