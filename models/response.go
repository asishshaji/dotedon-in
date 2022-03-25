package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type StudentResponse struct {
	ID               primitive.ObjectID `json:"id"`
	Email            string             `json:"email"`
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
	ID           primitive.ObjectID `json:"_id"`
	Name         string             `json:"name"`
	Title        string             `json:"title"`
	Organization string             `json:"organization"`
	Domain       string             `json:"domain"`
	CreatedAt    primitive.DateTime `json:"created_at"`
	Image        string             `json:"image"`
	Videos       []Videos           `json:"videos"`
}

type TaskStudentResponse struct {
	ID        primitive.ObjectID `json:"_id"`
	Semester  string             `json:"semester"`
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
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	FileURL string             `json:"fileurl"`
	Status  Status             `json:"status"`
	Comment string             `json:"comment"`
	Task    Task               `json:"task"`
	Student StudentTaskRespone `json:"student"`
	// UpdatedAt primitive.DateTime `json:"updated_at"`
}

type StudentLoginResponse struct {
	Student StudentResponse `json:"student"`
	Jwt     string          `json:"jwt"`
}

type Data struct {
	Domains  []string `json:"domains"`
	Colleges []string `json:"colleges"`
	Courses  []string `json:"courses"`
}

type NotificationResponse struct {
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	CreatedAt primitive.DateTime `json:"created_at"`
}
