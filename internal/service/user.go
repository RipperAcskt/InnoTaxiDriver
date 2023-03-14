package service

import (
	"context"

	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
)

type UserRepo interface {
	GetUserById(ctx context.Context, id string) (*model.Driver, error)
	UpdateDriverById(ctx context.Context, driver model.Driver) error
	DeleteDriverById(ctx context.Context, id string) error
}

type UserService struct {
	UserRepo
	cfg *config.Config
}

func NewUserSevice(cassandra UserRepo, cfg *config.Config) *UserService {
	return &UserService{cassandra, cfg}
}

func (user *UserService) GetProfile(ctx context.Context, id string) (*model.Driver, error) {
	return user.GetUserById(ctx, id)
}

func (user *UserService) UpdateProfile(ctx context.Context, driver model.Driver) error {
	return user.UpdateDriverById(ctx, driver)
}

func (user *UserService) DeleteProfile(ctx context.Context, id string) error {
	return user.DeleteDriverById(ctx, id)
}
