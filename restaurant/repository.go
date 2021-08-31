package restaurant

import (
	"context"
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrRestaurantIdWrong = errors.New("restaurant not found")
)

type Repository interface {
	CreateRestaurant(ctx context.Context, rest *Restaurant) error
	CreateFood(ctx context.Context, f *Food) error
	GetMenu(ctx context.Context, rid int) (*Menu, error)
}

type repository struct {
	db *gorm.DB
}

func NewRestaurantRepository(dsn string) (Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &repository{db}, nil
}

func (r *repository) CreateRestaurant(ctx context.Context, rest *Restaurant) error {
	if err := r.db.WithContext(ctx).Model(&Restaurant{}).Create(*rest).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) CreateFood(ctx context.Context, f *Food) error {
	if err := r.db.WithContext(ctx).Model(&Food{}).Create(*f).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) GetMenu(ctx context.Context, rid int) (*Menu, error) {
	var foods = &[]Food{}
	if err := r.db.WithContext(ctx).Model(&Food{}).Where("rid = ?", rid).Find(foods).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRestaurantIdWrong
		} else {
			return nil, err
		}
	}
	return &Menu{rid, *foods}, nil
}
