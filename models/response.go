package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type StudentResponse struct {
	ID               primitive.ObjectID `json:"id"`
	Username         string             `json:"username"`
	FirstName        string             `json:"first_name"`
	PreferedType     PreferedType       `json:"type"`
	LastName         string             `json:"last_name"`
	MiddleName       string             `json:"middle_name"`
	CreatedAt        primitive.DateTime `json:"created_at"`
	UpdatedAt        primitive.DateTime `json:"updated_at"`
	DOB              string             `json:"dob"`
	Gender           string             `json:"gender"`
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
}

type MentorResponse struct {
	Name         string             `json:"name"`
	Title        string             `json:"title"`
	Organization string             `json:"orginization"`
	Domain       string             `json:"domain"`
	CreatedAt    primitive.DateTime `json:"created_at"`
}

type TaskStudentResponse struct {
	ID        primitive.ObjectID `json:"_id"`
	Title     string             `json:"title"`
	Detail    string             `json:"detail"`
	Status    Status             `json:"status"`
	FileURL   string             `json:"file_url"`
	Comments  string             `json:"comments"`
	UpdatedAt string             `json:"updated_at"`
}

type StudentTaskRespone struct {
	Username string             `json:"username"`
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
}

type TaskSubmissionsAdminResponse struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Username string             `json:"username"`
	Title    string             `json:"string"`
	Detail   string             `json:"detail"`
	FileURL  string             `json:"fileurl"`
	Status   Status             `json:"status"`
	Comments string             `json:"comment"`
	Task     Task               `json:"task"`
	Student  StudentTaskRespone `json:"student"`
	// UpdatedAt primitive.DateTime `json:"updated_at"`
}
