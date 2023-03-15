package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Service interface {
	UserRegister(ctx context.Context, name, password string) (*User, error)
	UserLogin(ctx context.Context, name, password string) (string, string, error)
}

type userService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &userService{repo: r}
}

func (s *userService) UserRegister(ctx context.Context, name, password string) (*User, error) {
	exists, err := s.repo.CheckUserExist(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("check user exist failed: %v", err)
	} else if exists {
		return nil, NewUserError("this name has been registered")
	}

	hashPassword, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash user:%v's password failed: %v", name, err)
	}

	user, err := s.repo.CreateUser(ctx, name, hashPassword)
	if err != nil {
		return nil, fmt.Errorf("create user:%v failed: %v", name, err)
	}
	return user, nil
}

func (s *userService) UserLogin(ctx context.Context, name, password string) (string, string, error) {
	exist, err := s.repo.CheckUserExist(ctx, name)
	if err != nil {
		return "", "", fmt.Errorf("check user exist failed: %v", err)
	} else if !exist {
		return "", "", NewUserError("user doesn't exist")
	}

	user, err := s.repo.GetUser(ctx, name)
	if err != nil {
		return "", "", fmt.Errorf("fetch user:%v failed: %v", name, err)
	}

	if ok := comparePassword(user.Password, password); !ok {
		return "", "", NewUserError("password is wrong")
	}

	var token, refreshToken string
	if token, err = generateToken(user); err != nil {
		return "", "", fmt.Errorf("jwt format failed: %v", err)
	}

	return token, refreshToken, nil
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

func generateToken(u *User) (string, error) {
	claim := &jwt.MapClaims{
		"exp":      time.Now().Add(time.Hour).Unix(),
		"userId":   strconv.Itoa(u.ID),
		"userName": u.Name,
	}

	token, err := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claim).SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return token, nil
}
