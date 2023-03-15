package restaurant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Restaurant struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type Menu struct {
	Rid   int    `json:"rid"`
	Foods []Food `json:"foods"`
}

type Food struct {
	Fid         int     `json:"fid"`
	Rid         int     `json:"rid"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

type Service interface {
	CreateRestaurant(ctx context.Context, r *Restaurant) (int, error)
	CreateFood(ctx context.Context, f *Food) error
	GetRestaurantMenu(ctx context.Context, rid int) (*Menu, error)
	SearchRestaurant(ctx context.Context, location string) ([]*Restaurant, error)
}

type restaurantService struct {
	repo Repository
}

var (
	ErrRestaurantIdWrong = errors.New("restaurant isn't exist")
)

func NewService(r Repository) Service {
	return &restaurantService{r}
}

func (s *restaurantService) GetRestaurantMenu(ctx context.Context, rid int) (*Menu, error) {
	exist, err := s.repo.CheckRestaurantExist(ctx, rid)
	if err != nil {
		return nil, fmt.Errorf("check restaurant exist err: %v", err)
	} else if !exist {
		return nil, ErrRestaurantIdWrong
	}

	menu, err := s.repo.GetMenu(ctx, rid)
	if err != nil {
		return nil, fmt.Errorf("get menu err: %v", err)
	}
	return menu, nil
}

func (s *restaurantService) CreateRestaurant(ctx context.Context, r *Restaurant) (int, error) {
	rid, err := s.repo.CreateRestaurant(ctx, r)
	if err != nil {
		return 0, fmt.Errorf("create restaurant failed: %v", err)
	}
	return rid, nil
}

func (s *restaurantService) CreateFood(ctx context.Context, f *Food) error {
	if err := s.repo.CreateFood(ctx, f); err != nil {
		return fmt.Errorf("create food failed: %v", err)
	}
	return nil
}

func (s *restaurantService) SearchRestaurant(ctx context.Context, location string) ([]*Restaurant, error) {
	// use api to translate location to latitude and longitude
	mapUrl := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%v&key=%v",
		location, GetConfig().ApiKey)
	resp, err := http.Get(mapUrl)
	if err != nil {
		return nil, fmt.Errorf("google api visit err: %v", err)
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	detail := &mapResponse{}
	if err := json.Unmarshal(response, detail); err != nil {
		return nil, err
	}

	latitude, longitude := detail.Results[0].Geometry.Location.Lat, detail.Results[0].Geometry.Location.Lng
	restaurants, err := s.repo.SearchRestaurant(ctx, latitude, longitude)
	if err != nil {
		return nil, fmt.Errorf("search restaurant failed: %v", err)
	}
	return restaurants, nil
}
