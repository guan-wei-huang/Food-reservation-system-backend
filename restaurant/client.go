package restaurant

import (
	"context"
	pb "reserve_restaurant/restaurant/pb/restaurant"

	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.RestaurantServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewRestaurantServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) GetRestaurantMenu(ctx context.Context, rid int) (*Menu, error) {
	r, err := c.service.GetRestaurantMenu(ctx, &pb.MenuRequest{Rid: int32(rid)})
	if err != nil {

	}
}
