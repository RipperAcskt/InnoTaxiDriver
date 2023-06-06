package model

import "github.com/google/uuid"

const (
	StatusFree    = "free"
	StatusBusy    = "busy"
	StatusDeleted = "deleted"
)

type Driver struct {
	ID          uuid.UUID
	Name        string
	PhoneNumber string
	Email       string
	Password    string
	Rating      float32
	TaxiType    string
}
