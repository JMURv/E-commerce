package config

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB
var JWTsecretKey string
var ReviewServiceURL string
var UserSVC string

func Connect() {
	var err error

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	log.Println("Env file has been loaded")

	jwt := os.Getenv("JWTSECRET")
	dsn := os.Getenv("DSN")
	ReviewServiceURL = os.Getenv("REVIEW_SVC")
	UserSVC = os.Getenv("USER_SVC")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database")

	db = database
	JWTsecretKey = jwt
}

func GetDB() *gorm.DB {
	return db
}
