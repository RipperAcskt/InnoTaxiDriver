package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/RipperAcskt/innotaxidriver/internal/model"
)

type OrderRepo interface {
	UpdateStatus(ctx context.Context, drivers []*model.Driver) ([]*model.Driver, error)
}

type OrderService struct {
	OrderRepo
}

func NewOrdersList(repo OrderRepo) *OrderService {
	return &OrderService{repo}
}

func (o *OrderService) SyncDrivers(ctx context.Context, drivers []*model.Driver) ([]*model.Driver, error) {
	newDrivers, err := o.UpdateStatus(ctx, drivers)
	if err != nil && !errors.Is(err, ErrDriverDoesNotExists) {
		return nil, fmt.Errorf("find free faild: %w", err)
	}
	return newDrivers, nil

}
