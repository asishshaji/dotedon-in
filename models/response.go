package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserResponse struct {
	Username         string             `json:"username"`
	FirstName        string             `json:"first_name"`
	LastName         string             `json:"last_name"`
	MiddleName       string             `json:"middle_name"`
	CreatedAt        primitive.DateTime `json:"created_at"`
	UpdatedAt        primitive.DateTime `json:"updated_at"`
	DOB              string             `json:"dob"`
	Gender           Gender             `json:"gender"`
	PhoneNumber      string             `json:"phone_number"`
	PhoneNumberAlt   string             `json:"phone_number_alt"`
	College          string             `json:"college"`
	Course           string             `json:"course"`
	Specialization   string             `json:"specialization"`
	HasArrears       bool               `json:"has_arrears"`
	Place            string             `json:"place"`
	Semester         string             `json:"semester"`
	District         string             `json:"district"`
	State            string             `json:"state"`
	Country          string             `json:"country"`
	DateOfJoining    string             `json:"date_of_joining"`
	CourseEndingDate string             `json:"course_ending_date"`
	Mentors          []MentorResponse   `json:"mentors"`
}

type MentorResponse struct {
	Name         string             `json:"name"`
	Title        string             `json:"title"`
	Organization string             `json:"orginization"`
	Domain       string             `json:"domain"`
	CreatedAt    primitive.DateTime `json:"created_at"`
}
