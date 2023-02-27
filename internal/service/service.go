package service

import (
	"github.com/RipperAcskt/innotaxidriver/config"
)

type Service struct {
	*AuthService
}

func New(cassandra AuthRepo, client UserSerivce, cfg *config.Config) *Service {
	return &Service{
		AuthService: NewAuthSevice(cassandra, client, cfg),
	}
}
