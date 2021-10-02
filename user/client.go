package user

import (
	"context"
	"errors"
	"log"
	pb "reserve_restaurant/user/pb/user"

	"google.golang.org/grpc"
)

var (
	ErrInternalServer = errors.New("internal server error")
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.UserServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := pb.NewUserServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) NewUser(ctx context.Context, name, password string) (int, error) {
	log.Println("new user: ", name, " ", password)
	r, err := c.service.NewUser(ctx, &pb.NewUserRequest{Name: name, Password: password})
	if err != nil {
		return 0, ErrInternalServer
	}

	if r.Err != "" {
		return 0, errors.New(r.Err)
	}
	return int(r.Id), nil
}

func (c *Client) UserLogin(ctx context.Context, name, password string) (string, string, error) {
	r, err := c.service.UserLogin(ctx, &pb.LoginRequest{Name: name, Password: password})
	if err != nil {
		return "", "", ErrInternalServer
	}

	if r.Err != "" {
		return "", "", errors.New(r.Err)
	}
	return r.Token, r.RefreshToken, nil
}
