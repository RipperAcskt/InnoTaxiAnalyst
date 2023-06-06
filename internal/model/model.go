package model

import "github.com/google/uuid"

type AnalysType struct {
	AnalysType string `json:"analys_type"`
}

type ModelType string

const (
	UserType   ModelType = "user"
	DriverType ModelType = "driver"
	OrderType  ModelType = "order"
)

func New(s string) ModelType {
	return ModelType(s)
}

func (t ModelType) ToString() string {
	return string(t)
}

type User struct {
	ID          uuid.UUID `json:"-"`
	UserID      int64     `json:"ID"`
	Name        string
	PhoneNumber string
	Email       string
	Raiting     float64
}

type Driver struct {
	ID          uuid.UUID `json:"-"`
	DriverID    uuid.UUID `json:"ID"`
	Name        string
	PhoneNumber string
	Email       string
	Raiting     float64
	TaxiType    string
}

type Order struct {
	ID           uuid.UUID `json:"-"`
	OrderID      string    `json:"ID"`
	UserID       string
	DriverID     uuid.UUID
	DriverName   string
	DriverPhone  string
	DriverRating float64
	TaxiType     string
	From         string
	To           string
	Date         string
	Status       string
}
