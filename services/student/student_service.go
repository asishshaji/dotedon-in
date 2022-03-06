package student_service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/asishshaji/dotedon-api/models"
	student_repository "github.com/asishshaji/dotedon-api/repositories/student"
	"github.com/asishshaji/dotedon-api/utils"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentService struct {
	studentRepo student_repository.IStudentRepository
	l           *log.Logger
	rClient     *redis.Client
}

func NewStudentService(l *log.Logger, uR student_repository.IStudentRepository, rClient *redis.Client) IStudentService {
	return StudentService{
		studentRepo: uR,
		l:           l,
		rClient:     rClient,
	}
}

func (authService StudentService) RegisterStudent(ctx context.Context, userDto *models.StudentDTO) error {

	// TODO make username as index so this call be aboided
	userExists := authService.studentRepo.CheckStudentExistsWithStudentName(ctx, userDto.Username)

	if userExists {
		return models.ErrStudentExists
	}

	user := userDto.ToStudent()
	user.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		authService.l.Println(err)
		return err
	}

	user.Password = hashedPassword
	user.Mentors = make([]primitive.ObjectID, 0)

	err = authService.studentRepo.RegisterStudent(ctx, &user)
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

func (authService StudentService) GetMentors(ctx context.Context, userid primitive.ObjectID) ([]*models.MentorResponse, error) {

	fmt.Println("Hello World")
	mentorIdsFollowedByUser, err := authService.studentRepo.GetMentorIDsFollowedByStudent(ctx, userid)
	authService.l.Println(mentorIdsFollowedByUser)
	if err != nil {
		return nil, err
	}

	// TODO get mentors not in mentorIsFollowedByUser
	mentorDtos, err := authService.studentRepo.GetMentorsNotInIDS(ctx, mentorIdsFollowedByUser)
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
	task.Status = models.ACTIVE
	task.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	task.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	err = sS.studentRepo.TaskSubmission(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (sS StudentService) GetTasks(ctx context.Context, studentId primitive.ObjectID) ([]models.TaskStudentResponse, error) {

	taskStudentResponse := []models.TaskStudentResponse{}

	student, err := sS.studentRepo.GetStudentByID(ctx, studentId)

	if err != nil {
		return nil, err
	}

	taskSubmission, err := sS.studentRepo.GetTaskSubmissions(ctx, studentId)
	if err != nil {
		return nil, err
	}

	tasks, err := sS.studentRepo.GetTasks(ctx, string(student.PreferedType))
	if err != nil {
		return nil, err
	}

	for _, t := range tasks {

		fileUrl, comment, status := getFileAndCommentsForTaskIdAndUserId(taskSubmission, t.Id, studentId)
		taskStudentResponse = append(taskStudentResponse, models.TaskStudentResponse{
			ID:        t.Id,
			Title:     t.Title,
			Detail:    t.Detail,
			Status:    status,
			FileURL:   fileUrl,
			Comments:  comment,
			UpdatedAt: "",
			Semester:  t.Semester,
		})
	}

	return taskStudentResponse, nil
}

func getFileAndCommentsForTaskIdAndUserId(tS []models.TaskSubmission, taskID, userID primitive.ObjectID) (string, string, models.Status) {
	for _, t := range tS {
		if t.TaskId == taskID && t.UserId == userID {
			return t.FileURL, t.Comment, t.Status
		}
	}
	return "", "", models.INACTIVE // havent started yet
}
