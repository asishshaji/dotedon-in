package authentication_repository

import (
	"context"

	"github.com/asishshaji/dotedon-api/models"
)

type IUserAuthenticationRepository interface {
	RegisterUser(context.Context, *models.User) error
	CheckUserExistsWithUserName(ctx context.Context, username string) bool
}
