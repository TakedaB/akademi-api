package service

import (
	"errors"

	"github.com/TakedaB/akademi-api/internal/model"
	"github.com/TakedaB/akademi-api/internal/repository"
)

var ErrInvalidCredentials = errors.New("email ou senha inválidos")
var ErrUserNotFound = errors.New("usuário não encontrado")

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if !CheckPassword(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	return GenerateToken(user.ID, string(user.Role))
}

func (s *AuthService) GetProfile(userID string) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
