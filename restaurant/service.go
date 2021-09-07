package restaurant

import (
	"context"
)

type Restaurant struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

type Menu struct {
	Rid   int `json:"rid"`
	Foods []Food
}

type Food struct {
	Fid         int     `json:"fid"`
	Rid         int     `json:"rid"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

type Service interface {
	CreateRestaurant(ctx context.Context, r *Restaurant) error
	CreateFood(ctx context.Context, f *Food) error
	GetRestaurantMenu(ctx context.Context, rid int) (*Menu, error)
	SearchRestaurant(ctx context.Context, location string) ([]*Restaurant, error)
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

func (s *restaurantService) SearchRestaurant(ctx context.Context, location string) ([]*Restaurant, error) {

}
