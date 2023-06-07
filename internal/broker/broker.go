package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Broker struct {
	conn *kafka.Conn
	cfg  *config.Config
}

type Driver struct {
	ID          uuid.UUID
	Name        string
	PhoneNumber string
	Email       string
	Rating      float32
	TaxiType    string
}

func New(cfg *config.Config) (*Broker, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.BROKER_HOST, "driver", 0)
	if err != nil {
		return nil, fmt.Errorf("dial leader failed: %w", err)
	}

	return &Broker{
		conn: conn,
		cfg:  cfg,
	}, nil
}

func (b *Broker) Write(driver model.Driver) error {
	packedDriver := Driver{
		ID:          driver.ID,
		Name:        driver.Name,
		PhoneNumber: driver.PhoneNumber,
		Email:       driver.Email,
		Rating:      driver.Rating,
		TaxiType:    driver.TaxiType,
	}

	data, err := json.Marshal(packedDriver)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	_, err = b.conn.Write(data)
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	return nil
}
