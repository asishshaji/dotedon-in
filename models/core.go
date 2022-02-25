package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentJWTClaims struct {
	StudentId primitive.ObjectID
	jwt.StandardClaims
}

type AdminJWTClaims struct {
	AdminId primitive.ObjectID
	IsAdmin bool
	jwt.StandardClaims
}

type Response struct {
	Message interface{} `json:"message"`
}

type Gender int32

func (g Gender) String() string {
	switch g {
	case MALE:
		return "male"
	case FEMALE:
		return "female"
	}
	return ""
}

const (
	MALE   Gender = 1
	FEMALE Gender = 2
)

type Admin struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `json:"username"`
	Password string             `json:"-"`
}

type PreferedType string

const (
	EDTECH PreferedType = "ed-tech"
	RETAIL PreferedType = "retail"
)

type Student struct {
	ID               primitive.ObjectID   `bson:"_id" json:"_id"`
	Username         string               `json:"username" validate:"required"`
	FirstName        string               `json:"first_name" validate:"required"`
	PreferedType     PreferedType         `json:"type"`
	LastName         string               `json:"last_name" validate:"required"`
	MiddleName       string               `json:"middle_name"`
	CreatedAt        primitive.DateTime   `json:"-"`
	UpdatedAt        primitive.DateTime   `json:"-"`
	Password         string               `json:"password" validate:"required,min=4"`
	DOB              string               `json:"dob" validate:"required"`
	Gender           Gender               `json:"gender" validate:"required"`
	PhoneNumber      string               `json:"phone_number" validate:"required"`
	PhoneNumberAlt   string               `json:"phone_number_alt"`
	College          string               `json:"college" validate:"required"`
	Course           string               `json:"course" validate:"required"`
	Specialization   string               `json:"specialization" validate:"required"`
	HasArrears       bool                 `json:"has_arrears" validate:"required"`
	Place            string               `json:"place" validate:"required"`
	Semester         string               `json:"semester" validate:"required"`
	District         string               `json:"district" validate:"required"`
	State            string               `json:"state" validate:"required"`
	Country          string               `json:"country" validate:"required"`
	DateOfJoining    string               `json:"date_of_joining"`
	CourseEndingDate string               `json:"course_ending_date"`
	Mentors          []primitive.ObjectID `json:"-"`
}

type Students []Student

func (students Students) ToStudentResponse() []StudentResponse {
	studentReponse := []StudentResponse{}

	for _, stu := range students {
		studentReponse = append(studentReponse, StudentResponse{
			ID:               stu.ID,
			Username:         stu.Username,
			FirstName:        stu.FirstName,
			PreferedType:     stu.PreferedType,
			LastName:         stu.LastName,
			MiddleName:       stu.MiddleName,
			CreatedAt:        stu.CreatedAt,
			UpdatedAt:        stu.UpdatedAt,
			DOB:              stu.DOB,
			Gender:           Gender(stu.Gender).String(),
			PhoneNumber:      stu.PhoneNumber,
			PhoneNumberAlt:   stu.PhoneNumberAlt,
			College:          stu.College,
			Course:           stu.Course,
			Specialization:   stu.Specialization,
			HasArrears:       stu.HasArrears,
			Place:            stu.Place,
			Semester:         stu.Semester,
			District:         stu.District,
			State:            stu.State,
			Country:          stu.Country,
			DateOfJoining:    stu.DateOfJoining,
			CourseEndingDate: stu.CourseEndingDate,
		})
	}

	return studentReponse

}

var ErrStudentExists = fmt.Errorf("Student already exists")
var ErrInvalidCredentials = fmt.Errorf("Invalid credentials")
var ErrNoStudentExists = fmt.Errorf("No Student with given studentname")
var ErrNoStudentWithIdExists = fmt.Errorf("No Student with id exists")
var ErrParsingStudent = fmt.Errorf("Error parsing student data from database")

var ErrNoAdminWithUsername = fmt.Errorf("No admin with username exists")

func (Student *Student) ValidateStudent() error {
	validate := validator.New()

	return validate.Struct(Student)
}

type MentorDTO struct {
	Name         string             `json:"name" validate:"required"`
	Title        string             `json:"title" validate:"required"`
	Organization string             `json:"orginization" validate:"required"`
	Domain       string             `json:"domain" validate:"required"`
	CreatedAt    primitive.DateTime `json:"-"`
	UpdatedAt    primitive.DateTime `json:"-"`
}

func (dto *MentorDTO) ToResponse() *MentorResponse {

	return &MentorResponse{
		Name:         dto.Name,
		Title:        dto.Title,
		Organization: dto.Organization,
		Domain:       dto.Organization,
		CreatedAt:    dto.CreatedAt,
	}
}

type Task struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Semester  string             `json:"semester"`
	Type      string             `json:"type"`  // TYPE CAN BE RETAIL, ED-Tech
	Title     string             `json:"title"` // title of task
	Detail    string             `json:"detail"`
	CreatedAt primitive.DateTime `json:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at"`
	CreatorID primitive.ObjectID `json:"creator_id"`
}

type TaskUpdate struct {
	Id        primitive.ObjectID `json:"id"`
	Semester  string             `json:"semester" validate:"required"`
	Type      string             `json:"type" validate:"required"`
	Title     string             `json:"title" validate:"required"` // title of task
	Detail    string             `json:"detail" validate:"required"`
	UpdatedAt primitive.DateTime `json:"updated_at"`
}

func (task *TaskUpdate) Validate() error {
	validate := validator.New()

	return validate.Struct(task)
}

type Status string

const (
	ACTIVE    Status = "active"    // submitted
	COMPLETED Status = "completed" // verified by admin
	INACTIVE  Status = "inactive"  // not started
	REJECTED  Status = "rejected"
)

func (s Status) String() string {
	switch s {
	case ACTIVE:
		return "active"
	case COMPLETED:
		return "completed"
	case INACTIVE:
		return "inactive"
	case REJECTED:
		return "rejected"
	}
	return ""
}

type TaskSubmissionDTO struct {
	TaskId  string `json:"task_id"`
	Comment string `json:"comment"`
	FileURL string `json:"file_url"`
}

type TaskSubmission struct {
	UserId    primitive.ObjectID `json:"userid"`
	TaskId    primitive.ObjectID `json:"taskid"`
	Comment   string             `json:"comment"`
	FileURL   string             `json:"fileurl"`
	Status    Status             `json:"status"`
	CreatedAt primitive.DateTime `json:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at"`
}

var ErrTaskSubmissionExists = fmt.Errorf("Task submission already exists")

type Type struct {
	Name      string
	CreatedOn primitive.DateTime
}

var ErrNoValidRecordFound = fmt.Errorf("No valid document found")
