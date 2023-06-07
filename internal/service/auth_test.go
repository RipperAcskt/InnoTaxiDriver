package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/broker"
	user "github.com/RipperAcskt/innotaxidriver/internal/client"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
	"github.com/RipperAcskt/innotaxidriver/internal/service/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestSingUp(t *testing.T) {
	type mockBehavior func(s *mocks.MockAuthRepo, user model.Driver, b *mocks.MockBroker)
	type fileds struct {
		authRepo    *mocks.MockAuthRepo
		userService *mocks.MockUserSerivce
		broker      *mocks.MockBroker
	}
	test := []struct {
		name         string
		user         model.Driver
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "correct user",
			user: model.Driver{
				Name:        "Ivan",
				PhoneNumber: "+7455456",
				Email:       "ripper@algsdh",
				Password:    "12345",
				TaxiType:    "econom",
			},
			mockBehavior: func(s *mocks.MockAuthRepo, user model.Driver, b *mocks.MockBroker) {
				s.EXPECT().CreateDriver(context.Background(), user).Return(nil)
				s.EXPECT().CheckDriverByPhoneNumber(context.Background(), "+7455456").Return(&model.Driver{ID: uuid.MustParse("eba0d15e-c710-4c08-af1d-dedbcf5ad6ca")}, nil)
				b.EXPECT().Write(model.Driver{
					ID:          uuid.MustParse("eba0d15e-c710-4c08-af1d-dedbcf5ad6ca"),
					Name:        "Ivan",
					PhoneNumber: "+7455456",
					Email:       "ripper@algsdh",
					TaxiType:    "econom",
				}).Return(nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fileds{
				authRepo:    mocks.NewMockAuthRepo(ctrl),
				userService: mocks.NewMockUserSerivce(ctrl),
				broker:      mocks.NewMockBroker(ctrl),
			}

			service := service.Service{
				AuthService: service.NewAuthSevice(f.authRepo, f.broker, f.userService, &config.Config{}),
			}

			tmpPass := tt.user.Password
			tt.user.Password, _ = service.GenerateHash(tt.user.Password)
			tt.mockBehavior(f.authRepo, tt.user, f.broker)

			tt.user.Password = tmpPass
			err := service.SingUp(context.Background(), tt.user)
			assert.IsEqual(err, tt.err)
		})
	}
}

func TestSingIn(t *testing.T) {
	type mockBehavior func(s *mocks.MockAuthRepo, phone_number string)
	type fileds struct {
		authRepo    *mocks.MockAuthRepo
		userService *mocks.MockUserSerivce
	}
	test := []struct {
		name         string
		user         model.Driver
		mockBehavior mockBehavior
		token        *user.Token
		err          error
	}{
		{
			name: "correct password",
			user: model.Driver{
				PhoneNumber: "2",
				Password:    "2",
			},
			mockBehavior: func(s *mocks.MockAuthRepo, phone_number string) {
				s.EXPECT().CheckDriverByPhoneNumber(context.Background(), phone_number).Return(&model.Driver{
					PhoneNumber: "2",
					Password:    string([]byte{49, 50, 52, 106, 107, 104, 115, 100, 97, 102, 51, 52, 50, 53, 218, 75, 146, 55, 186, 204, 205, 241, 156, 7, 96, 202, 183, 174, 196, 168, 53, 144, 16, 176}),
				}, nil)
			},
			token: nil,
			err:   nil,
		},
		{
			name: "incorrect password",
			user: model.Driver{
				PhoneNumber: "+7455456",
				Password:    "123456",
			},
			mockBehavior: func(s *mocks.MockAuthRepo, phone_number string) {
				s.EXPECT().CheckDriverByPhoneNumber(context.Background(), phone_number).Return(&model.Driver{}, nil)
			},
			token: nil,
			err:   fmt.Errorf("incorrect password"),
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fileds{
				authRepo:    mocks.NewMockAuthRepo(ctrl),
				userService: mocks.NewMockUserSerivce(ctrl),
			}
			authService := service.NewAuthSevice(f.authRepo, &broker.Broker{}, f.userService, &config.Config{})

			tt.mockBehavior(f.authRepo, tt.user.PhoneNumber)

			service := service.Service{
				AuthService: authService,
			}
			token, _ := service.SingIn(context.Background(), tt.user)
			assert.Equal(t, token, tt.token)
		})
	}

}

func TestGenerateHash(t *testing.T) {
	type fileds struct {
		authRepo    *mocks.MockAuthRepo
		userService *mocks.MockUserSerivce
	}
	test := []struct {
		name     string
		password string
		hash     string
		err      error
	}{
		{
			name:     "password",
			password: "2",
			hash:     string([]byte{218, 75, 146, 55, 186, 204, 205, 241, 156, 7, 96, 202, 183, 174, 196, 168, 53, 144, 16, 176}),
			err:      nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fileds{
				authRepo:    mocks.NewMockAuthRepo(ctrl),
				userService: mocks.NewMockUserSerivce(ctrl),
			}
			authService := service.NewAuthSevice(f.authRepo, &broker.Broker{}, f.userService, &config.Config{})

			service := service.Service{
				AuthService: authService,
			}

			hash, err := service.GenerateHash(tt.password)
			assert.Equal(t, hash, tt.hash)
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestCheckToken(t *testing.T) {
	test := []struct {
		name  string
		exist error
	}{
		{
			name:  "check token",
			exist: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			_, err := service.Verify("0", &config.Config{})
			assert.NotEqual(t, err, tt.exist)
		})
	}
}
func TestVerify(t *testing.T) {
	cfg := &config.Config{
		HS256_SECRET: "QWERTfg53gxb2",
	}

	test := []struct {
		name   string
		token  string
		userId string
		err    error
	}{
		{
			name:   "verify token expired",
			token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzY4Nzk5NDIsInR5cGUiOiJ1c2VyIiwidXNlcl9pZCI6MX0.qwiL4bupjm9O-ZnKpIcB8-erQytBJgkWlxnwPmRmv-c",
			userId: "",
			err:    service.ErrTokenExpired,
		},
		{
			name:   "verify token ok",
			token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Nzk0NzAxNTQsInR5cGUiOiJ1c2VyIiwidXNlcl9pZCI6MX0.r5vZu9eOds5kti9UjQFXx8AYLHZC23YLtVVnr8dgx24",
			userId: "",
			err:    nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id, err := service.Verify(tt.token, cfg)
			assert.IsEqual(err, tt.err)
			assert.Equal(t, id, tt.userId)
		})
	}
}
