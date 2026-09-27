package testutil

import (
	"context"
	"testing"

	"github.com/s-588/tms/internal/config"
	"github.com/s-588/tms/internal/http"
)

func SetupHTTPServer(t *testing.T) *http.Server {
	t.Helper()
	ctx := context.Background()
	cfg := config.ServerConfig{
		HTTPPort: "0",
		HTTPS:    false,
	}
	srv := http.New(ctx, SetupTestDB(t), cfg)
	return srv
}
