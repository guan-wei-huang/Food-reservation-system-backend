package order

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidAccess = errors.New("this page is invalid")
)

type Product struct {
	Fid      int
	Name     string
	Price    float32
	Quantity int
}

type Order struct {
	Id        int
	Rid       int
	Uid       int
	Products  *[]Product
	CreatedAt time.Time
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
		return 0, err
	}
	return oid, nil
}

func (s *orderService) GetOrder(ctx context.Context, oid, uid int) (*Order, error) {
	order, err := s.repo.GetOrder(ctx, oid)
	if err != nil {
		return nil, err
	}

	if order.Uid != uid {
		return nil, ErrInvalidAccess
	}
	return order, nil
}

func (s *orderService) GetOrderForUser(ctx context.Context, id int) (*[]Order, error) {
	orders, err := s.repo.GetOrderForUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
