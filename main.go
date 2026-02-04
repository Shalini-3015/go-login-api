package main

import (

"go-login-api-task/config"
"go-login-api-task/models"
"go-login-api-task/router"

)

func main() {
	config.ConnectDatabase()
	migrateModels()
	r := router.SetupRouter()
	r.Run(":8080")
}

func migrateModels() {
	config.DB.AutoMigrate(
		&models.LoginUser{},
		
	)
}
