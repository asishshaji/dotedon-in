package student_repository

import (
	"context"
	"log"

	"github.com/asishshaji/dotedon-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type studentRepo struct {
	studentCollection *mongo.Collection
	mentorCollection  *mongo.Collection
	l                 *log.Logger
}

func NewStudentAuthRepo(l *log.Logger, db *mongo.Database) IStudentRepository {
	return studentRepo{
		studentCollection: db.Collection("students"),
		mentorCollection:  db.Collection("mentors"),
		l:                 l,
	}
}

func (uR studentRepo) RegisterStudent(ctx context.Context, uM *models.Student) error {
	res, err := uR.studentCollection.InsertOne(ctx, uM)

	if err != nil {
		uR.l.Println("Error inserting student")
		return err
	}

	uR.l.Println("Inserted new product with ID : ", res.InsertedID)

	return nil
}

func (uR studentRepo) CheckStudentExistsWithStudentName(ctx context.Context, studentname string) bool {
	err := uR.studentCollection.FindOne(ctx, bson.M{"studentname": studentname}).Err()
	if err == mongo.ErrNoDocuments {
		return false
	}
	return true
}

func (uR studentRepo) GetStudentByStudentname(ctx context.Context, studentname string) *models.Student {
	student := new(models.Student)
	res := uR.studentCollection.FindOne(ctx, bson.M{"studentname": studentname})

	if res.Err() == mongo.ErrNoDocuments {
		uR.l.Println("Invalid studentname and Password")
		return nil
	}

	err := res.Decode(student)
	if err != nil {
		uR.l.Println("Error decoding student")
		return nil
	}

	return student

}

func (uR studentRepo) GetMentors(ctx context.Context) ([]*models.MentorDTO, error) {

	mentors := []*models.MentorDTO{}

	cursor, err := uR.mentorCollection.Find(ctx, bson.M{})

	if err != nil {
		uR.l.Println(err)

		return nil, err
	}

	if err = cursor.All(ctx, &mentors); err != nil {
		uR.l.Println(err)
		return nil, err
	}

	return mentors, nil

}

func (uR studentRepo) AddMentorToStudent(ctx context.Context, studentId primitive.ObjectID, mentorId primitive.ObjectID) error {

	options := bson.M{
		"$addToSet": bson.M{
			"mentors": mentorId,
		},
	}

	res, err := uR.studentCollection.UpdateByID(ctx, studentId, options)
	if err != nil {
		uR.l.Println(err)
		return err
	}

	if res.MatchedCount == 0 {
		uR.l.Println("No student found with id:", studentId)
		return models.ErrNoStudentWithIdExists
	}

	return nil
}