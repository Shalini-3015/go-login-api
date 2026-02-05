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

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *AuthService) UserLogin(email, password string) (string, error) {
	
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email")
	}

	if err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return "", errors.New("invalid password")
	}

	token, err := generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}


func (s *AuthService) RegisterUser(email, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := &models.LoginUser{
		Email:    email,
		Password: hashedPassword,
	}

	return s.userRepo.CreateUser(user)
}




func generateJWT(user *models.LoginUser) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
