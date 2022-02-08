package authentication_service

import (
	"context"
	"log"

	"github.com/asishshaji/dotedon-api/models"
	authentication_repository "github.com/asishshaji/dotedon-api/repositories/authentication"
)

type AuthenticationService struct {
	userRepo authentication_repository.IUserAuthenticationRepository
	l        *log.Logger
}

func NewAuthenticationService(l *log.Logger, uR authentication_repository.IUserAuthenticationRepository) IAuthenticationService {
	return AuthenticationService{
		userRepo: uR,
		l:        l,
	}
}

func (authService AuthenticationService) RegisterUser(ctx context.Context, user *models.User) error {
	userExists := authService.userRepo.CheckUserExistsWithUserName(ctx, user.Username)
	if userExists {
		return models.ErrUserExists
	}

	err := user.ValidateUser()
	if err != nil {
		authService.l.Println("Error validating user model", err)
		return err
	}

	err = authService.userRepo.RegisterUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
