package service

import (
	"github.com/RipperAcskt/innotaxidriver/config"
	grpc "github.com/RipperAcskt/innotaxidriver/internal/gateway"
)

type Service struct {
	*AuthService
}

func New(cassandra AuthRepo, client *grpc.Client, cfg *config.Config) *Service {
	return &Service{
		AuthService: NewAuthSevice(cassandra, client, cfg),
	}
}
