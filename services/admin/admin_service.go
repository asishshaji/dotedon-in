package admin_service

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/asishshaji/dotedon-api/models"
	admin_repository "github.com/asishshaji/dotedon-api/repositories/admin"
	"github.com/asishshaji/dotedon-api/utils"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminService struct {
	l         *log.Logger
	adminRepo admin_repository.IAdminRepository
}

func NewAdminService(l *log.Logger, adminRepo admin_repository.IAdminRepository) IAdminService {
	return AdminService{
		l:         l,
		adminRepo: adminRepo,
	}
}

func (aS AdminService) Login(ctx context.Context, username, password string) (string, error) {
	admin, err := aS.adminRepo.GetAdmin(ctx, username)

	if err != nil {
		return "", models.ErrNoAdminWithUsername
	}

	authenticate := utils.CheckPasswordHash(password, admin.Password)

	if !authenticate {
		return "", models.ErrInvalidCredentials
	}

	adminClaims := &models.AdminJWTClaims{
		admin.ID,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	tokenMethod := jwt.NewWithClaims(jwt.SigningMethodHS256, adminClaims)
	t, err := tokenMethod.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		aS.l.Println(err)
		return "", err
	}

	return t, nil
}

func (aS AdminService) AddTask(ctx context.Context, task models.Task, creatorID primitive.ObjectID) error {
	task.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	task.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	task.CreatorID = creatorID
	task.Id = primitive.NewObjectIDFromTimestamp(time.Now())

	err := aS.adminRepo.AddTask(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (aS AdminService) UpdateTask(ctx context.Context, task models.TaskUpdate) error {
	task.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	err := aS.adminRepo.UpdateTask(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (aS AdminService) GetTasks(ctx context.Context) ([]models.Task, error) {
	return aS.adminRepo.GetTasks(ctx)
}
