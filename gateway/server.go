package main

import (
	"context"
	order "reserve_restaurant/order"
	restaurant "reserve_restaurant/restaurant"
	user "reserve_restaurant/user"
	"time"
)

type Service interface {
}

type apiGatewayService struct {
	orderClient      *order.Client
	userClient       *user.Client
	restaurantClient *restaurant.Client
}

func NewGatewayServer(orderUrl, userUrl, restUrl string) (Service, error) {
	orderClient, err := order.NewClient(orderUrl)
	if err != nil {
		return nil, err
	}

	userClient, err := user.NewClient(userUrl)
	if err != nil {
		return nil, err
	}

	restClient, err := restaurant.NewClient(restUrl)
	if err != nil {
		return nil, err
	}

	return &apiGatewayService{
		orderClient:      orderClient,
		userClient:       userClient,
		restaurantClient: restClient,
	}, nil
}

func (s *apiGatewayService) NewUser(ctx context.Context, name, password string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	uid, err := s.userClient.NewUser(ctx, name, password)
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func (s *apiGatewayService) UserLogin(ctx context.Context, name, password string) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	token, refreshToken, err := s.userClient.UserLogin(ctx, name, password)
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}

func (s *apiGatewayService) GetRestaurantMenu(ctx context.Context, rid int) (*restaurant.Menu, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	m, err := s.restaurantClient.GetRestaurantMenu(ctx, rid)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *apiGatewayService) CreateFood(ctx context.Context, rid int, name, description string, price float32) error {
	f := &restaurant.Food{
		Rid:         rid,
		Name:        name,
		Description: description,
		Price:       price,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	if err := s.restaurantClient.CreateFood(ctx, f); err != nil {
		return err
	}
	return nil
}

func (s *apiGatewayService) CreateRestaurant(ctx context.Context, name, description, location string) error {
	r := &restaurant.Restaurant{
		Name:        name,
		Description: description,
		Location:    location,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	if err := s.restaurantClient.CreateRestaurant(ctx, r); err != nil {
		return err
	}
	return nil
}

func (s *apiGatewayService) CreateOrder(ctx context.Context, order *order.Order) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	id, err := s.orderClient.CreateOrder(ctx, order)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *apiGatewayService) GetOrder(ctx context.Context, id int) (*order.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	order, err := s.orderClient.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *apiGatewayService) GetOrderForUser(ctx context.Context, id int) ([]*order.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	orders, err := s.orderClient.GetOrderForUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
