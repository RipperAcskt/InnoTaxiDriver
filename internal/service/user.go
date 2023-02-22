package service

import (
	"github.com/RipperAcskt/innotaxidriver/config"
)

type UserRepo interface {
	DeleteDriverById(id string) error
}

type UserService struct {
	UserRepo
	cfg *config.Config
}

func NewUserSevice(cassandra UserRepo, cfg *config.Config) *UserService {
	return &UserService{cassandra, cfg}
}

func (user *UserService) DeleteProfile(id string) error {
	return user.DeleteDriverById(id)
}
