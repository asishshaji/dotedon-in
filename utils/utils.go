package utils

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EnvironmentConfig struct {
	ServerPort string
	DBURL      string
	DBName     string
	l          *log.Logger
}

func LoadEnv(l *log.Logger) *EnvironmentConfig {
	if err := godotenv.Load(); err != nil {
		l.Fatalln("Error loading env file")
	}

	return &EnvironmentConfig{
		ServerPort: os.Getenv("SERVER_PORT"),
		DBURL:      os.Getenv("DB_URL"),
		DBName:     os.Getenv("DB_NAME"),
		l:          l,
	}
}

func (env *EnvironmentConfig) ConnectToDB() *mongo.Database {
	env.l.Println("Starting connection to db")

	client, err := mongo.NewClient(options.Client().ApplyURI(env.DBURL))

	if err != nil {
		env.l.Fatalln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		env.l.Fatalln(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		env.l.Fatalln(err)
	}

	env.l.Println("Connected to db")

	return client.Database(env.DBName)

}
