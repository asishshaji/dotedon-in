package user_service

import (
	"context"

	"github.com/asishshaji/dotedon-api/models"
)

type IUserService interface {
	RegisterUser(ctx context.Context, user *models.User) error
	LoginUser(ctx context.Context, username, password string) (string, error)
	GetMentors(ctx context.Context) ([]*models.MentorResponse, error)
}
