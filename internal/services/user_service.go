package services

import (
	"errors"
	"financeAppAPI/internal/models"
	"financeAppAPI/internal/repositories"
	"financeAppAPI/internal/utils"
	"time"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo}
}

func (s *UserService) Register(username, email, password string) (*models.User, error) {
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

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *UserService) CreateUser(user *models.User) error {
	if _, err := s.userRepo.GetUserByEmail(user.Email); err == nil {
		return errors.New("пользователь с таким email уже существует")
	}
	// user.PasswordHash уже должен быть хешем
	user.CreatedAt = time.Now()
	return s.userRepo.CreateUser(user)
}
