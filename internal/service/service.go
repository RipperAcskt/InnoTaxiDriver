package service

import (
	"github.com/RipperAcskt/innotaxidriver/config"
)

//go:generate mockgen -destination=mocks/mock_auth.go -package=mocks github.com/RipperAcskt/innotaxidriver/internal/service AuthRepo
//go:generate mockgen -destination=mocks/mock_user.go -package=mocks github.com/RipperAcskt/innotaxidriver/internal/service UserRepo
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
