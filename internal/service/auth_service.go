package service

import (
	"errors"

	"github.com/TakedaB/akademi-api/internal/repository"
)

var ErrInvalidCredentials = errors.New("email ou senha inválidos")

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
