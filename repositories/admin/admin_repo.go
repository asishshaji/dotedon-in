package admin_repository

import (
	"context"
	"errors"
	"log"

	"github.com/asishshaji/dotedon-api/models"
	"github.com/asishshaji/dotedon-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminRepository struct {
	l                 *log.Logger
	adminCollection   *mongo.Collection
	taskCollection    *mongo.Collection
	typeCollection    *mongo.Collection
	studentCollection *mongo.Collection
}

func NewAdminRepository(l *log.Logger, db *mongo.Database) IAdminRepository {
	return AdminRepository{
		l:                 l,
		adminCollection:   db.Collection("admin"),
		taskCollection:    db.Collection("tasks"),
		typeCollection:    db.Collection("types"),
		studentCollection: db.Collection("students"),
	}
}

func (adminRepo AdminRepository) GetAdmin(ctx context.Context, username string) (*models.Admin, error) {
	admin := new(models.Admin)

	res := adminRepo.adminCollection.FindOne(ctx, bson.M{"username": username})

	if res.Err() == mongo.ErrNoDocuments {
		adminRepo.l.Println("No admin with username", username, "exists")
		return nil, res.Err()
	}

	err := res.Decode(admin)

	if err != nil {
		adminRepo.l.Println(err)
		return nil, err
	}

	return admin, nil
}

func (aR AdminRepository) AddTask(ctx context.Context, task models.Task) error {

	res, err := aR.taskCollection.InsertOne(ctx, task)

	if err != nil {
		aR.l.Println(err)
		return err
	}

	aR.l.Println("Inserted task with ID", res.InsertedID)

	return nil
}

func (aR AdminRepository) UpdateTask(ctx context.Context, task models.TaskUpdate) error {
	opts := options.Update().SetUpsert(true)

	up, err := utils.ToDoc(task)

	aR.l.Println(up)

	if err != nil {
		return err
	}

	doc := bson.M{"$set": up}

	res, err := aR.taskCollection.UpdateByID(ctx, task.Id, doc, opts)
	if err != nil {
		aR.l.Println(err)
		return err
	}
	aR.l.Println(res.UpsertedID)
	return nil
}

func (aR AdminRepository) GetTasks(ctx context.Context) ([]models.Task, error) {
	tasks := []models.Task{}

	cursor, err := aR.taskCollection.Find(ctx, bson.M{})
	if err != nil {
		aR.l.Println(err)
		return nil, err
	}

	if err = cursor.All(ctx, &tasks); err != nil {
		aR.l.Println(err)
		return nil, err
	}

	return tasks, nil
}

func (aR AdminRepository) CreateType(ctx context.Context, typeT models.Type) error {

	return nil
}

func (aR AdminRepository) GetUsers(ctx context.Context) (models.Students, error) {
	students := new([]models.Student)
	cursor, err := aR.studentCollection.Find(ctx, bson.M{})
	if err != nil {
		aR.l.Println(err)
		return nil, err
	}
	if err = cursor.All(ctx, students); err != nil {
		aR.l.Println(err)
		return nil, err
	}

	return *students, nil
}
func (aR AdminRepository) DeleteTask(ctx context.Context, taskId primitive.ObjectID) error {
	res, err := aR.taskCollection.DeleteOne(ctx, bson.M{
		"_id": taskId,
	})
	if res.DeletedCount == 0 {
		aR.l.Println("No task found to be deleted")
		return errors.New("no task found with given id")

	}
	if err != nil {
		aR.l.Println(err)
		return err
	}
	return nil
}
