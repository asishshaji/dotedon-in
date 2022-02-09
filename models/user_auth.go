package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	Message interface{} `json:"message"`
}

type Gender int32

const (
	MALE   Gender = 1
	FEMALE Gender = 2
)

type User struct {
	Username         string               `json:"username" validate:"required"`
	FirstName        string               `json:"first_name" validate:"required"`
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

var ErrUserExists = fmt.Errorf("User already exists")
var ErrInvalidCredentials = fmt.Errorf("Invalid credentials")
var ErrNoUserExists = fmt.Errorf("No user with given username")

func (user *User) ValidateUser() error {
	validate := validator.New()

	return validate.Struct(user)
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
