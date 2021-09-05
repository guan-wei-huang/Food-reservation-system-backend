package restaurant

import (
	"context"
	"errors"
	pb "reserve_restaurant/restaurant/pb/restaurant"

	"google.golang.org/grpc"
)

var (
	ErrInternalServer = errors.New("internal server error")
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
		return nil, ErrInternalServer
	}

	if r.Error != "" {
		return nil, errors.New(r.Error)
	}

	foods := []Food{}
	for _, f := range r.Menu.Food {
		ff := &Food{
			Fid:         int(f.Fid),
			Rid:         int(f.Rid),
			Name:        f.FoodName,
			Description: f.Description,
			Price:       f.Price,
		}
		foods = append(foods, *ff)
	}

	menu := &Menu{
		Rid:   int(r.Menu.Rid),
		Foods: foods,
	}
	return menu, nil
}

func (c *Client) CreateFood(ctx context.Context, f *Food) error {
	food := &pb.Food{
		Rid:         int32(f.Rid),
		FoodName:    f.Name,
		Description: f.Description,
		Price:       f.Price,
	}

	r, err := c.service.CreateFood(ctx, &pb.CreateFoodRequest{Food: food})
	if err != nil {
		return ErrInternalServer
	}

	if !r.Complete {
		return errors.New(r.Error)
	}
	return nil
}

func (c *Client) CreateRestaurant(ctx context.Context, r *Restaurant) error {
	restaurant := &pb.Restaurant{
		Name:        r.Name,
		Description: r.Description,
		Location:    r.Location,
	}

	rsp, err := c.service.CreateRestaurant(ctx, &pb.CreateRestReq{Rest: restaurant})
	if err != nil {
		return ErrInternalServer
	}

	if !rsp.Complete {
		return errors.New(rsp.Error)
	}
	return nil
}
