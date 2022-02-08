package authentication_service

import (
	"context"

	"github.com/asishshaji/dotedon-api/models"
)

type IAuthenticationService interface {
	RegisterUser(ctx context.Context, user *models.User) error
	LoginUser(ctx context.Context, username, password string) (string, error)
}
