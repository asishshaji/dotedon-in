package student_repository

import (
	"context"
	"log"

	"github.com/asishshaji/dotedon-api/models"
	"github.com/asishshaji/dotedon-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type studentRepo struct {
	studentCollection        *mongo.Collection
	mentorCollection         *mongo.Collection
	taskSubmissionCollection *mongo.Collection
	taskCollection           *mongo.Collection
	domainCollection         *mongo.Collection
	collegeCollection        *mongo.Collection
	coursesCollection        *mongo.Collection
	tokenCollection          *mongo.Collection
	notificationCollection   *mongo.Collection
	l                        *log.Logger
}

func NewStudentAuthRepo(l *log.Logger, db *mongo.Database) IStudentRepository {
	return studentRepo{
		studentCollection:        db.Collection("students"),
		mentorCollection:         db.Collection("mentor"),
		taskSubmissionCollection: db.Collection("task_submission"),
		domainCollection:         db.Collection("domains"),
		taskCollection:           db.Collection("tasks"),
		coursesCollection:        db.Collection("courses"),
		collegeCollection:        db.Collection("colleges"),
		tokenCollection:          db.Collection("tokens"),
		notificationCollection:   db.Collection("notifications"),
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

func (uR studentRepo) CheckStudentExistsWithEmail(ctx context.Context, email string) bool {
	res := uR.studentCollection.FindOne(ctx, bson.M{
		"email": email,
	}).Err()
	return res != mongo.ErrNoDocuments
}

func (uR studentRepo) GetStudentByEmail(ctx context.Context, email string) *models.Student {
	student := new(models.Student)
	res := uR.studentCollection.FindOne(ctx, bson.M{"email": email})

	if res.Err() == mongo.ErrNoDocuments {
		uR.l.Println("Invalid studentname and password")
		return nil
	}

	err := res.Decode(student)
	if err != nil {
		uR.l.Println("Error decoding student")
		return nil
	}

	return student

}

func (uR studentRepo) GetMentorsNotInIDS(ctx context.Context, ids []primitive.ObjectID) ([]*models.Mentor, error) {

	mentors := []*models.Mentor{}

	cursor, err := uR.mentorCollection.Find(ctx, bson.M{
		"_id": bson.M{"$nin": ids},
	})

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

func (uR studentRepo) AddDomainToStudent(ctx context.Context, userID primitive.ObjectID, domain string) error {
	options := bson.M{
		"$addToSet": bson.M{
			"domains": domain,
		},
	}
	_, err := uR.studentCollection.UpdateByID(ctx, userID, options)
	if err != nil {
		uR.l.Println(err)
		return err
	}
	return nil
}

func (uR studentRepo) GetMentorByID(ctx context.Context, mentorId primitive.ObjectID) (models.Mentor, error) {
	mentor := models.Mentor{}

	res := uR.mentorCollection.FindOne(ctx, bson.M{
		"_id": mentorId,
	})

	if err := res.Decode(&mentor); err != nil {
		uR.l.Println(err)
		return mentor, err
	}

	return mentor, nil
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

	if res.ModifiedCount == 0 {
		uR.l.Println("Already following mentor")
		return models.ErrAlreadyFollowing
	}

	return nil
}

func (sR studentRepo) UpdateTaskSubmission(ctx context.Context, task models.TaskSubmission) error {

	opts := options.Update().SetUpsert(true)
	up, err := utils.ToDoc(task)

	if err != nil {
		return err
	}
	doc := bson.M{"$set": up}

	res, err := sR.taskSubmissionCollection.UpdateOne(ctx, bson.M{
		"userid": task.UserId,
		"taskid": task.TaskId,
	}, doc, opts)

	if err != nil {
		sR.l.Println(err)
		return err
	}

	sR.l.Println("Inserted submission with ID", res.UpsertedID)

	return nil

}

func (sR studentRepo) UpdateStudent(ctx context.Context, student models.Student) error {
	opts := options.Update().SetUpsert(true)
	up, err := utils.ToDoc(student)

	if err != nil {
		return err
	}
	doc := bson.M{"$set": up}

	res, err := sR.studentCollection.UpdateOne(ctx, bson.M{
		"_id": student.ID,
	}, doc, opts)

	if err != nil {
		sR.l.Println(err)
		return err
	}

	sR.l.Println("Updated student with ID", res.UpsertedID)

	return nil

}

func checkIfTaskSubmissionExists(ctx context.Context, collection *mongo.Collection, taskId, userId primitive.ObjectID) bool {
	res := collection.FindOne(ctx, bson.M{
		"taskid": taskId,
		"userid": userId,
	})

	return res.Err() != mongo.ErrNoDocuments

}

func (sR studentRepo) CreateTaskSubmission(ctx context.Context, task models.TaskSubmission) error {

	if checkIfTaskSubmissionExists(ctx, sR.taskSubmissionCollection, task.TaskId, task.UserId) {
		sR.l.Println(models.ErrTaskSubmissionExists)
		return models.ErrTaskSubmissionExists
	}

	res, err := sR.taskSubmissionCollection.InsertOne(ctx, task)
	if err != nil {
		sR.l.Println(err)
		return err
	}
	sR.l.Println("Inserted task submission with id", res.InsertedID)

	return nil

}
func (sR studentRepo) GetTasks(ctx context.Context, domains, sems []string) ([]models.Task, error) {
	tasks := []models.Task{}

	cursor, err := sR.taskCollection.Find(ctx, bson.M{
		"$and": bson.A{
			bson.D{
				{
					"type", bson.D{{
						"$in", domains,
					}},
				},
			},
			bson.D{
				{
					"semester", bson.D{{
						"$in", sems,
					}},
				},
			},
		},
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

func (sR studentRepo) GetMentorIDsFollowedByStudent(ctx context.Context, userid primitive.ObjectID) ([]primitive.ObjectID, error) {

	cursor, err := sR.studentCollection.Aggregate(ctx, mongo.Pipeline{

		bson.D{{
			"$match", bson.D{{
				"_id", userid,
			}},
		}},

		bson.D{{
			"$project", bson.D{{
				"_id", 0,
			}, {
				"mentors", 1,
			}},
		}},
	})
	if err != nil {
		sR.l.Println(err)
		return nil, err
	}

	mentors := []models.TT{}

	if err = cursor.All(ctx, &mentors); err != nil {
		sR.l.Println(err)
		return nil, err
	}
	if len(mentors) == 0 {
		return []primitive.ObjectID{}, nil
	}

	men := mentors[0].Mentors
	if men != nil {
		return men, nil
	}

	return []primitive.ObjectID{}, nil
}

func (sR studentRepo) GetDomains(ctx context.Context) ([]models.StaticModel, error) {
	domains := []models.StaticModel{}

	cursor, err := sR.domainCollection.Find(ctx, bson.M{})
	if err != nil {
		return domains, err
	}
	if err = cursor.All(ctx, &domains); err != nil {
		return domains, err
	}

	return domains, nil
}

func (sR studentRepo) GetColleges(ctx context.Context) ([]models.StaticModel, error) {
	c := []models.StaticModel{}

	cursor, err := sR.collegeCollection.Find(ctx, bson.M{})
	if err != nil {
		return c, err
	}
	if err = cursor.All(ctx, &c); err != nil {
		return c, err
	}

	return c, nil
}

func (sR studentRepo) GetNotifications(ctx context.Context, uid primitive.ObjectID) ([]models.NotificationEntity, error) {
	notifications := []models.NotificationEntity{}

	cursor, err := sR.notificationCollection.Find(ctx, bson.M{
		"user_id": uid,
	})

	if err != nil {
		return notifications, err
	}

	if err = cursor.All(ctx, &notifications); err != nil {
		return notifications, err
	}

	return notifications, nil

}

func (sR studentRepo) GetCourses(ctx context.Context) ([]models.StaticModel, error) {
	c := []models.StaticModel{}

	cursor, err := sR.coursesCollection.Find(ctx, bson.M{})
	if err != nil {
		return c, err
	}
	if err = cursor.All(ctx, &c); err != nil {
		return c, err
	}

	return c, nil
}

func (sR studentRepo) InsertToken(ctx context.Context, tK models.Token) error {
	opts := options.Update().SetUpsert(true)

	_, err := sR.tokenCollection.UpdateOne(ctx, bson.M{
		"user_id": tK.UserId,
	}, bson.M{
		"$set": tK,
	}, opts)

	if err != nil {
		return err
	}
	return nil
}
