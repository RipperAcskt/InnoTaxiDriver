package model

import "github.com/google/uuid"

const (
	StatusCreated = "created"
	StatusDeleted = "deleted"
)

type Driver struct {
	ID          uuid.UUID
	Name        string
	PhoneNumber string
	Email       string
	Password    string
	Raiting     float32
	TaxiType    string
}
