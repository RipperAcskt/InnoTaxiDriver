package service

import (
	"github.com/RipperAcskt/innotaxidriver/config"
)

type Repo interface {
	UserRepo
	AuthRepo
}

type Service struct {
	*AuthService
	*UserService
}

func New(cassandra Repo, client UserSerivce, cfg *config.Config) *Service {
	return &Service{
		AuthService: NewAuthSevice(cassandra, client, cfg),
		UserService: NewUserSevice(cassandra, cfg),
	}
}
