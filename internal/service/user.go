package service

import (
	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
)

type UserRepo interface {
	GetUserById(id string) (*model.Driver, error)
}

type UserService struct {
	UserRepo
	cfg *config.Config
}

func NewUserSevice(cassandra UserRepo, cfg *config.Config) *UserService {
	return &UserService{cassandra, cfg}
}

func (user *UserService) GetProfile(id string) (*model.Driver, error) {
	return user.GetUserById(id)
}
