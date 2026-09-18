package testutil

import (
	"context"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/s-588/tms/internal/config"
	"github.com/s-588/tms/internal/db"
	"github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// SetupTestDB starts a PostGIS container, runs migrations,
// and returns a ready-to-use db.DB. It registers cleanup automatically.
func SetupTestDB(t *testing.T) db.DB {
	t.Helper()
	goose.SetLogger(goose.NopLogger())
	ctx := context.Background()

	pgcontainer, err := postgres.Run(ctx, "postgis/postgis",
		postgres.WithDatabase("test_tms"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("failed to start Postgres container: %v", err)
	}

	t.Cleanup(func() {
		if err := pgcontainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	})
	log.SetDefault(log.NewNoopLogger())

	host, err := pgcontainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container's host: %v", err)
	}

	port, err := pgcontainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get container's port: %v", err)
	}

	cfg := config.DBConfig{
		Addr:     host,
		User:     "postgres",
		Password: "postgres",
		DB:       "test_tms",
		Port:     port.Port(),
	}

	db, err := db.New(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to connect to container: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}
