package router

import (
	"go-login-api-task/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	authController := controller.NewAuthController()

	
	r.POST("/login", authController.Login)

	return r
}
