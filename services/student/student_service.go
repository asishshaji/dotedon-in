package student_service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/asishshaji/dotedon-api/models"
	student_repository "github.com/asishshaji/dotedon-api/repositories"
	file_service "github.com/asishshaji/dotedon-api/services/file"
	"github.com/asishshaji/dotedon-api/utils"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentService struct {
	studentRepo student_repository.IStudentRepository
	l           *log.Logger
	rClient     *redis.Client
	fileService file_service.IFileService
}

func NewStudentService(l *log.Logger, uR student_repository.IStudentRepository, rClient *redis.Client, fileService file_service.IFileService) IStudentService {
	return StudentService{
		studentRepo: uR,
		l:           l,
		rClient:     rClient,
		fileService: fileService,
	}
}

func (authService StudentService) RegisterStudent(ctx context.Context, userDto *models.StudentDTO) error {

	// TODO make username as index so this call be aboided
	userExists := authService.studentRepo.CheckStudentExistsWithEmail(ctx, userDto.Email)

	if userExists {
		return models.ErrStudentExists
	}

	user := userDto.ToStudent()
	user.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	hasedpassword, err := utils.Hashpassword(user.Password)

	if err != nil {
		authService.l.Println(err)
		return err
	}

	user.Password = hasedpassword
	user.Mentors = make([]primitive.ObjectID, 0)

	err = authService.studentRepo.RegisterStudent(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

func (authService StudentService) LoginStudent(ctx context.Context, email, password string) (models.StudentLoginResponse, error) {

	user := authService.studentRepo.GetStudentByEmail(ctx, email)
	if user == nil {
		return models.StudentLoginResponse{}, models.ErrNoStudentExists
	}

	authenticate := utils.CheckpasswordHash(password, user.Password)
	if !authenticate {
		return models.StudentLoginResponse{}, models.ErrInvalidCredentials
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
		return models.StudentLoginResponse{}, err
	}

	return models.StudentLoginResponse{
		Student: user.ToStudentResponse(),
		Jwt:     t,
	}, nil

}

func (sS StudentService) UpdateStudent(ctx context.Context, student models.StudentDTO) error {
	s := student.ToStudent()
	fmt.Println(s.ID)
	return sS.studentRepo.UpdateStudent(ctx, s)
}
func (authService StudentService) GetMentors(ctx context.Context, userid primitive.ObjectID) ([]*models.MentorResponse, error) {

	mentorIdsFollowedByUser, err := authService.studentRepo.GetMentorIDsFollowedByStudent(ctx, userid)
	if err != nil {
		return nil, err
	}

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

func (sS StudentService) UpdateTaskSubmission(ctx context.Context, taskDto models.TaskSubmissionDTO, userID primitive.ObjectID) error {

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
	task.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	err = sS.studentRepo.UpdateTaskSubmission(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (sS StudentService) CreateTaskSubmission(ctx context.Context, taskDto models.TaskSubmissionDTO, userID primitive.ObjectID) error {
	taskObjID, err := primitive.ObjectIDFromHex(taskDto.TaskId)
	// TODO check if task exists

	if err != nil {
		sS.l.Println(err)
		return err
	}

	task := models.TaskSubmission{}
	task.TaskId = taskObjID
	task.FileURL = taskDto.FileURL
	task.UserId = userID
	task.Comment = taskDto.Comment
	task.Status = models.ACTIVE
	task.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	task.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return sS.studentRepo.CreateTaskSubmission(ctx, task)
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

	studentSem := student.Semester
	var sems []string
	var s string = "S1"
	var idx = 1

	for s != studentSem {
		s = "S"
		s = s + fmt.Sprint(idx)
		idx++
		sems = append(sems, s)
	}
	// TODO GET TASKS BASED ON SEM
	for _, t := range tasks {

		for _, b := range sems {
			if b == t.Semester {
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
		}

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

func (sS StudentService) GetStudent(ctx context.Context, studentId primitive.ObjectID) (models.StudentResponse, error) {
	student, err := sS.studentRepo.GetStudentByID(ctx, studentId)
	if err != nil {
		return models.StudentResponse{}, err
	}
	return student.ToStudentResponse(), nil
}

func (sS StudentService) GetData(ctx context.Context) (models.Data, error) {
	data := models.Data{}

	res, err := sS.rClient.Get(ctx, "static_data").Result()
	if err != redis.Nil {
		sS.l.Println("Getting data from cache")

		err = json.Unmarshal([]byte(res), &data)
		if err != nil {
			return data, err
		}
		return data, nil
	}

	domainEntities, err := sS.studentRepo.GetDomains(ctx)
	if err != nil {
		return data, err
	}
	collegeEntities, err := sS.studentRepo.GetColleges(ctx)
	if err != nil {
		return data, err
	}
	courseEntities, err := sS.studentRepo.GetCourses(ctx)
	if err != nil {
		return data, err
	}

	for _, v := range domainEntities {
		data.Domains = append(data.Domains, v.Name)
	}

	for _, v := range collegeEntities {
		data.Colleges = append(data.Colleges, v.Name)
	}

	for _, v := range courseEntities {
		data.Courses = append(data.Courses, v.Name)

	}

	b, err := json.Marshal(&data)
	if err != nil {
		return data, err
	}

	r := sS.rClient.Set(ctx, "static_data", b, time.Hour*2)
	if r.Err() != nil {
		return data, r.Err()
	}

	return data, nil
}

func (sS StudentService) UploadFile(ctx context.Context, file multipart.File) (string, error) {

	return sS.fileService.UploadFile(ctx, file)
}

func (sS StudentService) InsertToken(ctx context.Context, tK models.TokenDto, uId primitive.ObjectID) error {
	t := models.Token{
		UserId: uId,
		Token:  tK.Token,
	}
	return sS.studentRepo.InsertToken(ctx, t)
}
