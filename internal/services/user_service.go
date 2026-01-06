package services

import (
	"github.com/stewicca/angagrar-backend/internal/models"
	"github.com/stewicca/angagrar-backend/internal/repositories"
)

type UserService interface {
	GetProfile(userID uint) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(userID uint) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}
