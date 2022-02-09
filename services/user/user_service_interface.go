package user_service

import (
	"context"

	"github.com/asishshaji/dotedon-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUserService interface {
	RegisterUser(ctx context.Context, user *models.User) error
	LoginUser(ctx context.Context, username, password string) (string, error)
	GetMentors(ctx context.Context) ([]*models.MentorResponse, error)
	AddMentorToUser(ctx context.Context, userId, mentorId primitive.ObjectID) error
}
