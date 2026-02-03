package main

import (
	"fmt"
	"log"

	"go-login-api-task/config"
	"go-login-api-task/models"
	"go-login-api-task/service"
)

func main() {

	config.ConnectDatabase()

	

	
	hashedPassword, err := service.HashPassword("password123")
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	

	if err := config.DB.Create(&user).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("User created with hashed password")
}
