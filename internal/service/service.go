package service

import (
	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/handler/grpc"
)

type Repo interface {
	UserRepo
	AuthRepo
}

type Service struct {
	*AuthService
	*UserService
}

func New(cassandra Repo, client *grpc.Client, cfg *config.Config) *Service {
	return &Service{
		AuthService: NewAuthSevice(cassandra, client, cfg),
		UserService: NewUserSevice(cassandra, cfg),
	}
}
