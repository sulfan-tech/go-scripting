package main

import (
	"context"
	"go-scripting/scripts/create_dummy_user_staging/service"
	"log"

	"github.com/joho/godotenv"
)

var dummyUserService service.UserDummyImpl

func init() {
	env := "../../.env"
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file:", err.Error())
	}

	dummyUserService = service.NewInstanceDummyUserService()
}

func main() {
	err := dummyUserService.CreateUserDummyFS(context.Background(), 1000)
	if err != nil {
		log.Printf("Error creating user : %v", err)
	}
}
