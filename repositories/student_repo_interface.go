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
	AddDomainToStudent(ctx context.Context, userID primitive.ObjectID, domain string) error

	GetMentorIDsFollowedByStudent(ctx context.Context, userid primitive.ObjectID) ([]primitive.ObjectID, error)
	GetMentorsNotInIDS(c context.Context, ids []primitive.ObjectID) ([]*models.Mentor, error)
	AddMentorToStudent(ctx context.Context, userId primitive.ObjectID, mentorId primitive.ObjectID) error
	GetMentorByID(ctx context.Context, mentorId primitive.ObjectID) (models.Mentor, error)

	UpdateTaskSubmission(ctx context.Context, task models.TaskSubmission) error
	GetTasksBySemestersAndDomains(ctx context.Context, domains, sems []string) ([]models.Task, error)
	GetTaskByID(ctx context.Context, taskID primitive.ObjectID) (models.Task, error)
	GetTaskSubmissionsBySemesters(ctx context.Context, userId primitive.ObjectID, semesters []string) ([]models.TaskSubmission, error)
	CreateTaskSubmission(ctx context.Context, task models.TaskSubmission) error

	GetStudentByID(ctx context.Context, studentID primitive.ObjectID) (*models.Student, error)
	GetSubmissionCountStat(ctx context.Context, studentID primitive.ObjectID, status models.Status) (int64, error)

	GetDomains(ctx context.Context) ([]models.StaticModel, error)
	GetColleges(ctx context.Context) ([]models.StaticModel, error)
	GetCourses(ctx context.Context) ([]models.StaticModel, error)

	InsertToken(ctx context.Context, tK models.Token) error

	GetNotifications(ctx context.Context, uid primitive.ObjectID) ([]models.NotificationEntity, error)
}
