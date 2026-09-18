package grpc

import (
	"context"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	pb "github.com/s-588/tms/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateClient creates a new client from the request.
func (s *Server) CreateClient(ctx context.Context, req *pb.CreateClientRequest) (*pb.Client, error) {
	args := db.CreateClientArgs{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}
	client, err := s.DB.CreateClient(ctx, args)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoClient(client), nil
}

// GetClient retrieves a client by ID.
func (s *Server) GetClient(ctx context.Context, req *pb.GetClientRequest) (*pb.Client, error) {
	client, err := s.DB.GetClient(ctx, req.ClientId)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoClient(client), nil
}

// ListClients returns a paginated list of clients matching the filter.
func (s *Server) ListClients(ctx context.Context, req *pb.ListClientsRequest) (*pb.ListClientsResponse, error) {
	filter := models.ClientFilter{}
	if req.Name != nil {
		filter.Name.SetValue(*req.Name)
	}
	if req.Email != nil {
		filter.Email.SetValue(*req.Email)
	}
	if req.Phone != nil {
		filter.Phone.SetValue(*req.Phone)
	}
	if req.EmailVerified != nil {
		filter.EmailVerified.SetValue(*req.EmailVerified)
	}
	if req.SortBy != "" {
		filter.SortBy.SetValue(req.SortBy)
	}
	if req.SortOrder != "" {
		filter.SortOrder.SetValue(req.SortOrder)
	}

	clients, totalPages, err := s.DB.GetClients(ctx, req.Page, filter)
	if err != nil {
		return nil, mapError(err)
	}

	resp := &pb.ListClientsResponse{TotalPages: totalPages}
	for _, c := range clients {
		resp.Clients = append(resp.Clients, toProtoClient(c))
	}
	return resp, nil
}

// UpdateClient updates an existing client's information.
func (s *Server) UpdateClient(ctx context.Context, req *pb.UpdateClientRequest) (*emptypb.Empty, error) {
	err := s.DB.UpdateClient(ctx, db.UpdateClientArgs{
		ClientID: req.ClientId,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}

// DeleteClient soft-deletes a client by ID.
func (s *Server) DeleteClient(ctx context.Context, req *pb.DeleteClientRequest) (*emptypb.Empty, error) {
	err := s.DB.SoftDeleteClient(ctx, req.ClientId)
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}
