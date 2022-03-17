package main

import (
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	r "stablex/repository/mongodb"
	"strconv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mongoURL := os.Getenv("MONGO_URL")
	mongodb := os.Getenv("MONGO_DB")
	mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))

	repo, err := r.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
	if err != nil {
		log.Fatal(err)
	}

	docs := []interface{}{
		bson.M{
			"first_name": "Jhon",
			"last_name":  "Doe",
			"position":   "technician",
			"actions":    []interface{}{
				//bson.D{
				//	{"type", "setup"}, {"created_at", "1647455473"},
				//},
				//bson.D{
				//	{"type", "start"}, {"created_at", "1647455473"},
				//},
			},
		},
	}

	repo.InsertOperators(docs)
}
