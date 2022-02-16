package student_service

import (
	"context"

	"github.com/asishshaji/dotedon-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStudentService interface {
	RegisterStudent(ctx context.Context, student *models.Student) error
	LoginStudent(ctx context.Context, studentname, password string) (string, error)
	GetMentors(ctx context.Context) ([]*models.MentorResponse, error)
	AddMentorToStudent(ctx context.Context, studentId, mentorId primitive.ObjectID) error
	TaskSubmission(ctx context.Context, taskDto models.TaskSubmissionDTO, userID primitive.ObjectID) error
}
