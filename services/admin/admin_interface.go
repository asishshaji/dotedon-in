package admin_service

import (
	"context"

	"github.com/asishshaji/dotedon-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAdminService interface {
	Login(ctx context.Context, username, password string) (string, error)
	AddTask(ctx context.Context, task models.Task, creatorID primitive.ObjectID) error
	UpdateTask(ctx context.Context, task models.TaskUpdate) error
	DeleteTask(c context.Context, taskId primitive.ObjectID) error
	GetTasks(ctx context.Context) ([]models.Task, error)
	GetUsers(ctx context.Context) ([]models.StudentResponse, error)
}
