package service

import (
	"crypto/sha1"
	"fmt"

	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/restapi/operations/auth"
)

var (
	ErrUserAlreadyExists = fmt.Errorf("user alredy exists")

	StatusCreated = "created"
)

type AuthRepo interface {
	CreateDriver(driver auth.PostDriverSingUpBody) error
}

type AuthService struct {
	AuthRepo
	cfg *config.Config
}

func NewAuthSevice(cassandra AuthRepo, cfg *config.Config) *AuthService {
	return &AuthService{cassandra, cfg}
}

func (s *AuthService) SingUp(driver auth.PostDriverSingUpBody) error {
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
