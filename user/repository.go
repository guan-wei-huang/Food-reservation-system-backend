package user

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	GetUser(ctx context.Context, name string) (*User, error)
	CreateUser(ctx context.Context, name, password string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewUserRepository(dsn string) (Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &repository{db}, nil
}

func (r *repository) GetUser(ctx context.Context, name string) (*User, error) {
	var user = &User{Name: name}
	if err := r.db.WithContext(ctx).Model(&User{}).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) CreateUser(ctx context.Context, name, password string) (*User, error) {
	var user = &User{Name: name, Password: password}
	if err := r.db.WithContext(ctx).Model(&User{}).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
