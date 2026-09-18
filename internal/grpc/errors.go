package grpc

import (
	"errors"

	"github.com/s-588/tms/internal/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mapError converts domain errors to appropriate gRPC status errors.
// It preserves gRPC status errors if they are already in that format.
func mapError(err error) error {
	if err == nil {
		return nil
	}

	// Check for known domain errors.
	switch {
	case errors.Is(err, db.ErrDuplicateEmail),
		errors.Is(err, db.ErrDuplicatePhone),
		errors.Is(err, db.ErrDuplicateLicense),
		errors.Is(err, db.ErrDuplicatePrice),
		errors.Is(err, db.ErrDuplicateNodeAddress):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, db.ErrIncorrectPhone):
		return status.Error(codes.InvalidArgument, err.Error())

	default:
		// If the error is already a gRPC status, pass it through.
		if _, ok := status.FromError(err); ok {
			return err
		}
		// Otherwise treat as internal server error.
		return status.Error(codes.Internal, err.Error())
	}
}
