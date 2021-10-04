package order

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidAccess = errors.New("this page is invalid")
)

type Product struct {
	Fid      int     `json:"fid" form:"fid"`
	Name     string  `json:"name" form:"name"`
	Price    float32 `json:"price" form:"price"`
	Quantity int     `json:"quantity" form:"quantity"`
}

type Order struct {
	Id        int       `json:"id" form:"id"`
	Rid       int       `json:"rid" form:"rid"`
	Uid       int       `json:"uid" form:"uid"`
	Products  []Product `json:"products" form:"products"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
}

type Service interface {
	CreateOrder(ctx context.Context, order *Order) (int, error)
	GetOrder(ctx context.Context, oid, uid int) (*Order, error)
	GetOrderForUser(ctx context.Context, id int) (*[]Order, error)
}

type orderService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &orderService{repo: r}
}

func (s *orderService) CreateOrder(ctx context.Context, order *Order) (int, error) {
	oid, err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return 0, fmt.Errorf("create order failed: %v", err)
	}
	return oid, nil
}

func (s *orderService) GetOrder(ctx context.Context, oid, uid int) (*Order, error) {
	order, err := s.repo.GetOrder(ctx, oid)
	if errors.Is(err, ErrOrderInvalid) {
		return nil, ErrOrderInvalid
	} else if err != nil {
		return nil, fmt.Errorf("get order failed: %v", err)
	}

	if order.Uid != uid {
		return nil, ErrInvalidAccess
	}
	return order, nil
}

func (s *orderService) GetOrderForUser(ctx context.Context, id int) (*[]Order, error) {
	orders, err := s.repo.GetOrderForUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get orders for user failed: %v", err)
	}
	return orders, nil
}
