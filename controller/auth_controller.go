package controller

import (
	"net/http"

	"go-login-api-task/service"

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
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}


	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	token, err := a.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
