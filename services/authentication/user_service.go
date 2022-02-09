package authentication_service

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/asishshaji/dotedon-api/models"
	authentication_repository "github.com/asishshaji/dotedon-api/repositories/authentication"
	"github.com/asishshaji/dotedon-api/utils"
	"github.com/dgrijalva/jwt-go"
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

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		authService.l.Println(err)
		return err
	}

	user.Password = hashedPassword

	err = authService.userRepo.RegisterUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (authService AuthenticationService) LoginUser(ctx context.Context, username, password string) (string, error) {

	user := authService.userRepo.GetUserByUsername(ctx, username)
	if user == nil {
		return "", models.ErrNoUserExists
	}

	authenticate := utils.CheckPasswordHash(password, user.Password)
	if !authenticate {
		return "", models.ErrInvalidCredentials
	}

	tokenMethod := jwt.New(jwt.SigningMethodHS256)

	claims := tokenMethod.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := tokenMethod.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		authService.l.Println(err)
		return "", err
	}

	return t, nil

}
