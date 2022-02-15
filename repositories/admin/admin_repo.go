package admin_repository

import (
	"context"
	"log"

	"github.com/asishshaji/dotedon-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepository struct {
	l               *log.Logger
	adminCollection *mongo.Collection
}

func NewAdminRepository(l *log.Logger, admiCollection *mongo.Collection) IAdminRepository {
	return AdminRepository{
		l:               l,
		adminCollection: admiCollection,
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
