package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fedo3nik/GamePortal_IdentityService/config"
	"github.com/fedo3nik/GamePortal_IdentityService/internal/application/service"
	"github.com/fedo3nik/GamePortal_IdentityService/internal/interface/controller"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const contextTimeout = 10

func initClient(ctx context.Context, connURI string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(connURI))
	if err != nil {
		log.Panicf("initClient error: %v", err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Panicf("Connect error: %v", err)
	}

	return client
}

func main() {
	c, err := config.NewConfig()
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout*time.Second)

	fmt.Println(c.ConnURI)

	mongoClient := initClient(ctx, c.ConnURI)
	userService := service.NewUserService(mongoClient, c.DB)

	signUpHandler := controller.NewHTTPSignUpHandler(userService)
	signInHandler := controller.NewHTTPSignInHandler(userService)
	addWarningHandler := controller.NewHTTPAddWarnHandler(userService)
	remWarningHandler := controller.NewHTTPRemHandler(userService)
	isBannedHandler := controller.NewHTTPIsBannedHandler(userService)

	handler := mux.NewRouter()

	handler.Handle("/user/sign-up", signUpHandler).Methods("POST")
	handler.Handle("/user/sign-in", signInHandler).Methods("POST")
	handler.Handle("/user/warn/{id}", addWarningHandler).Methods("PUT")
	handler.Handle("/user/remove-warning/{id}", remWarningHandler).Methods("PUT")
	handler.Handle("/user/is-banned/{id}", isBannedHandler).Methods("GET")

	err = http.ListenAndServe(c.Host+c.Port, handler)
	if err != nil {
		log.Panicf("Listen %v error: %v", c.Host+c.Port, err)
	}

	defer func() {
		cancel()

		err = mongoClient.Disconnect(ctx)
		if err != nil {
			log.Printf("Disconnect error: %v", err)
		}
	}()
}
