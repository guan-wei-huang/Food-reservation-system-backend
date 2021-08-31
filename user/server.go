package user

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "reserve_restaurant/user/pb/user"

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
	pb.RegisterUserServiceServer(serv, &grpcServer{s})
	if err = serv.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *grpcServer) NewUser(ctx context.Context, r *pb.NewUserRequest) (*pb.NewUserResponse, error) {
	user, err := s.service.UserRegister(ctx, r.Name, r.Password)
	if err != nil {
		if _, ok := err.(*UserError); !ok {
			log.Fatal(err)
			return nil, err
		}
		return &pb.NewUserResponse{Err: err.Error()}, nil
	}
	return &pb.NewUserResponse{Id: int32(user.ID)}, nil
}

func (s *grpcServer) UserLogin(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, refreshToken, err := s.service.UserLogin(ctx, r.Name, r.Password)
	if err != nil {
		if _, ok := err.(*UserError); !ok {
			log.Fatal(err)
			return nil, err
		}
		return &pb.LoginResponse{Err: err.Error()}, nil
	}
	return &pb.LoginResponse{Token: token, RefreshToken: refreshToken}, nil
}
