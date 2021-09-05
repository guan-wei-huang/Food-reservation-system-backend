package order

import (
	"context"
	"errors"
	pb "reserve_restaurant/order/pb/order"

	"google.golang.org/grpc"
)

var (
	ErrInternalServer = errors.New("internal server error")
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := pb.NewOrderServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) CreateOrder(ctx context.Context, order *Order) (int, error) {
	r, err := c.service.CreateOrder(ctx, &pb.CreateOrderRequest{Order: formatOrder(order)})
	if err != nil {
		return 0, ErrInternalServer
	}

	if r.Error != "" {
		return 0, errors.New(r.Error)
	}
	return int(r.Id), nil
}

func (c *Client) GetOrder(ctx context.Context, id int) (*Order, error) {
	r, err := c.service.GetOrder(ctx, &pb.GetOrderRequest{Id: int32(id)})
	if err != nil {
		return nil, ErrInternalServer
	}

	if r.Error != "" {
		return nil, errors.New(r.Error)
	}

	return parseOrder(r.Order), nil
}

func (c *Client) GetOrderForUser(ctx context.Context, id int) ([]*Order, error) {
	r, err := c.service.GetOrderForUser(ctx, &pb.GetOrderRequest{Uid: int32(id)})
	if err != nil {
		return nil, ErrInternalServer
	}

	if r.Error != "" {
		return nil, errors.New(r.Error)
	}

	orders := []*Order{}
	for _, o := range r.Orders {
		orders = append(orders, parseOrder(o))
	}
	return orders, nil
}
