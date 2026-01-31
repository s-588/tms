package models

import (
	"time"
)

type Client struct {
	ClientID             int32
	Name                 string
	Email                string
	EmailVerified        bool
	EmailToken           *string
	EmailTokenExpiration time.Time
	Phone                string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
	Orders               []Order
}

type Employee struct {
	EmployeeID int32
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

type Transport struct {
	TransportID     int32
	EmployeeID      *int32
	Model           string
	LicensePlate    *string
	PayloadCapacity int32
	FuelID          int32
	FuelConsumption int32
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}

type Fuel struct {
	FuelID    int32
	Name      string
	Supplier  *string
	Price     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Order struct {
	OrderID    int32
	Distance   int32
	Weight     int32
	TotalPrice string
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
	Transports []Transport
}

type Price struct {
	PriceID   int32
	CargoType string
	Cost      string
	Weight    int32
	Distance  int32
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
