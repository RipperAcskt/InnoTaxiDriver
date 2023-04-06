package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/RipperAcskt/innotaxidriver/internal/model"
)

type OrderRepo interface {
	FindFree(ctx context.Context) (*model.Driver, error)
}

type OrderService struct {
	OrderRepo
}

type OrderFound struct {
	driver *model.Driver
}

func NewOrdersList(repo OrderRepo) *OrderService {
	return &OrderService{repo}
}

func (o *OrderService) FindDriver(ctx context.Context, id string) (*model.Driver, error) {

	driver, err := o.FindFree(ctx)
	if err != nil && !errors.Is(err, ErrDriverDoesNotExists) {
		return nil, fmt.Errorf("find free faild: %w", err)
	}
	return driver, nil

}
