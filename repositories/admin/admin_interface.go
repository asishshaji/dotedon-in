package admin_repository

import (
	"context"

	"github.com/asishshaji/dotedon-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAdminRepository interface {
	GetAdmin(ctx context.Context, username string) (*models.Admin, error)
	AddTask(ctx context.Context, task models.Task) error
	UpdateTask(ctx context.Context, task models.TaskUpdate) error
	DeleteTask(ctx context.Context, taskId primitive.ObjectID) error
	GetTasks(ctx context.Context) ([]models.Task, error)
	CreateType(ctx context.Context, typeT models.Type) error
	GetUsers(ctx context.Context) (models.Students, error)
	GetTaskSubmissions(c context.Context) ([]models.TaskSubmissionsAdminResponse, error)
	GetTaskSubmissionsForUser(c context.Context, userid primitive.ObjectID) ([]models.TaskSubmissionsAdminResponse, error)
	EditTaskSubmissionStatus(c context.Context, status models.Status, taskid primitive.ObjectID) error
}
