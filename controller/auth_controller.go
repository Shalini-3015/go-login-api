package controller

import (
	"net/http"

	"go-login-api-task/service"
    "go-login-api-task/dto/auth"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}


func NewAuthController() *AuthController {
	return &AuthController{
		authService: service.NewAuthService(),
	}
}


func (a *AuthController) Login(c *gin.Context) {
	var req auth.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	token, err := a.authService.UserLogin(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, auth.LoginResponse{
		AccessToken: token,
	})
}
