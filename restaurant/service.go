package restaurant

import (
	"context"
)

type Restaurant struct {
	ID          int
	Name        string
	Description string
	location    string
}

type Menu struct {
	Rid   int
	Foods []Food
}

type Food struct {
	Fid         int
	Rid         int
	Name        string
	Description string
	Price       float32
}

type Service interface {
	CreateRestaurant(ctx context.Context, r *Restaurant) error
	CreateFood(ctx context.Context, f *Food) error
	GetRestaurantMenu(ctx context.Context, rid int) (*Menu, error)
}

type restaurantService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &restaurantService{r}
}

func (s *restaurantService) GetRestaurantMenu(ctx context.Context, rid int) (*Menu, error) {
	menu, err := s.repo.GetMenu(ctx, rid)
	if err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *restaurantService) CreateRestaurant(ctx context.Context, r *Restaurant) error {
	if err := s.repo.CreateRestaurant(ctx, r); err != nil {
		return err
	}
	return nil
}

func (s *restaurantService) CreateFood(ctx context.Context, f *Food) error {
	if err := s.repo.CreateFood(ctx, f); err != nil {
		return err
	}
	return nil
}
