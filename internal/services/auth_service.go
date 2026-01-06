package services

import (
	"github.com/google/uuid"
	"github.com/stewicca/angagrar-backend/internal/models"
	"github.com/stewicca/angagrar-backend/internal/repositories"
	"github.com/stewicca/angagrar-backend/pkg/utils"
)

type AuthService interface {
	CreateGuest() (*models.User, string, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) CreateGuest() (*models.User, string, error) {
	guestID := uuid.New().String()

	user := &models.User{
		GuestID: guestID,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(user.ID, guestID, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
