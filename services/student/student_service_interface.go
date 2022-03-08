package student_service

import (
	"context"
	"mime/multipart"

	"github.com/asishshaji/dotedon-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStudentService interface {
	RegisterStudent(ctx context.Context, student *models.StudentDTO) error
	LoginStudent(ctx context.Context, studentname, password string) (string, error)

	GetMentors(ctx context.Context, userid primitive.ObjectID) ([]*models.MentorResponse, error)
	AddMentorToStudent(ctx context.Context, studentId, mentorId primitive.ObjectID) error

	UpdateTaskSubmission(ctx context.Context, taskDto models.TaskSubmissionDTO, userID primitive.ObjectID, file multipart.File) error
	GetTasks(ctx context.Context, studentId primitive.ObjectID) ([]models.TaskStudentResponse, error)
	CreateTaskSubmission(ctx context.Context, task models.TaskSubmissionDTO, userID primitive.ObjectID, file multipart.File) error
}
