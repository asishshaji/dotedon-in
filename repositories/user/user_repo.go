package user_repository

import (
	"context"
	"log"

	"github.com/asishshaji/dotedon-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	userCollection   *mongo.Collection
	mentorCollection *mongo.Collection
	l                *log.Logger
}

func NewUserAuthRepo(l *log.Logger, db *mongo.Database) IUserRepository {
	return userRepo{
		userCollection:   db.Collection("users"),
		mentorCollection: db.Collection("mentors"),
		l:                l,
	}
}

func (uR userRepo) RegisterUser(ctx context.Context, uM *models.User) error {
	res, err := uR.userCollection.InsertOne(ctx, uM)

	if err != nil {
		uR.l.Println("Error inserting user")
		return err
	}

	uR.l.Println("Inserted new product with ID : ", res.InsertedID)

	return nil
}

func (uR userRepo) CheckUserExistsWithUserName(ctx context.Context, username string) bool {
	err := uR.userCollection.FindOne(ctx, bson.M{"username": username}).Err()
	if err == mongo.ErrNoDocuments {
		return false
	}
	return true
}

func (uR userRepo) GetUserByUsername(ctx context.Context, username string) *models.User {
	user := new(models.User)
	res := uR.userCollection.FindOne(ctx, bson.M{"username": username})

	if res.Err() == mongo.ErrNoDocuments {
		uR.l.Println("Invalid username and Password")
		return nil
	}

	err := res.Decode(user)
	if err != nil {
		uR.l.Println("Error decoding user")
		return nil
	}

	return user

}

func (uR userRepo) GetMentors(ctx context.Context) ([]*models.MentorDTO, error) {

	mentors := []*models.MentorDTO{}

	cursor, err := uR.mentorCollection.Find(ctx, bson.M{})

	if err != nil {
		uR.l.Fatalln(err)

		return nil, err
	}

	if err = cursor.All(ctx, &mentors); err != nil {
		uR.l.Fatalln(err)
		return nil, err
	}

	return mentors, nil

}
