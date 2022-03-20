package main

import (
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	r "stablex/repository/mongodb"
	"strconv"
	"time"
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
		//bson.M{
		//	"first_name": "Jhon",
		//	"last_name":  "Doe",
		//	"position":   "technician",
		//	"actions": []interface{}{
		//		bson.D{
		//			{"type", "setup"}, {"created_at", time.Now().UTC()},
		//		},
		//		bson.D{
		//			{"type", "start"}, {"created_at", time.Now().UTC()},
		//		},
		//	},
		//},
		bson.M{
			"first_name": "Mike",
			"last_name":  "Foo",
			"position":   "technician",
			"actions": []interface{}{
				bson.D{
					{"type", "setup"}, {"created_at", time.Now().UTC()},
				},
				bson.D{
					{"type", "start"}, {"created_at", time.Now().AddDate(0, 1, 0).UTC()},
				},
			},
		},
	}

	repo.InsertOperators(docs)
}
