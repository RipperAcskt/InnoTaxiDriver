package service

import (
	"github.com/RipperAcskt/innotaxidriver/config"
)

//go:generate mockgen -destination=mocks/mock_auth.go -package=mocks github.com/RipperAcskt/innotaxidriver/internal/service AuthRepo
//go:generate mockgen -destination=mocks/mock_user.go -package=mocks github.com/RipperAcskt/innotaxidriver/internal/service DriverRepo
//go:generate mockgen -destination=mocks/mock_jwt.go -package=mocks github.com/RipperAcskt/innotaxidriver/internal/service UserSerivce
//go:generate mockgen -destination=mocks/mock_broker.go -package=mocks github.com/RipperAcskt/innotaxidriver/internal/service Broker
type Repo interface {
	DriverRepo
	AuthRepo
	OrderRepo
}

type Service struct {
	*AuthService
	*DriverService
	Order *OrderService
}

func New(cassandra Repo, broker Broker, client UserSerivce, cfg *config.Config) *Service {
	return &Service{
		AuthService:   NewAuthSevice(cassandra, broker, client, cfg),
		DriverService: NewDriverService(cassandra, cfg),
		Order:         NewOrdersList(cassandra),
	}
}
