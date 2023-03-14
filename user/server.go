package user

import (
	"context"
	"fmt"
	"net"

	pb "reserve_restaurant/user/pb/user"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type grpcServer struct {
	logger  *zap.SugaredLogger
	service Service
}

func ListenGRPC(s Service, config *Config) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return err
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{config.LogFile}
	logger, _ := zapConfig.Build()
	sugar := logger.Sugar()
	defer sugar.Sync() // make sure buffer clean

	serv := grpc.NewServer()
	pb.RegisterUserServiceServer(serv, &grpcServer{sugar, s})
	if err = serv.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *grpcServer) NewUser(ctx context.Context, r *pb.NewUserRequest) (*pb.NewUserResponse, error) {
	user, err := s.service.UserRegister(ctx, r.Name, r.Password)
	if err != nil {
		if _, ok := err.(*UserError); !ok {
			s.logger.Error(err.Error())
			return nil, err
		}
		s.logger.Infof("client NewUser failed: %v", err)
		return &pb.NewUserResponse{Err: err.Error()}, nil
	}
	return &pb.NewUserResponse{Id: int32(user.ID)}, nil
}

func (s *grpcServer) UserLogin(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, refreshToken, err := s.service.UserLogin(ctx, r.Name, r.Password)
	if err != nil {
		if _, ok := err.(*UserError); !ok {
			s.logger.Error(err.Error())
			return nil, err
		}
		s.logger.Infof("client UserLogin failed: %v", err)
		return &pb.LoginResponse{Err: err.Error()}, nil
	}
	return &pb.LoginResponse{Token: token, RefreshToken: refreshToken}, nil
}
