package client

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxidriver/config"
	pb "github.com/RipperAcskt/innotaxidriver/pkg/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
	cfg    *config.Config
}

type Token struct {
	AccessToken  string
	RefreshToken string
}

func New(cfg *config.Config) (*Client, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(cfg.GRPC_AUTH_HOST, opts...)

	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	client := pb.NewAuthServiceClient(conn)

	return &Client{client, conn, cfg}, nil
}

func (c *Client) GetJWT(id uuid.UUID) (*Token, error) {
	request := &pb.Params{
		DriverID: id.String(),
		Type:     "driver",
	}
	response, err := c.client.GetJWT(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("get jwt failed: %w", err)
	}
	return &Token{response.AccessToken, response.RefreshToken}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
