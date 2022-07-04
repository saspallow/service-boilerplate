package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"service-boilerplate/authentication"
	"service-boilerplate/server"
)

type Config struct {
	Port string
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		log.Println("We are getting the env values")
	}

	config := Config{
		Port: os.Getenv("PORT"),
	}

	timeoutContext := time.Duration(500) * time.Second

	// Repository
	var authenticationRepo authentication.Repository
	authenticationRepo = authentication.NewRepository(timeoutContext)

	// Services
	var authenticationSvc authentication.Service
	authenticationSvc = authentication.NewService(
		authenticationRepo,
		timeoutContext,
	)

	srv := server.New(
		authenticationSvc,
	)

	errs := make(chan error, 2)
	go func() {
		log.Println("transport", "http", "address", ":"+config.Port, "msg", "listening")
		errs <- http.ListenAndServe(":"+config.Port, srv)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Println("terminated", <-errs)
}
