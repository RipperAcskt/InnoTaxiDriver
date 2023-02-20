package service

import (
	"github.com/RipperAcskt/innotaxidriver/config"
)

type Service struct {
	*AuthService
}

func New(cassandra AuthRepo, cfg *config.Config) *Service {
	return &Service{
		AuthService: NewAuthSevice(cassandra, cfg),
	}
}
