package service

import (
	"context"

	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
)

type DriverRepo interface {
	GetDriverById(ctx context.Context, id string) (*model.Driver, error)
	UpdateDriverById(ctx context.Context, driver model.Driver) error
	DeleteDriverById(ctx context.Context, id string) error
	SetRatingById(ctx context.Context, id string, raiting float32) error
	GetRating(ctx context.Context) ([]model.Driver, error)
}

type DriverService struct {
	DriverRepo
	cfg *config.Config
}

func NewDriverService(cassandra DriverRepo, cfg *config.Config) *DriverService {
	return &DriverService{cassandra, cfg}
}

func (Driver *DriverService) GetProfile(ctx context.Context, id string) (*model.Driver, error) {
	return Driver.GetDriverById(ctx, id)
}

func (Driver *DriverService) UpdateProfile(ctx context.Context, driver model.Driver) error {
	return Driver.UpdateDriverById(ctx, driver)
}

func (Driver *DriverService) DeleteProfile(ctx context.Context, id string) error {
	return Driver.DeleteDriverById(ctx, id)
}

func (Driver *DriverService) SetRaiting(ctx context.Context, id string, raiting float32) error {
	return Driver.SetRatingById(ctx, id, raiting)
}

func (Driver *DriverService) GetRaiting(ctx context.Context) ([]model.Driver, error) {
	return Driver.GetRating(ctx)
}
