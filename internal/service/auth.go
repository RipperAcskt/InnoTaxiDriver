package service

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxidriver/config"
	grpc "github.com/RipperAcskt/innotaxidriver/internal/gateway"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
	"github.com/golang-jwt/jwt"
)

var (
	ErrDriverAlreadyExists = fmt.Errorf("driver alredy exists")
	ErrDriverDoesNotExists = fmt.Errorf("driver does not exist")
	ErrIncorrectPassword   = fmt.Errorf("incorrect password")
	ErrTokenExpired        = fmt.Errorf("token expired")
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

func (s *AuthService) SingIn(driver model.Driver) (*grpc.Token, error) {
	driverDB, err := s.CheckUserByPhoneNumber(driver.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("check user by phone number failed: %w", err)
	}

	hash := sha1.New()
	_, err = hash.Write([]byte(driver.Password))
	if err != nil {
		return nil, fmt.Errorf("write failed: %w", err)
	}

	if driverDB.Password != string(hash.Sum([]byte(s.cfg.SALT))) {
		return nil, ErrIncorrectPassword
	}

	token, err := s.client.GetJWT(driverDB.ID)
	if err != nil {
		return nil, fmt.Errorf("get jwt failed: %w", err)
	}

	return token, nil
}

func Verify(token string, cfg *config.Config) (string, error) {
	tokenJwt, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.HS256_SECRET), nil
		},
	)

	if err != nil {
		return "", fmt.Errorf("token parse failed: %w", err)
	}

	claims, ok := tokenJwt.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("jwt map claims failed")
	}

	if !claims.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
		return "", ErrTokenExpired
	}
	return string(claims["user_id"].(string)), nil
}

func (s *AuthService) Refresh(driver model.Driver) (*grpc.Token, error) {
	token, err := s.client.GetJWT(driver.ID)
	if err != nil {
		return nil, fmt.Errorf("get jwt failed: %w", err)
	}

	return token, nil
}
