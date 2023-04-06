package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
	"github.com/RipperAcskt/innotaxiorder/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	order      *service.OrderService
	listener   net.Listener
	grpcServer *grpc.Server
	log        *zap.Logger
	cfg        *config.Config
}

func New(order *service.OrderService, log *zap.Logger, cfg *config.Config) *Server {
	return &Server{order, nil, nil, log, cfg}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.cfg.GRPC_HOST)

	if err != nil {
		return fmt.Errorf("listen failed: %w", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	s.listener = listener
	s.grpcServer = grpcServer

	proto.RegisterOrderServiceServer(grpcServer, s)
	err = grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("serve failed: %w", err)
	}

	return nil
}

func (s *Server) FindDriver(c context.Context, params *proto.Params) (proto.Response, error) {
	driver, err := s.order.FindDriver(c, params.OrderID)
	if err != nil {
		return proto.Response{}, fmt.Errorf("find driver failed: %w", err)
	}
	response := proto.Response{
		ID:          string(driver.ID),
		Name:        driver.Name,
		PhoneNumber: driver.PhoneNumber,
		Raiting:     driver.Raiting,
	}
	return response, nil
}

func (s *Server) Stop() error {
	s.log.Info("Shuttig down grpc...")

	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("listener close failed: %w", err)
	}

	s.grpcServer.Stop()
	s.log.Info("Grpc server exiting.")
	return nil
}
