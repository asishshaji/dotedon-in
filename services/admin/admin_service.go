package admin_service

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/asishshaji/dotedon-api/models"
	admin_repository "github.com/asishshaji/dotedon-api/repositories/admin"
	"github.com/asishshaji/dotedon-api/utils"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminService struct {
	l         *log.Logger
	adminRepo admin_repository.IAdminRepository
	rClient   *redis.Client
}

func NewAdminService(l *log.Logger, adminRepo admin_repository.IAdminRepository, rClient *redis.Client) IAdminService {
	return AdminService{
		l:         l,
		adminRepo: adminRepo,
		rClient:   rClient,
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

func (aS AdminService) GetUsers(ctx context.Context) ([]models.StudentResponse, error) {
	studentModels, err := aS.adminRepo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	studentResponse := studentModels.ToStudentResponse()

	return studentResponse, nil
}

func (aS AdminService) DeleteTask(c context.Context, taskId primitive.ObjectID) error {
	return aS.adminRepo.DeleteTask(c, taskId)
}
