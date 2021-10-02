package restaurant

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "reserve_restaurant/restaurant/pb/restaurant"

	"google.golang.org/grpc"
)

var (
	GOOGLE_API_KEY string
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
			log.Println("get restaurant menu failed:", err)
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

	rid, err := s.service.CreateRestaurant(ctx, rest)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &pb.GeneralResponse{Complete: true, Id: int32(rid)}, nil
}

func (s *grpcServer) SearchRestaurant(ctx context.Context, r *pb.SearchRestaurantReq) (*pb.SearchRestaurantResp, error) {
	rests, err := s.service.SearchRestaurant(ctx, r.Location)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	restaurants := []*pb.Restaurant{}
	for _, r := range rests {
		restaurants = append(restaurants, formatRestaurant(r))
	}
	return &pb.SearchRestaurantResp{Restaurants: restaurants}, nil
}

func formatRestaurant(rest *Restaurant) *pb.Restaurant {
	r := &pb.Restaurant{
		Id:          int32(rest.ID),
		Name:        rest.Name,
		Description: rest.Description,
		Location:    rest.Location,
	}
	return r
}

func parseRestaurant(pbr *pb.Restaurant) *Restaurant {
	r := &Restaurant{
		ID:          int(pbr.Id),
		Name:        pbr.Name,
		Description: pbr.Description,
		Location:    pbr.Location,
	}
	return r
}
