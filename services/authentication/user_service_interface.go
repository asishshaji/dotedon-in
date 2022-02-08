package authentication_service

import (
	"context"

	"github.com/asishshaji/dotedon-api/models"
)

type IAuthenticationService interface {
	RegisterUser(ctx context.Context, user *models.User) error
}
