package user

import (
	"context"
	"fmt"

	pb "github.com/RipperAcskt/innotaxi/pkg/proto"
	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type User struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
	cfg    *config.Config
}

type Token struct {
	AccessToken  string
	RefreshToken string
}

func New(cfg *config.Config) (*User, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(cfg.GRPC_USER_SERVICE_HOST, opts...)

	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	client := pb.NewAuthServiceClient(conn)

	return &User{client, conn, cfg}, nil
}

func (u *User) GetJWT(id uuid.UUID) (*Token, error) {
	request := &pb.Params{
		DriverID: id.String(),
		Type:     "driver",
	}
	response, err := u.client.GetJWT(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("get jwt failed: %w", err)
	}
	return &Token{response.AccessToken, response.RefreshToken}, nil
}

func (u *User) Close() error {
	return u.conn.Close()
}
