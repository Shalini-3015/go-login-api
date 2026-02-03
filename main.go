package main

import (

"go-login-api-task/config"
"go-login-api-task/models"
"go-login-api-task/router"

)

func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{})
	r := router.SetupRouter()
	r.Run(":8080")
}