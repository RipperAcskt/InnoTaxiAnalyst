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
	Name        string    `json:"Name"`
	PhoneNumber string    `json:"PhoneNumber"`
	Email       string    `json:"Email"`
	Raiting     float64   `json:"Raiting"`
}

type Driver struct {
	ID          uuid.UUID `json:"-"`
	DriverID    uuid.UUID `json:"ID"`
	Name        string    `json:"Name"`
	PhoneNumber string    `json:"PhoneNumber"`
	Email       string    `json:"Email"`
	Raiting     float64   `json:"Raiting"`
	TaxiType    string    `json:"TaxiType"`
}

type Order struct {
	ID           uuid.UUID `json:"-"`
	OrderID      string    `json:"ID"`
	UserID       string    `json:"UserID"`
	DriverID     uuid.UUID `json:"DriverID"`
	DriverName   string    `json:"DriverName"`
	DriverPhone  string    `json:"DriverPhone"`
	DriverRating float64   `json:"DriverRating"`
	TaxiType     string    `json:"TaxiType"`
	From         string    `json:"From"`
	To           string    `json:"To"`
	Date         string    `json:"Date"`
	Status       string    `json:"Status"`
}

type Rating struct {
	Type   string  `json:"-"`
	ID     string  `json:"ID"`
	Rating float32 `json:"Rating"`
}
