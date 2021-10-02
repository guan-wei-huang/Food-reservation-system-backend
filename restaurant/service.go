package restaurant

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Restaurant struct {
	ID          int     `json:"id" form:"id"`
	Name        string  `json:"name" form:"name"`
	Description string  `json:"description" form:"description"`
	Location    string  `json:"location" form:"location"`
	Latitude    float64 `json:"latitude" form:"latitude"`
	Longitude   float64 `json:"longitude" form:"longtitude"`
}

type Menu struct {
	Rid   int `json:"rid"`
	Foods []Food
}

type Food struct {
	Fid         int     `json:"fid" form:"fid"`
	Rid         int     `json:"rid" form:"rid"`
	Name        string  `json:"name" form:"name"`
	Description string  `json:"description" form:"description"`
	Price       float32 `json:"price" form:"price"`
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

func (s *restaurantService) CreateRestaurant(ctx context.Context, r *Restaurant) (int, error) {
	//  TODO: translate location to latitude and longtitude
	rid, err := s.repo.CreateRestaurant(ctx, r)
	if err != nil {
		return 0, err
	}
	return rid, nil
}

func (s *restaurantService) CreateFood(ctx context.Context, f *Food) error {
	if err := s.repo.CreateFood(ctx, f); err != nil {
		return err
	}
	return nil
}

func (s *restaurantService) SearchRestaurant(ctx context.Context, location string) ([]*Restaurant, error) {
	// use api to translate location to latitude and longitude
	mapUrl := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%v&key=%v", location, GOOGLE_API_KEY)
	resp, err := http.Get(mapUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
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
		return nil, err
	}
	return restaurants, nil
}
