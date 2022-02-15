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
	GetMentors(ctx context.Context) ([]*models.MentorDTO, error)
	AddMentorToStudent(ctx context.Context, userId primitive.ObjectID, mentorId primitive.ObjectID) error
}
