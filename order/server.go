package order

import (
	"context"
	"fmt"
	"net"

	pb "reserve_restaurant/order/pb/order"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcServer struct {
	logger  *zap.SugaredLogger
	service Service
}

func ListenGRPC(s Service, config *Config) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return nil
	}

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{sugar, s})
	if err = serv.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *grpcServer) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	oid, err := s.service.CreateOrder(ctx, parseOrder(r.Order))
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &pb.CreateOrderResponse{Id: int32(oid)}, nil
}

func (s *grpcServer) GetOrder(ctx context.Context, r *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := s.service.GetOrder(ctx, int(r.Id), int(r.Uid))
	if err != nil {
		switch err {
		case ErrInvalidAccess, ErrOrderInvalid:
			s.logger.Infof("client GetOrder failed: %v", err)
			return &pb.GetOrderResponse{Error: err.Error()}, nil
		default:
			s.logger.Error(err.Error())
			return nil, err
		}
	}
	return &pb.GetOrderResponse{Order: formatOrder(order)}, nil
}

func (s *grpcServer) GetOrderForUser(ctx context.Context, r *pb.GetOrderRequest) (*pb.GetOrderForUserResponse, error) {
	orders, err := s.service.GetOrderForUser(ctx, int(r.Uid))
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	os := []*pb.Order{}
	for _, o := range *orders {
		os = append(os, formatOrder(&o))
	}

	return &pb.GetOrderForUserResponse{Orders: os}, nil
}

func parseOrder(po *pb.Order) *Order {
	order := &Order{
		Id:        int(po.Id),
		Rid:       int(po.Rid),
		Uid:       int(po.Uid),
		CreatedAt: po.CreateAt.AsTime(),
	}

	products := []Product{}
	for _, p := range po.Products {
		product := Product{
			Fid:      int(p.Fid),
			Price:    p.Price,
			Name:     p.Name,
			Quantity: int(p.Quantity),
		}
		products = append(products, product)
	}
	order.Products = products
	return order
}

func formatOrder(o *Order) *pb.Order {
	order := &pb.Order{
		Id:       int32(o.Id),
		Rid:      int32(o.Rid),
		Uid:      int32(o.Uid),
		CreateAt: timestamppb.New(o.CreatedAt),
	}

	products := []*pb.Order_Product{}
	for _, p := range o.Products {
		product := &pb.Order_Product{
			Fid:      int32(p.Fid),
			Name:     p.Name,
			Price:    p.Price,
			Quantity: int32(p.Quantity),
		}
		products = append(products, product)
	}
	order.Products = products
	return order
}
