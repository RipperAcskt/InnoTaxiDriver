package grpc_client

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxidriver/config"
	pb "github.com/RipperAcskt/innotaxidriver/internal/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client pb.AuthServiceClient
	cfg    *config.Config
}

func New(cfg *config.Config) (*Client, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(cfg.GRPC_HOST, opts...)

	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	return &Client{client, cfg}, nil
}

func (c *Client) GetJWT(id uuid.UUID) error {
	request := &pb.Params{
		DriverID: id.String(),
		Type:     "driver",
	}
	response, err := c.client.GetJWT(context.Background(), request)
	if err != nil {
		return fmt.Errorf("get jwt failed: %w", err)
	}
	fmt.Println(response)
	return nil
}
