package service

import (
	"github.com/RipperAcskt/innotaxidriver/config"
)

type Repo interface {
	DriverRepo
	AuthRepo
}

type Service struct {
	*AuthService
	*DriverService
}

func New(cassandra Repo, client UserSerivce, cfg *config.Config) *Service {
	return &Service{
		AuthService:   NewAuthSevice(cassandra, client, cfg),
		DriverService: NewDriverService(cassandra, cfg),
	}
}
