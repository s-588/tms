package grpc

import (
	"context"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	pb "github.com/s-588/tms/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateEmployee adds a new employee.
func (s *Server) CreateEmployee(ctx context.Context, req *pb.CreateEmployeeRequest) (*pb.Employee, error) {
	args := db.CreateEmployeeArgs{
		Name:              req.Name,
		Status:            fromProtoEmployeeStatus(req.Status),
		JobTitle:          fromProtoEmployeeJobTitle(req.JobTitle),
		HireDate:          fromProtoTimestamp(req.HireDate),
		Salary:            fromProtoDecimal(req.Salary),
		LicenseIssued:     fromProtoTimestamp(req.LicenseIssued),
		LicenseExpiration: fromProtoTimestamp(req.LicenseExpiration),
	}
	emp, err := s.DB.CreateEmployee(ctx, args)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoEmployee(emp), nil
}

// GetEmployee retrieves an employee by ID.
func (s *Server) GetEmployee(ctx context.Context, req *pb.GetEmployeeRequest) (*pb.Employee, error) {
	emp, err := s.DB.GetEmployeeByID(ctx, req.EmployeeId)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoEmployee(emp), nil
}

// ListEmployees returns a paginated list of employees matching the filter.
func (s *Server) ListEmployees(ctx context.Context, req *pb.ListEmployeesRequest) (*pb.ListEmployeesResponse, error) {
	filter := models.EmployeeFilter{}
	if req.Name != nil {
		filter.Name.SetValue(*req.Name)
	}
	if req.JobTitle != nil {
		filter.JobTitle.SetValue(fromProtoEmployeeJobTitle(*req.JobTitle))
	}
	if req.Status != nil {
		filter.Status.SetValue(fromProtoEmployeeStatus(*req.Status))
	}
	if req.SalaryMin != nil {
		filter.SalaryMin.SetValue(fromProtoDecimal(*req.SalaryMin))
	}
	if req.SalaryMax != nil {
		filter.SalaryMax.SetValue(fromProtoDecimal(*req.SalaryMax))
	}
	if req.SortBy != "" {
		filter.SortBy.SetValue(req.SortBy)
	}
	if req.SortOrder != "" {
		filter.SortOrder.SetValue(req.SortOrder)
	}

	employees, totalPages, err := s.DB.GetEmployees(ctx, req.Page, filter)
	if err != nil {
		return nil, mapError(err)
	}

	resp := &pb.ListEmployeesResponse{TotalPages: totalPages}
	for _, e := range employees {
		resp.Employees = append(resp.Employees, toProtoEmployee(e))
	}
	return resp, nil
}

// UpdateEmployee modifies an existing employee.
func (s *Server) UpdateEmployee(ctx context.Context, req *pb.UpdateEmployeeRequest) (*emptypb.Empty, error) {
	err := s.DB.UpdateEmployee(ctx, db.UpdateEmployeeArgs{
		EmployeeID:        req.EmployeeId,
		Name:              req.Name,
		Status:            fromProtoEmployeeStatus(req.Status),
		JobTitle:          fromProtoEmployeeJobTitle(req.JobTitle),
		HireDate:          fromProtoTimestamp(req.HireDate),
		Salary:            fromProtoDecimal(req.Salary),
		LicenseIssued:     fromProtoTimestamp(req.LicenseIssued),
		LicenseExpiration: fromProtoTimestamp(req.LicenseExpiration),
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}

// DeleteEmployee soft-deletes an employee by ID.
func (s *Server) DeleteEmployee(ctx context.Context, req *pb.DeleteEmployeeRequest) (*emptypb.Empty, error) {
	err := s.DB.SoftDeleteEmployee(ctx, req.EmployeeId)
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}
