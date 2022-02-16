package student_service

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/asishshaji/dotedon-api/models"
	student_repository "github.com/asishshaji/dotedon-api/repositories/student"
	"github.com/asishshaji/dotedon-api/utils"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentService struct {
	studentRepo student_repository.IStudentRepository
	l           *log.Logger
}

func NewStudentService(l *log.Logger, uR student_repository.IStudentRepository) IStudentService {
	return StudentService{
		studentRepo: uR,
		l:           l,
	}
}

func (authService StudentService) RegisterStudent(ctx context.Context, user *models.Student) error {

	userExists := authService.studentRepo.CheckStudentExistsWithStudentName(ctx, user.Username)
	if userExists {
		return models.ErrStudentExists
	}

	err := user.ValidateStudent()
	if err != nil {
		authService.l.Println("Error validating user model", err)
		return err
	}

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		authService.l.Println(err)
		return err
	}

	user.Password = hashedPassword
	user.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	user.Mentors = make([]primitive.ObjectID, 0)

	err = authService.studentRepo.RegisterStudent(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (authService StudentService) LoginStudent(ctx context.Context, username, password string) (string, error) {

	user := authService.studentRepo.GetStudentByStudentname(ctx, username)
	if user == nil {
		return "", models.ErrNoStudentExists
	}

	authenticate := utils.CheckPasswordHash(password, user.Password)
	if !authenticate {
		return "", models.ErrInvalidCredentials
	}

	claims := &models.StudentJWTClaims{
		user.ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	tokenMethod := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := tokenMethod.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		authService.l.Println(err)
		return "", err
	}

	return t, nil

}

func (authService StudentService) GetMentors(ctx context.Context) ([]*models.MentorResponse, error) {
	mentorDtos, err := authService.studentRepo.GetMentors(ctx)
	if err != nil {
		return nil, err
	}

	mentorResponses := []*models.MentorResponse{}

	for _, dto := range mentorDtos {
		mentorResponses = append(mentorResponses, dto.ToResponse())
	}

	return mentorResponses, nil
}

func (authService StudentService) AddMentorToStudent(ctx context.Context, userId, mentorId primitive.ObjectID) error {
	err := authService.studentRepo.AddMentorToStudent(ctx, userId, mentorId)
	if err != nil {
		return err
	}
	return nil
}

func (sS StudentService) TaskSubmission(ctx context.Context, taskDto models.TaskSubmissionDTO, userID primitive.ObjectID) error {

	taskObjID, err := primitive.ObjectIDFromHex(taskDto.TaskId)

	if err != nil {
		sS.l.Println(err)
		return err
	}

	task := models.TaskSubmission{}
	task.TaskId = taskObjID
	task.UserId = userID
	task.Comment = taskDto.Comment
	task.FileURL = taskDto.FileURL
	task.Status = taskDto.Status

	err = sS.studentRepo.TaskSubmission(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (sS StudentService) GetTasks(ctx context.Context, studentId primitive.ObjectID) error {

	typeVar := ""

	tasks, err := sS.studentRepo.GetTasks(ctx, typeVar)
	sS.l.Println(tasks)
	if err != nil {
		return err
	}

	return nil
}
