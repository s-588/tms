package grpc

import (
	"context"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	pb "github.com/s-588/tms/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateNode creates a new geographical node.
func (s *Server) CreateNode(ctx context.Context, req *pb.CreateNodeRequest) (*pb.Node, error) {
	args := db.CreateNodeArgs{
		Name:    fromProtoOptionalString(req.Name),
		Address: req.Address,
		Geom: models.Point{
			X: req.Geom.X,
			Y: req.Geom.Y,
		},
	}
	n, err := s.DB.CreateNode(ctx, args)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoNode(n), nil
}

// GetNode retrieves a node by ID.
func (s *Server) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.Node, error) {
	n, err := s.DB.GetNodeByID(ctx, req.NodeId)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoNode(n), nil
}

// ListNodes returns a paginated list of nodes.
func (s *Server) ListNodes(ctx context.Context, req *pb.ListNodesRequest) (*pb.ListNodesResponse, error) {
	filter := models.NodeFilter{}
	if req.Name != nil {
		filter.Name.SetValue(*req.Name)
	}
	if req.SortBy != "" {
		filter.SortBy.SetValue(req.SortBy)
	}
	if req.SortOrder != "" {
		filter.SortOrder.SetValue(req.SortOrder)
	}

	nodes, totalPages, err := s.DB.GetNodes(ctx, req.Page, filter)
	if err != nil {
		return nil, mapError(err)
	}

	resp := &pb.ListNodesResponse{TotalPages: totalPages}
	for _, n := range nodes {
		resp.Nodes = append(resp.Nodes, toProtoNode(n))
	}
	return resp, nil
}

// UpdateNode modifies an existing node.
func (s *Server) UpdateNode(ctx context.Context, req *pb.UpdateNodeRequest) (*emptypb.Empty, error) {
	err := s.DB.UpdateNode(ctx, db.UpdateNodeArgs{
		NodeID: req.NodeId,
		Name:   fromProtoOptionalString(req.Name),
		Geom: models.Point{
			X: req.Geom.X,
			Y: req.Geom.Y,
		},
		Address: req.Address,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}

// DeleteNode soft-deletes a node by ID.
func (s *Server) DeleteNode(ctx context.Context, req *pb.DeleteNodeRequest) (*emptypb.Empty, error) {
	err := s.DB.SoftDeleteNode(ctx, req.NodeId)
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}

// CalculateDistance returns the distance in kilometers between two nodes.
func (s *Server) CalculateDistance(ctx context.Context, req *pb.CalculateDistanceRequest) (*pb.CalculateDistanceResponse, error) {
	distance, err := s.DB.CalculateDistance(ctx, req.NodeAId, req.NodeBId)
	if err != nil {
		return nil, mapError(err)
	}
	return &pb.CalculateDistanceResponse{DistanceKm: distance}, nil
}
