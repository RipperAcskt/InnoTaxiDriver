package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
	orderProto "github.com/RipperAcskt/innotaxiorder/pkg/proto"
	"github.com/google/uuid"

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
	return &Server{
		order:      order,
		listener:   nil,
		grpcServer: &grpc.Server{},
		log:        log,
		cfg:        cfg,
	}
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
	orderProto.RegisterOrderServiceServer(grpcServer, s)
	err = grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("serve failed: %w", err)
	}

	return nil
}

func (s *Server) SyncDriver(c context.Context, params *orderProto.Info) (*orderProto.Info, error) {
	drivers := make([]*model.Driver, 0, len(params.Drivers))
	for _, driver := range params.Drivers {
		id, err := uuid.Parse(driver.ID)
		if err != nil {
			return nil, fmt.Errorf("parse failed: %w", err)
		}
		tmp := &model.Driver{
			ID: id,
		}
		drivers = append(drivers, tmp)
	}
	newDrivers, err := s.order.SyncDrivers(c, drivers)
	if err != nil {
		return nil, fmt.Errorf("find driver failed: %w", err)
	}

	response := make([]*orderProto.Driver, 0)
	for _, driver := range newDrivers {
		tmp := &orderProto.Driver{
			ID:          driver.ID.String(),
			Name:        driver.Name,
			PhoneNumber: driver.PhoneNumber,
			TaxiType:    driver.TaxiType,
			Raiting:     driver.Raiting,
		}
		response = append(response, tmp)
	}
	return &orderProto.Info{Drivers: response}, nil
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
