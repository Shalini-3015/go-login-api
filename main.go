package main

import (
	"go-login-api-task/config"
	"go-login-api-task/models"
	"go-login-api-task/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.ConnectDatabase()
	migrateModels()
	router.SetupRouter(r)
	r.Run(":8080")
}

func migrateModels() {      
	config.DB.AutoMigrate(
		&models.LoginUser{},
		&models.Currency{},
		&models.ExchangeRate{},
	)
}
