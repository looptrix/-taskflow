package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"taskflow/internal/model"
	"taskflow/internal/repository"
)

type AuthService struct {
	users  *repository.UserRepository
	secret string
}

func NewAuthService(users *repository.UserRepository, secret string) *AuthService {
	return &AuthService{users: users, secret: secret}
}

func (s *AuthService) Register(name, email, password string) (*model.User, error) {
	if name == "" || email == "" || len(password) < 8 {
		return nil, errors.New("укажите имя, email и пароль длиной минимум 8 символов")
	}

	taken, err := s.users.IsEmailTaken(email)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, errors.New("пользователь с таким email уже существует")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hash),
	}
	if err := s.users.Create(user); err != nil {
		return nil, err
	}
	return user, nil

}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.users.FindByEmail(email)
	if err != nil {
		return "", errors.New("неверный email или пароль")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("неверный email или пароль")
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))

}
