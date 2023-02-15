package service

import (
	"crypto/sha1"
	"fmt"

	"github.com/RipperAcskt/innotaxidriver/config"
	grpc "github.com/RipperAcskt/innotaxidriver/internal/grpc"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
)

var (
	ErrDriverAlreadyExists = fmt.Errorf("driver alredy exists")
	ErrDriverDoesNotExists = fmt.Errorf("driver does not exist")
	ErrIncorrectPassword   = fmt.Errorf("incorrect password")
)

type AuthRepo interface {
	CreateDriver(driver model.Driver) error
	CheckUserByPhoneNumber(phone_number string) (*model.Driver, error)
}

type AuthService struct {
	AuthRepo
	client *grpc.Client
	cfg    *config.Config
}

func NewAuthSevice(cassandra AuthRepo, client *grpc.Client, cfg *config.Config) *AuthService {
	return &AuthService{cassandra, client, cfg}
}

func (s *AuthService) SingUp(driver model.Driver) error {
	var err error
	driver.Password, err = s.GenerateHash(driver.Password)
	if err != nil {
		return fmt.Errorf("generate hash failed: %w", err)
	}

	err = s.CreateDriver(driver)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GenerateHash(password string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", fmt.Errorf("write failed: %w", err)
	}
	return string(hash.Sum([]byte(s.cfg.SALT))), nil
}

func (s *AuthService) SingIn(driver model.Driver) error {
	driverDB, err := s.CheckUserByPhoneNumber(driver.PhoneNumber)
	if err != nil {
		return fmt.Errorf("check user by phone number failed: %w", err)
	}

	hash := sha1.New()
	_, err = hash.Write([]byte(driver.Password))
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	if driverDB.Password != string(hash.Sum([]byte(s.cfg.SALT))) {
		return ErrIncorrectPassword
	}

	err = s.client.GetJWT(driverDB.ID)
	if err != nil {
		return fmt.Errorf("get jwt failed: %w", err)
	}

	return nil
}
