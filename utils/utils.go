package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/asishshaji/dotedon-api/models"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func AdminAuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		admin := c.Get("user").(*jwt.Token)
		claims := admin.Claims.(*models.AdminJWTClaims)
		fmt.Println(claims)
		c.Logger().Print(claims.AdminId)
		if !claims.IsAdmin {
			return echo.ErrForbidden
		}

		c.Set("admin_id", claims.AdminId)
		return next(c)
	}
}

func StudentAuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.StudentJWTClaims)
		c.Logger().Print(claims.StudentId)

		c.Set("student_id", claims.StudentId)

		return next(c)
	}
}
