package service

import (
	"context"

	"github.com/RipperAcskt/innotaxidriver/config"
)

type UserRepo interface {
	DeleteDriverById(ctx context.Context, id string) error
}

type UserService struct {
	UserRepo
	cfg *config.Config
}

func NewUserSevice(cassandra UserRepo, cfg *config.Config) *UserService {
	return &UserService{cassandra, cfg}
}

func (user *UserService) DeleteProfile(ctx context.Context, id string) error {
	return user.DeleteDriverById(ctx, id)
}
