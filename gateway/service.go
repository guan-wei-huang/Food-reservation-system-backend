package main

import (
	"context"
	order "reserve_restaurant/order"
	restaurant "reserve_restaurant/restaurant"
	user "reserve_restaurant/user"
	"time"
)

type Service interface {
	NewUser(context.Context, string, string) (int, error)
	UserLogin(context.Context, string, string) (string, string, error)

	GetRestaurantMenu(context.Context, int) (*restaurant.Menu, error)
	CreateFood(context.Context, *restaurant.Food) error
	CreateRestaurant(context.Context, *restaurant.Restaurant) error
	SearchRestaurant(context.Context, string) ([]*restaurant.Restaurant, error)

	CreateOrder(context.Context, *order.Order) (int, error)
	GetOrder(context.Context, int) (*order.Order, error)
	GetOrderForUser(context.Context, int) ([]*order.Order, error)
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

func (s *apiGatewayService) CreateFood(ctx context.Context, f *restaurant.Food) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	if err := s.restaurantClient.CreateFood(ctx, f); err != nil {
		return err
	}
	return nil
}

func (s *apiGatewayService) CreateRestaurant(ctx context.Context, r *restaurant.Restaurant) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	if err := s.restaurantClient.CreateRestaurant(ctx, r); err != nil {
		return err
	}
	return nil
}

func (s *apiGatewayService) SearchRestaurant(ctx context.Context, location string) ([]*restaurant.Restaurant, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	rests, err := s.restaurantClient.SearchRestaurant(ctx, location)
	if err != nil {
		return nil, err
	}
	return rests, nil
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
