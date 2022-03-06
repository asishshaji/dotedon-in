package student_repository

import (
	"context"

	"github.com/asishshaji/dotedon-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStudentRepository interface {
	RegisterStudent(context.Context, *models.Student) error
	CheckStudentExistsWithStudentName(ctx context.Context, username string) bool
	GetStudentByStudentname(ctx context.Context, username string) *models.Student
	GetMentorIDsFollowedByStudent(ctx context.Context, userid primitive.ObjectID) ([]primitive.ObjectID, error)
	GetMentorsNotInIDS(c context.Context, ids []primitive.ObjectID) ([]*models.Mentor, error)
	AddMentorToStudent(ctx context.Context, userId primitive.ObjectID, mentorId primitive.ObjectID) error
	TaskSubmission(ctx context.Context, task models.TaskSubmission) error
	GetTasks(ctx context.Context, typeVar string) ([]models.Task, error)
	GetTaskSubmissions(ctx context.Context, userId primitive.ObjectID) ([]models.TaskSubmission, error)
	GetStudentByID(ctx context.Context, studentID primitive.ObjectID) (*models.Student, error)
}
