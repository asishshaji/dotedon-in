package models

import (
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

type Admin struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `json:"username"`
	password string             `json:"-"`
}

type Student struct {
	ID               primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Email            string               `json:"email" validate:"required" bson:",omitempty"`
	FirstName        string               `json:"first_name" validate:"required" bson:",omitempty"`
	PreferedType     PreferedType         `json:"type" bson:",omitempty"`
	LastName         string               `json:"last_name" validate:"required" bson:",omitempty"`
	MiddleName       string               `json:"middle_name" bson:",omitempty"`
	CreatedAt        primitive.DateTime   `json:"-" bson:",omitempty"`
	UpdatedAt        primitive.DateTime   `json:"-" bson:",omitempty"`
	Password         string               `json:"password" validate:"required,min=4" bson:",omitempty"`
	DOB              string               `json:"dob" validate:"required" bson:",omitempty"`
	Gender           Gender               `json:"gender" validate:"required" bson:",omitempty"`
	PhoneNumber      string               `json:"phone_number" validate:"required" bson:",omitempty"`
	PhoneNumberAlt   string               `json:"phone_number_alt bson:",omitempty""`
	College          string               `json:"college" validate:"required" bson:",omitempty"`
	Course           string               `json:"course" validate:"required" bson:",omitempty"`
	Specialization   string               `json:"specialization" validate:"required" bson:",omitempty"`
	HasArrears       bool                 `json:"has_arrears" validate:"required" bson:",omitempty"`
	Place            string               `json:"place" validate:"required" bson:",omitempty"`
	Semester         string               `json:"semester" validate:"required" bson:",omitempty"`
	District         string               `json:"district" validate:"required" bson:",omitempty"`
	State            string               `json:"state" validate:"required" bson:",omitempty"`
	Country          string               `json:"country" validate:"required" bson:",omitempty"`
	DateOfJoining    string               `json:"date_of_joining" bson:",omitempty"`
	CourseEndingDate string               `json:"course_ending_date" bson:",omitempty"`
	Mentors          []primitive.ObjectID `json:"-" bson:",omitempty"`
}

func (stu Student) ToStudentResponse() StudentResponse {
	return StudentResponse{
		ID:               stu.ID,
		Email:            stu.Email,
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
	}
}

type Students []Student

func (students Students) ToStudentResponse() []StudentResponse {
	studentReponse := []StudentResponse{}

	for _, stu := range students {
		studentReponse = append(studentReponse, StudentResponse{
			ID:               stu.ID,
			Email:            stu.Email,
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

type Mentor struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" validate:"required"`
	Title        string             `json:"title" validate:"required"`
	Organization string             `json:"organization" validate:"required"`
	Domain       string             `json:"domain" validate:"required"`
	Image        string             `json:"image"`
	CreatedAt    primitive.DateTime `bson:",omitempty"`
	UpdatedAt    primitive.DateTime `bson:",omitempty"`
}

func (mentor *Mentor) Validate() error {
	validate := validator.New()

	return validate.Struct(mentor)
}

func (dto *Mentor) ToResponse() *MentorResponse {

	return &MentorResponse{
		ID:           dto.ID,
		Name:         dto.Name,
		Title:        dto.Title,
		Organization: dto.Organization,
		Domain:       dto.Domain,
		CreatedAt:    dto.CreatedAt,
		Image:        dto.Image,
	}
}

type Task struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Semester  string             `json:"semester"`
	Type      string             `json:"type"`  // TYPE CAN BE RETAIL, ED-Tech
	Title     string             `json:"title"` // title of task
	Detail    string             `json:"detail"`
	CreatedAt primitive.DateTime `json:"created_at" bson:",omitempty"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:",omitempty"`
	CreatorID primitive.ObjectID `json:"creator_id"`
}

type TaskSubmission struct {
	UserId    primitive.ObjectID `json:"userid"`
	TaskId    primitive.ObjectID `json:"taskid"`
	Comment   string             `json:"comment"`
	FileURL   string             `json:"fileurl" bson:",omitempty"`
	Feedback  string
	Status    Status             `json:"status"`
	CreatedAt primitive.DateTime `bson:",omitempty"`
	UpdatedAt primitive.DateTime `bson:",omitempty"`
}

type Type struct {
	Name      string
	CreatedOn primitive.DateTime
}

type TT struct {
	Mentors []primitive.ObjectID
}

type StaticModel struct {
	Name      string
	CreatedOn primitive.DateTime `bson:"created_at"`
}

type Token struct {
	UserId primitive.ObjectID `bson:"user_id"`
	Token  string             `bson:"token"`
}

type NotificationEntity struct {
	Title     string
	Content   string
	Image     string
	CreatedAt primitive.DateTime `bson:"created_at"`
	UserId    primitive.ObjectID `bson:"user_id"`
}

type NotificationMessage struct {
	UserToken string
	Contents  map[string]string
	Heading   map[string]string
	CreatedAt primitive.DateTime `bson:"created_at"`
}

func (nE NotificationEntity) ToNotificationResponse() NotificationResponse {
	return NotificationResponse{
		Title:     nE.Title,
		Content:   nE.Content,
		CreatedAt: nE.CreatedAt,
	}
}
