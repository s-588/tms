// Package grpc provides the gRPC server implementation for the TMS service.
package grpc

import (
	"time"

	"github.com/s-588/tms/cmd/models"
	pb "github.com/s-588/tms/proto"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// toProtoTimestamp converts a time.Time to a protobuf Timestamp.
// It returns nil if the time is zero.
func toProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

// fromProtoTimestamp converts a protobuf Timestamp to time.Time.
// It returns the zero time if the pointer is nil.
func fromProtoTimestamp(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

// toProtoDecimal converts a decimal.Decimal to its string representation.
func toProtoDecimal(d decimal.Decimal) string {
	return d.String()
}

// fromProtoDecimal parses a string into a decimal.Decimal.
// It returns decimal.Zero on empty string or parse error (error ignored).
func fromProtoDecimal(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	d, _ := decimal.NewFromString(s)
	return d
}

// toProtoOptionalString converts an Optional[string] to a *string pointer.
// It returns nil if the optional is not set.
func toProtoOptionalString(o models.Optional[string]) *string {
	if o.Set {
		return &o.Value
	}
	return nil
}

// fromProtoOptionalString converts a *string pointer to an Optional[string].
// It leaves the optional unset if the pointer is nil.
func fromProtoOptionalString(p *string) models.Optional[string] {
	var o models.Optional[string]
	if p != nil {
		o.SetValue(*p)
	}
	return o
}

// toProtoEmployeeStatus maps a domain EmployeeStatus to its protobuf enum.
func toProtoEmployeeStatus(s models.EmployeeStatus) pb.EmployeeStatus {
	switch s {
	case models.EmployeeStatusAvailable:
		return pb.EmployeeStatus_EMPLOYEE_STATUS_AVAILABLE
	case models.EmployeeStatusAssigned:
		return pb.EmployeeStatus_EMPLOYEE_STATUS_ASSIGNED
	case models.EmployeeStatusUnavailable:
		return pb.EmployeeStatus_EMPLOYEE_STATUS_UNAVAILABLE
	default:
		return pb.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED
	}
}

// fromProtoEmployeeStatus maps a protobuf EmployeeStatus to its domain enum.
// It defaults to Available for unspecified values.
func fromProtoEmployeeStatus(s pb.EmployeeStatus) models.EmployeeStatus {
	switch s {
	case pb.EmployeeStatus_EMPLOYEE_STATUS_AVAILABLE:
		return models.EmployeeStatusAvailable
	case pb.EmployeeStatus_EMPLOYEE_STATUS_ASSIGNED:
		return models.EmployeeStatusAssigned
	case pb.EmployeeStatus_EMPLOYEE_STATUS_UNAVAILABLE:
		return models.EmployeeStatusUnavailable
	default:
		return models.EmployeeStatusAvailable
	}
}

// toProtoEmployeeJobTitle maps a domain EmployeeJobTitle to protobuf enum.
func toProtoEmployeeJobTitle(t models.EmployeeJobTitle) pb.EmployeeJobTitle {
	switch t {
	case models.EmployeeJobTitleDriver:
		return pb.EmployeeJobTitle_EMPLOYEE_JOB_TITLE_DRIVER
	case models.EmployeeJobTitleDispatcher:
		return pb.EmployeeJobTitle_EMPLOYEE_JOB_TITLE_DISPATCHER
	case models.EmployeeJobTitleMechanic:
		return pb.EmployeeJobTitle_EMPLOYEE_JOB_TITLE_MECHANIC
	case models.EmployeeJobTitleLogisticsManager:
		return pb.EmployeeJobTitle_EMPLOYEE_JOB_TITLE_LOGISTICS_MANAGER
	default:
		return pb.EmployeeJobTitle_EMPLOYEE_JOB_TITLE_UNSPECIFIED
	}
}

// fromProtoEmployeeJobTitle maps a protobuf EmployeeJobTitle to domain enum.
// It defaults to Driver for unspecified values.
func fromProtoEmployeeJobTitle(t pb.EmployeeJobTitle) models.EmployeeJobTitle {
	switch t {
	case pb.EmployeeJobTitle_EMPLOYEE_JOB_TITLE_DRIVER:
		return models.EmployeeJobTitleDriver
	case pb.EmployeeJobTitle_EMPLOYEE_JOB_TITLE_DISPATCHER:
		return models.EmployeeJobTitleDispatcher
	case pb.EmployeeJobTitle_EMPLOYEE_JOB_TITLE_MECHANIC:
		return models.EmployeeJobTitleMechanic
	case pb.EmployeeJobTitle_EMPLOYEE_JOB_TITLE_LOGISTICS_MANAGER:
		return models.EmployeeJobTitleLogisticsManager
	default:
		return models.EmployeeJobTitleDriver
	}
}

// toProtoOrderStatus maps domain OrderStatus to protobuf enum.
func toProtoOrderStatus(s models.OrderStatus) pb.OrderStatus {
	switch s {
	case models.OrderStatusPending:
		return pb.OrderStatus_ORDER_STATUS_PENDING
	case models.OrderStatusAssigned:
		return pb.OrderStatus_ORDER_STATUS_ASSIGNED
	case models.OrderStatusInProgress:
		return pb.OrderStatus_ORDER_STATUS_IN_PROGRESS
	case models.OrderStatusCompleted:
		return pb.OrderStatus_ORDER_STATUS_COMPLETED
	case models.OrderStatusCancelled:
		return pb.OrderStatus_ORDER_STATUS_CANCELLED
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

// fromProtoOrderStatus maps protobuf OrderStatus to domain enum.
// It defaults to Pending for unspecified values.
func fromProtoOrderStatus(s pb.OrderStatus) models.OrderStatus {
	switch s {
	case pb.OrderStatus_ORDER_STATUS_PENDING:
		return models.OrderStatusPending
	case pb.OrderStatus_ORDER_STATUS_ASSIGNED:
		return models.OrderStatusAssigned
	case pb.OrderStatus_ORDER_STATUS_IN_PROGRESS:
		return models.OrderStatusInProgress
	case pb.OrderStatus_ORDER_STATUS_COMPLETED:
		return models.OrderStatusCompleted
	case pb.OrderStatus_ORDER_STATUS_CANCELLED:
		return models.OrderStatusCancelled
	default:
		return models.OrderStatusPending
	}
}

// toProtoInspectionStatus maps domain InspectionStatus to protobuf enum.
func toProtoInspectionStatus(s models.InspectionStatus) pb.InspectionStatus {
	switch s {
	case models.InspectionStatusReady:
		return pb.InspectionStatus_INSPECTION_STATUS_READY
	case models.InspectionStatusRepair:
		return pb.InspectionStatus_INSPECTION_STATUS_REPAIR
	case models.InspectionStatusOverdue:
		return pb.InspectionStatus_INSPECTION_STATUS_OVERDUE
	default:
		return pb.InspectionStatus_INSPECTION_STATUS_UNSPECIFIED
	}
}

// fromProtoInspectionStatus maps protobuf InspectionStatus to domain enum.
// It defaults to Ready for unspecified values.
func fromProtoInspectionStatus(s pb.InspectionStatus) models.InspectionStatus {
	switch s {
	case pb.InspectionStatus_INSPECTION_STATUS_READY:
		return models.InspectionStatusReady
	case pb.InspectionStatus_INSPECTION_STATUS_REPAIR:
		return models.InspectionStatusRepair
	case pb.InspectionStatus_INSPECTION_STATUS_OVERDUE:
		return models.InspectionStatusOverdue
	default:
		return models.InspectionStatusReady
	}
}

// toProtoClient converts a domain Client to a protobuf Client message.
func toProtoClient(c models.Client) *pb.Client {
	return &pb.Client{
		ClientId:      c.ClientID,
		Name:          c.Name,
		Email:         c.Email,
		EmailVerified: c.EmailVerified,
		Phone:         c.Phone,
		Score:         int32(c.Score),
		CreatedAt:     toProtoTimestamp(c.CreatedAt),
		UpdatedAt:     toProtoTimestamp(c.UpdatedAt),
		DeletedAt:     toProtoTimestamp(c.DeletedAt),
	}
}

// toProtoEmployee converts a domain Employee to a protobuf Employee message.
func toProtoEmployee(e models.Employee) *pb.Employee {
	return &pb.Employee{
		EmployeeId:        e.EmployeeID,
		Name:              e.Name,
		Status:            toProtoEmployeeStatus(e.Status),
		JobTitle:          toProtoEmployeeJobTitle(e.JobTitle),
		HireDate:          toProtoTimestamp(e.HireDate),
		Salary:            toProtoDecimal(e.Salary),
		LicenseIssued:     toProtoTimestamp(e.LicenseIssued),
		LicenseExpiration: toProtoTimestamp(e.LicenseExpiration),
		CreatedAt:         toProtoTimestamp(e.CreatedAt),
		UpdatedAt:         toProtoTimestamp(e.UpdatedAt),
		DeletedAt:         toProtoTimestamp(e.DeletedAt),
	}
}

// toProtoTransport converts a domain Transport to a protobuf Transport message.
func toProtoTransport(t models.Transport) *pb.Transport {
	return &pb.Transport{
		TransportId:     t.TransportID,
		Model:           t.Model,
		LicensePlate:    t.LicensePlate,
		PayloadCapacity: t.PayloadCapacity,
		FuelConsumption: t.FuelConsumption,
		CreatedAt:       toProtoTimestamp(t.CreatedAt),
		UpdatedAt:       toProtoTimestamp(t.UpdatedAt),
		DeletedAt:       toProtoTimestamp(t.DeletedAt),
	}
}

// toProtoOrder converts a domain Order to a protobuf Order message.
func toProtoOrder(o models.Order) *pb.Order {
	return &pb.Order{
		OrderId:               o.OrderID,
		ClientId:              o.ClientID,
		TransportId:           o.TransportID,
		EmployeeId:            o.EmployeeID,
		Grade:                 int32(o.Grade),
		Distance:              o.Distance,
		Weight:                o.Weight,
		TotalPrice:            toProtoDecimal(o.TotalPrice),
		PriceId:               o.PriceID,
		Status:                toProtoOrderStatus(o.Status),
		NodeIdStart:           o.NodeIDStart,
		NodeIdEnd:             o.NodeIDEnd,
		CreatedAt:             toProtoTimestamp(o.CreatedAt),
		UpdatedAt:             toProtoTimestamp(o.UpdatedAt),
		DeletedAt:             toProtoTimestamp(o.DeletedAt),
		ClientName:            o.ClientName,
		EmployeeName:          o.EmployeeName,
		TransportLicensePlate: o.TransportLicensePlate,
		PriceCargoType:        o.PriceCargoType,
		NodeStartName:         o.NodeStartName,
		NodeEndName:           o.NodeEndName,
	}
}

// toProtoPrice converts a domain Price to a protobuf Price message.
func toProtoPrice(p models.Price) *pb.Price {
	return &pb.Price{
		PriceId:   p.PriceID,
		CargoType: p.CargoType,
		Weight:    toProtoDecimal(p.Weight),
		Distance:  toProtoDecimal(p.Distance),
		CreatedAt: toProtoTimestamp(p.CreatedAt),
		UpdatedAt: toProtoTimestamp(p.UpdatedAt),
		DeletedAt: toProtoTimestamp(p.DeletedAt),
	}
}

// toProtoNode converts a domain Node to a protobuf Node message.
func toProtoNode(n models.Node) *pb.Node {
	return &pb.Node{
		NodeId: n.NodeID,
		Name:   n.Name,
		Geom: &pb.Point{
			X: n.Geom.X,
			Y: n.Geom.Y,
		},
		Address:   n.Address, // Address is assumed to be present in models.Node (if not, adjust)
		CreatedAt: toProtoTimestamp(n.CreatedAt),
		UpdatedAt: toProtoTimestamp(n.UpdatedAt),
		DeletedAt: toProtoTimestamp(n.DeletedAt),
	}
}

// toProtoInspection converts a domain Inspection to a protobuf Inspection message.
func toProtoInspection(i models.Inspection) *pb.Inspection {
	return &pb.Inspection{
		InspectionId:         i.InspectionID,
		TransportId:          i.TransportID,
		InspectionDate:       toProtoTimestamp(i.InspectionDate),
		InspectionExpiration: toProtoTimestamp(i.InspectionExpiration),
		Status:               toProtoInspectionStatus(i.Status),
		CreatedAt:            toProtoTimestamp(i.CreatedAt),
		UpdatedAt:            toProtoTimestamp(i.UpdatedAt),
		DeletedAt:            toProtoTimestamp(i.DeletedAt),
	}
}

// toProtoInsurance converts a domain Insurance to a protobuf Insurance message.
func toProtoInsurance(i models.Insurance) *pb.Insurance {
	return &pb.Insurance{
		InsuranceId:         i.InsuranceID,
		TransportId:         i.TransportID,
		InsuranceDate:       toProtoTimestamp(i.InsuranceDate),
		InsuranceExpiration: toProtoTimestamp(i.InsuranceExpiration),
		Payment:             toProtoDecimal(i.Payment),
		Coverage:            toProtoDecimal(i.Coverage),
		CreatedAt:           toProtoTimestamp(i.CreatedAt),
		UpdatedAt:           toProtoTimestamp(i.UpdatedAt),
		DeletedAt:           toProtoTimestamp(i.DeletedAt),
	}
}
