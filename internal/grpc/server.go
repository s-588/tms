package grpc

import (
	"fmt"
	"net"

	"github.com/s-588/tms/internal/db"
	pb "github.com/s-588/tms/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server combines all gRPC service implementations and holds a database handle.
type Server struct {
	pb.UnimplementedClientServiceServer
	pb.UnimplementedEmployeeServiceServer
	pb.UnimplementedOrderServiceServer
	pb.UnimplementedTransportServiceServer
	pb.UnimplementedPriceServiceServer
	pb.UnimplementedNodeServiceServer
	pb.UnimplementedInspectionServiceServer
	pb.UnimplementedInsuranceServiceServer
	DB db.DB
}

// NewServer creates a new gRPC server with the given database connection.
func NewServer(db db.DB) *Server {
	return &Server{DB: db}
}

// RegisterServices registers all service implementations with the provided gRPC server.
// It also enables reflection for debugging purposes.
func (s *Server) RegisterServices(grpcServer *grpc.Server) {
	pb.RegisterClientServiceServer(grpcServer, s)
	pb.RegisterEmployeeServiceServer(grpcServer, s)
	pb.RegisterOrderServiceServer(grpcServer, s)
	pb.RegisterTransportServiceServer(grpcServer, s)
	pb.RegisterPriceServiceServer(grpcServer, s)
	pb.RegisterNodeServiceServer(grpcServer, s)
	pb.RegisterInspectionServiceServer(grpcServer, s)
	pb.RegisterInsuranceServiceServer(grpcServer, s)

	reflection.Register(grpcServer)
}

// Run starts the gRPC server on the specified port.
// It blocks until the server stops or an error occurs.
func (s *Server) Run(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	grpcServer := grpc.NewServer()
	s.RegisterServices(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("gRPC server error: %w", err)
	}
	return nil
}
