package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxidriver/config"
	user "github.com/RipperAcskt/innotaxidriver/internal/client"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	ErrDriverAlreadyExists = fmt.Errorf("driver alredy exists")
	ErrDriverDoesNotExists = fmt.Errorf("driver does not exist")
	ErrIncorrectPassword   = fmt.Errorf("incorrect password")
	ErrTokenExpired        = fmt.Errorf("token expired")
	ErrTokenClaims         = fmt.Errorf("jwt map claims failed")
	ErrTokenId             = fmt.Errorf("jwt get user id failed failed")
)

type AuthRepo interface {
	CreateDriver(ctx context.Context, driver model.Driver) error
	CheckDriverByPhoneNumber(ctx context.Context, phone_number string) (*model.Driver, error)
}

type UserSerivce interface {
	GetJWT(ctx context.Context, id uuid.UUID) (*user.Token, error)
}

type AuthService struct {
	AuthRepo
	user UserSerivce
	cfg  *config.Config
}

func NewAuthSevice(cassandra AuthRepo, user UserSerivce, cfg *config.Config) *AuthService {
	return &AuthService{cassandra, user, cfg}
}

func (s *AuthService) SingUp(ctx context.Context, driver model.Driver) error {
	var err error
	driver.Password, err = s.GenerateHash(driver.Password)
	if err != nil {
		return fmt.Errorf("generate hash failed: %w", err)
	}

	err = s.CreateDriver(ctx, driver)
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

func (s *AuthService) SingIn(ctx context.Context, driver model.Driver) (*user.Token, error) {
	driverDB, err := s.CheckDriverByPhoneNumber(ctx, driver.PhoneNumber)
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

	token, err := s.user.GetJWT(ctx, driverDB.ID)
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
		return "", ErrTokenClaims
	}

	if !claims.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
		return "", ErrTokenExpired
	}
	id, ok := claims["user_id"]
	if !ok {
		return "", ErrTokenId
	}
	str, ok := id.(string)
	if !ok {
		return "", ErrTokenId
	}
	return string(str), nil
}

func (s *AuthService) Refresh(ctx context.Context, driver model.Driver) (*user.Token, error) {
	token, err := s.user.GetJWT(ctx, driver.ID)
	if err != nil {
		return nil, fmt.Errorf("get jwt failed: %w", err)
	}

	return token, nil
}
