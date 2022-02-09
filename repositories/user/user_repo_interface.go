package user_repository

import (
	"context"

	"github.com/asishshaji/dotedon-api/models"
)

type IUserRepository interface {
	RegisterUser(context.Context, *models.User) error
	CheckUserExistsWithUserName(ctx context.Context, username string) bool
	GetUserByUsername(ctx context.Context, username string) *models.User
	GetMentors(ctx context.Context) ([]*models.MentorDTO, error)
}
