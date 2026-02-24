package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type ClientStatus string
type EmployeeStatus string
type OrderStatus string
type InspectionStatus string
type EmployeeJobTitle string

const (
	ClientStatusActive   ClientStatus = "active"
	ClientStatusInactive ClientStatus = "inactive"
	ClientStatusBlocked  ClientStatus = "blocked"
	ClientStatusNew      ClientStatus = "new"

	EmployeeStatusAvailable   EmployeeStatus = "available"
	EmployeeStatusAssigned    EmployeeStatus = "assigned"
	EmployeeStatusUnavailable EmployeeStatus = "unavailable"

	EmployeeJobTitleDriver           EmployeeJobTitle = "driver"
	EmployeeJobTitleDispatcher       EmployeeJobTitle = "dispatcher"
	EmployeeJobTitleMechanic         EmployeeJobTitle = "mechanic"
	EmployeeJobTitleLogisticsManager EmployeeJobTitle = "logistics_manager"

	OrderStatusPending    OrderStatus = "pending"
	OrderStatusAssigned   OrderStatus = "assigned"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"

	InspectionStatusReady   InspectionStatus = "ready"
	InspectionStatusRepair  InspectionStatus = "repair"
	InspectionStatusOverdue InspectionStatus = "overdue"
)

type Client struct {
	ClientID             int32
	Name                 string
	Email                string
	EmailVerified        bool
	EmailToken           string
	EmailTokenExpiration time.Time
	Phone                string
	Score                uint8
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
	Orders               []Order
}

type Employee struct {
	EmployeeID        int32
	Name              string
	Status            EmployeeStatus
	JobTitle          EmployeeJobTitle
	HireDate          time.Time
	Salary            decimal.Decimal
	LicenseIssued     time.Time
	LicenseExpiration time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
}

type Transport struct {
	TransportID      int32
	Model            string
	LicensePlate     string
	PayloadCapacity  int32
	FuelConsumption  int32
	InspectionPassed bool
	InspectionDate   time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}

type Order struct {
	OrderID     int32
	ClientID    int32
	TransportID int32
	EmployeeID  int32
	PriceID     int32
	Grade       uint8
	Distance    int32
	Weight      int32
	TotalPrice  decimal.Decimal
	Status      OrderStatus
	NodeIDStart *int32
	NodeIDEnd   *int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type Price struct {
	PriceID   int32
	CargoType string
	Weight    decimal.Decimal
	Distance  decimal.Decimal
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// New models for additional tables (if needed)
type Insurance struct {
	InsuranceID         int32
	TransportID         int32
	InsuranceDate       time.Time
	InsuranceExpiration time.Time
	Payment             decimal.Decimal
	Coverage            decimal.Decimal
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           time.Time
}

type Inspection struct {
	InspectionID         int32
	TransportID          int32
	InspectionDate       time.Time
	InspectionExpiration time.Time
	Status               InspectionStatus
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
}

type Node struct {
	NodeID    int32
	Name      string
	Geom      Point
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Point struct {
	X, Y float64
}
