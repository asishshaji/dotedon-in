package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gender int32

const (
	MALE   Gender = 0
	FEMALE Gender = 1
)

type User struct {
	Username         string             `json:"username" validate:"required"`
	FirstName        string             `json:"first_name" validate:"required"`
	LastName         string             `json:"last_name" validate:"required"`
	MiddleName       string             `json:"middle_name"`
	CreatedAt        primitive.DateTime `json:"created_at"`
	UpdatedAt        primitive.DateTime `json:"updated_at"`
	Password         string             `json:"-"`
	DOB              string             `json:"dob" validate:"required"`
	Gender           Gender             `json:"gender" validate:"required"`
	PhoneNumber      string             `json:"phone_number" validate:"required"`
	PhoneNumberAlt   string             `json:"phone_number_alt"`
	College          string             `json:"college" validate:"required"`
	Course           string             `json:"course" validate:"required"`
	Specialization   string             `json:"specialization" validate:"required"`
	HasArrears       bool               `json:"has_arrears" validate:"required"`
	Place            string             `json:"place" validate:"required"`
	Semester         string             `json:"semester" validate:"required"`
	District         string             `json:"district" validate:"required"`
	State            string             `json:"state" validate:"required"`
	Country          string             `json:"country" validate:"required"`
	DateOfJoining    string             `json:"date_of_joining"`
	CourseEndingDate string             `json:"course_ending_date"`
}

var ErrUserExists = fmt.Errorf("User already exists")
var ErrInvalidCredentials = fmt.Errorf("Invalid credentials")
var ErrNoUserExists = fmt.Errorf("No user with given username")

func (user *User) ValidateUser() error {
	validate := validator.New()

	return validate.Struct(user)
}
