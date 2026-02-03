package service

import (
	"errors"
	"time"

	"go-login-api-task/models"
	"go-login-api-task/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("secret_key")

type AuthService struct {
	userRepo *repository.UserRepository
}

// Constructor
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

// Login logic
func (s *AuthService) Login(email, password string) (string, error) {
	// Get user from DB
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT
	token, err := generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// JWT creation
func generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
