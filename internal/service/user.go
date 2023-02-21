package service

import (
	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
)

type UserRepo interface {
	UpdateDriverById(driver model.Driver) error
}

type UserService struct {
	UserRepo
	cfg *config.Config
}

func NewUserSevice(cassandra UserRepo, cfg *config.Config) *UserService {
	return &UserService{cassandra, cfg}
}

func (user *UserService) UpdateProfile(driver model.Driver) error {
	return user.UpdateDriverById(driver)
}
