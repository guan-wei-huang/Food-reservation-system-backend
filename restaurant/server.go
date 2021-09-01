package restaurant

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "reserve_restaurant/restaurant/pb/restaurant"

	"google.golang.org/grpc"
)

type grpcServer struct {
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterRestaurantServiceServer(serv, &grpcServer{s})
	if err = serv.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *grpcServer) GetRestaurantMenu(ctx context.Context, r *pb.MenuRequest) (*pb.MenuResponse, error) {
	m, err := s.service.GetRestaurantMenu(ctx, int(r.Rid))

	if err != nil {
		switch err {
		case ErrRestaurantIdWrong:
			return &pb.MenuResponse{Error: err.Error()}, nil
		default:
			return nil, err
		}
	}

	menu := &pb.Menu{Rid: int32(m.Rid)}
	for _, f := range m.Foods {
		food := &pb.Food{
			Fid:         int32(f.Fid),
			Rid:         int32(f.Rid),
			FoodName:    f.Name,
			Description: f.Description,
			Price:       f.Price,
		}
		menu.Food = append(menu.Food, food)
	}

	return &pb.MenuResponse{Menu: menu}, nil
}

func (s *grpcServer) CreateFood(ctx context.Context, r *pb.CreateFoodRequest) (*pb.GeneralResponse, error) {
	f := &Food{
		Fid:         int(r.Food.Fid),
		Rid:         int(r.Food.Rid),
		Name:        r.Food.FoodName,
		Description: r.Food.Description,
		Price:       r.Food.Price,
	}

	if err := s.service.CreateFood(ctx, f); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &pb.GeneralResponse{Complete: true}, nil
}

func (s *grpcServer) CreateRestaurant(ctx context.Context, r *pb.CreateRestReq) (*pb.GeneralResponse, error) {
	rest := &Restaurant{
		Name:        r.Rest.Name,
		Description: r.Rest.Description,
		Location:    r.Rest.Location,
	}

	if err := s.service.CreateRestaurant(ctx, rest); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &pb.GeneralResponse{Complete: true}, nil
}
