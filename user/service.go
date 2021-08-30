package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Service interface {
	UserRegister(ctx context.Context, name, password string) (*User, error)
	UserLogin(ctx context.Context, name, password string) (*User, error)
}

type userService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &userService{repo: r}
}

func (s *userService) UserRegister(ctx context.Context, name, password string) (*User, error) {
	_, err := s.repo.GetUser(ctx, name)
	if err == nil {
		return nil, NewUserError("this name has been registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("fetch user:%v failed", name)
	}

	newPassword, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash user:%v's password failed", name)
	}

	user, err := s.repo.CreateUser(ctx, name, newPassword)
	if err != nil {
		return nil, fmt.Errorf("create user:%v failed", name)
	}
	return user, nil
}

func (s *userService) UserLogin(ctx context.Context, name, password string) (*User, error) {
	user, err := s.repo.GetUser(ctx, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewUserError("account doesn't exist")
		} else {
			return nil, fmt.Errorf("fetch user:%v failed", name)
		}
	}

	if ok := comparePassword(password, user.Password); !ok {
		return nil, NewUserError("password is wrong")
	}
	return user, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePassword(p1, p2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p1), []byte(p2))
	return err == nil
}