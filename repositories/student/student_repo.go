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
	studentCollection        *mongo.Collection
	mentorCollection         *mongo.Collection
	taskSubmissionCollection *mongo.Collection
	taskCollection           *mongo.Collection
	l                        *log.Logger
}

func NewStudentAuthRepo(l *log.Logger, db *mongo.Database) IStudentRepository {
	return studentRepo{
		studentCollection:        db.Collection("students"),
		mentorCollection:         db.Collection("mentors"),
		taskSubmissionCollection: db.Collection("task_submission"),
		taskCollection:           db.Collection("tasks"),
		l:                        l,
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
	res := uR.studentCollection.FindOne(ctx, bson.M{
		"username": studentname,
	}).Err()
	return res != mongo.ErrNoDocuments
}

func (uR studentRepo) GetStudentByStudentname(ctx context.Context, studentname string) *models.Student {
	student := new(models.Student)
	res := uR.studentCollection.FindOne(ctx, bson.M{"username": studentname})

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

func checkTaskSubmissionExists(ctx context.Context, collection *mongo.Collection, userId, taskId primitive.ObjectID) bool {
	res := collection.FindOne(ctx, bson.M{"userid": userId, "taskid": taskId})

	return res.Err() != mongo.ErrNoDocuments
}

func (sR studentRepo) TaskSubmission(ctx context.Context, task models.TaskSubmission) error {

	if checkTaskSubmissionExists(ctx, sR.taskSubmissionCollection, task.UserId, task.TaskId) {
		sR.l.Println(models.ErrTaskSubmissionExists)
		return models.ErrTaskSubmissionExists
	}

	res, err := sR.taskSubmissionCollection.InsertOne(ctx, task)

	if err != nil {
		sR.l.Println(err)
		return err
	}

	sR.l.Println("Inserted submission with ID", res.InsertedID)

	return nil

}

func (sR studentRepo) GetTasks(ctx context.Context, typeVar string) ([]models.Task, error) {
	tasks := []models.Task{}

	// TODO add filter queries: semester, type
	cursor, err := sR.taskCollection.Find(ctx, bson.M{
		"type": typeVar,
	})

	if err != nil {
		sR.l.Println(err)
		return nil, err
	}

	if err = cursor.All(ctx, &tasks); err != nil {
		sR.l.Println(err)
		return nil, err
	}

	return tasks, nil
}

func (sR studentRepo) GetTaskSubmissions(ctx context.Context, userId primitive.ObjectID) ([]models.TaskSubmission, error) {
	submissions := []models.TaskSubmission{}
	sR.l.Println(userId)
	cursor, err := sR.taskSubmissionCollection.Find(ctx, bson.M{
		"userid": userId,
	})
	if err != nil {
		sR.l.Println(err)
		return nil, err
	}

	if err = cursor.All(ctx, &submissions); err != nil {
		sR.l.Println(err)
		return nil, err
	}
	return submissions, nil
}

func (sR studentRepo) GetStudentByID(ctx context.Context, studentID primitive.ObjectID) (*models.Student, error) {

	student := new(models.Student)

	res := sR.studentCollection.FindOne(ctx, bson.M{
		"_id": studentID,
	})
	if res.Err() == mongo.ErrNoDocuments {
		sR.l.Println(res.Err())
		return nil, models.ErrNoStudentWithIdExists
	}

	err := res.Decode(&student)
	if err != nil {
		return nil, models.ErrParsingStudent
	}

	return student, nil

}
