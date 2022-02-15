package admin_service

import "context"

type IAdminService interface {
	Login(ctx context.Context, username, password string) (string, error)
}
