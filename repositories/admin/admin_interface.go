package admin_repository

import (
	"context"

	"github.com/asishshaji/dotedon-api/models"
)

type IAdminRepository interface {
	GetAdmin(ctx context.Context, username string) (*models.Admin, error)
}
