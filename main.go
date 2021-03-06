package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"stablex/api"
	"stablex/domain"
	mr "stablex/repository/mongodb"
	"stablex/repository/redis"
	"stablex/router"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// load env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// setup repo & redis client
	repo := getRepo()
	rdb, err := redis.NewRedisRepository(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// inject dependencies
	service := domain.NewOperatorService(repo)
	handler := api.NewHandler(service, rdb)
	appRouter := router.New(handler)

	srv := &http.Server{
		Addr:         httpPort(),
		Handler:      appRouter,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000")
		errs <- srv.ListenAndServe()

	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func getRepo() domain.OperatorRepository {
	mongoURL := os.Getenv("MONGO_URL")
	mongodb := os.Getenv("MONGO_DB")
	mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
	repo, err := mr.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
