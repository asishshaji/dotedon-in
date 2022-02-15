package admin_repository

import (
	"context"

	"github.com/asishshaji/dotedon-api/models"
)

type IAdminRepository interface {
	GetAdmin(ctx context.Context, username string) (*models.Admin, error)
	AddTask(ctx context.Context, task models.Task) error
	UpdateTask(ctx context.Context, task models.TaskUpdate) error
	GetTasks(ctx context.Context) ([]models.Task, error)
}
