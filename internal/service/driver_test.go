package service_test

import (
	"context"
	"testing"

	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
	"github.com/RipperAcskt/innotaxidriver/internal/service/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestGetProfile(t *testing.T) {
	type mockBehavior func(s *mocks.MockDriverRepo)
	test := []struct {
		name         string
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "get Driver",
			mockBehavior: func(s *mocks.MockDriverRepo) {
				s.EXPECT().GetDriverById(context.Background(), "").Return(&model.Driver{
					Name:        "2",
					PhoneNumber: "2",
					Email:       "2",
					Raiting:     0,
				}, nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			DriverRepo := mocks.NewMockDriverRepo(ctrl)
			DriverService := service.NewDriverService(DriverRepo, &config.Config{})

			tt.mockBehavior(DriverRepo)

			service := service.Service{
				DriverService: DriverService,
			}

			_, err := service.GetProfile(context.Background(), "")
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	type mockBehavior func(s *mocks.MockDriverRepo, Driver model.Driver)

	test := []struct {
		name         string
		Driver       model.Driver
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "update Driver",
			Driver: model.Driver{
				PhoneNumber: "+77777778",
				Email:       "ripper@mail.ru",
			},
			mockBehavior: func(s *mocks.MockDriverRepo, Driver model.Driver) {
				s.EXPECT().UpdateDriverById(context.Background(), Driver).Return(nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			DriverRepo := mocks.NewMockDriverRepo(ctrl)
			DriverService := service.NewDriverService(DriverRepo, &config.Config{})

			tt.mockBehavior(DriverRepo, tt.Driver)

			service := service.Service{
				DriverService: DriverService,
			}

			err := service.UpdateProfile(context.Background(), tt.Driver)
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestDeleteProfile(t *testing.T) {
	type mockBehavior func(s *mocks.MockDriverRepo)

	test := []struct {
		name         string
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "delete Driver",
			mockBehavior: func(s *mocks.MockDriverRepo) {
				s.EXPECT().DeleteDriverById(context.Background(), "").Return(nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			DriverRepo := mocks.NewMockDriverRepo(ctrl)
			DriverService := service.NewDriverService(DriverRepo, &config.Config{})

			tt.mockBehavior(DriverRepo)

			service := service.Service{
				DriverService: DriverService,
			}

			err := service.DeleteProfile(context.Background(), "")
			assert.Equal(t, err, tt.err)
		})
	}
}
