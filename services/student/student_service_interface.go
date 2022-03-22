package student_service

import (
	"context"
	"mime/multipart"

	"github.com/asishshaji/dotedon-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStudentService interface {
	RegisterStudent(ctx context.Context, student *models.StudentDTO) error
	LoginStudent(ctx context.Context, studentname, password string) (models.StudentLoginResponse, error)
	GetStudent(ctx context.Context, studentId primitive.ObjectID) (models.StudentResponse, error)
	UpdateStudent(ctx context.Context, student models.StudentDTO) error

	GetMentors(ctx context.Context, userid primitive.ObjectID) ([]*models.MentorResponse, error)
	AddMentorToStudent(ctx context.Context, studentId, mentorId primitive.ObjectID) error

	UpdateTaskSubmission(ctx context.Context, taskDto models.TaskSubmissionDTO, userID primitive.ObjectID) error
	GetTasks(ctx context.Context, studentId primitive.ObjectID) ([]models.TaskStudentResponse, error)
	CreateTaskSubmission(ctx context.Context, task models.TaskSubmissionDTO, userID primitive.ObjectID) error

	GetData(ctx context.Context) (models.Data, error)
	UploadFile(ctx context.Context, file multipart.File) (string, error)

	InsertToken(ctx context.Context, tK models.TokenDto, uId primitive.ObjectID) error
}
