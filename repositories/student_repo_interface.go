package student_repository

import (
	"context"

	"github.com/asishshaji/dotedon-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStudentRepository interface {
	RegisterStudent(context.Context, *models.Student) error
	CheckStudentExistsWithEmail(ctx context.Context, email string) bool
	GetStudentByEmail(ctx context.Context, email string) *models.Student
	UpdateStudent(ctx context.Context, student models.Student) error

	GetMentorIDsFollowedByStudent(ctx context.Context, userid primitive.ObjectID) ([]primitive.ObjectID, error)
	GetMentorsNotInIDS(c context.Context, ids []primitive.ObjectID) ([]*models.Mentor, error)
	AddMentorToStudent(ctx context.Context, userId primitive.ObjectID, mentorId primitive.ObjectID) error

	UpdateTaskSubmission(ctx context.Context, task models.TaskSubmission) error
	GetTasks(ctx context.Context, typeVar string) ([]models.Task, error)
	GetTaskSubmissions(ctx context.Context, userId primitive.ObjectID) ([]models.TaskSubmission, error)
	CreateTaskSubmission(ctx context.Context, task models.TaskSubmission) error

	GetStudentByID(ctx context.Context, studentID primitive.ObjectID) (*models.Student, error)

	GetDomains(ctx context.Context) ([]models.StaticModel, error)
	GetColleges(ctx context.Context) ([]models.StaticModel, error)
	GetCourses(ctx context.Context) ([]models.StaticModel, error)

	InsertToken(ctx context.Context, tK models.Token) error

	GetNotifications(ctx context.Context, uid primitive.ObjectID) ([]models.NotificationEntity, error)
}
