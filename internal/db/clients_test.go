package db_test

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/testutil"
)

func TestDB_Clients(t *testing.T) {
	database := testutil.SetupTestDB(t)

	testCreateClient(t, database)
	testGetClient(t, database)
}

func testCreateClient(t *testing.T, database db.DB) {
	ctx := context.Background()

	t.Run("creates client successfully", func(t *testing.T) {
		args := db.CreateClientArgs{
			Name:  "Test Tester",
			Email: "tester@test.test",
			Phone: "+1111111111",
		}

		client, err := database.CreateClient(ctx, args)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if client.ClientID == 0 {
			t.Errorf("expected client ID to be populated")
		}
		if client.Email != args.Email {
			t.Errorf("expected email %s, got %s", args.Email, client.Email)
		}
		if client.Phone != args.Phone {
			t.Errorf("expected phone %s, got %s", args.Phone, client.Phone)
		}
	})

	t.Run("error; user with duplicate email", func(t *testing.T) {
		args := db.CreateClientArgs{
			Name:  "Test Tester Original",
			Email: "duplicate_email@test.test",
			Phone: "+1211111111",
		}

		_, err := database.CreateClient(ctx, args)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		args = db.CreateClientArgs{
			Name:  "Test Tester Duplicate",
			Email: "duplicate_email@test.test",
			Phone: "+12222222222",
		}
		_, err = database.CreateClient(ctx, args)
		if !errors.Is(err, db.ErrDuplicateEmail) {
			t.Fatalf("expected \"%v\" error, got %v", db.ErrDuplicateEmail, err)
		}
	})

	t.Run("error; user with duplicate phone", func(t *testing.T) {
		args := db.CreateClientArgs{
			Name:  "Test Tester Original",
			Email: "original_email@test.test",
			Phone: "+1311111111",
		}

		_, err := database.CreateClient(ctx, args)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		args = db.CreateClientArgs{
			Name:  "Test Tester Duplicate",
			Email: "another_email@test.test",
			Phone: "+1311111111",
		}
		_, err = database.CreateClient(ctx, args)
		if !errors.Is(err, db.ErrDuplicatePhone) {
			t.Fatalf("expected \"%v\" error, got %v", db.ErrDuplicatePhone, err)
		}
	})

	t.Run("error; user with incorrect phone", func(t *testing.T) {
		args := db.CreateClientArgs{
			Name:  "Test Tester Original",
			Email: "original_email@test.test",
			Phone: "+14",
		}

		_, err := database.CreateClient(ctx, args)
		if !errors.Is(err, db.ErrIncorrectPhone) {
			t.Fatalf("expected \"%v\" error, got %v", db.ErrIncorrectPhone, err)
		}
	})
}

func testGetClient(t *testing.T, database db.DB) {
	ctx := context.Background()

	args := db.CreateClientArgs{
		Name:  "Get Test Tester",
		Email: "get@client.test",
		Phone: "+1121111111",
	}

	client, err := database.CreateClient(ctx, args)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("get client by id", func(t *testing.T) {
		c, err := database.GetClient(ctx, client.ClientID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(client, c) {
			t.Fatalf("received client different from inserted; got %v, want %v", c, client)
		}
	})

	t.Run("get client by email", func(t *testing.T) {
		c, err := database.GetClientByEmail(ctx, client.Email)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(client, c) {
			t.Fatalf("received client different from inserted; got %v, want %v", c, client)
		}
	})

	t.Run("get all clients", func(t *testing.T) {
		cs, _, err := database.GetClients(ctx, 1, models.ClientFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !slices.ContainsFunc(cs, func(c models.Client) bool {
			return reflect.DeepEqual(client, c)
		}) {
			t.Fatalf("slice doesn't contain wanted client; clients %v, want %v",
				cs, client)
		}
	})

	t.Run("get all clients; filter by email", func(t *testing.T) {
		args := db.CreateClientArgs{
			Name:  "Get Test Tester Filter",
			Email: "filter@email.test",
			Phone: "+1141111111",
		}

		client, err := database.CreateClient(ctx, args)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		cs, _, err := database.GetClients(ctx, 1, models.ClientFilter{
			Email: models.Optional[string]{
				Value: "filter@email.test",
				Set:   true,
			},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !slices.ContainsFunc(cs, func(c models.Client) bool {
			return reflect.DeepEqual(client, c)
		}) {
			t.Fatalf("slice doesn't contain wanted client; clients %v, want %v",
				cs, client)
		}
	})
}

func createTestClient(t *testing.T, database db.DB, email, phone string) models.Client {
	t.Helper()
	ctx := context.Background()

	client, err := database.CreateClient(ctx, db.CreateClientArgs{
		Name:  "Order Test Client",
		Email: email,
		Phone: phone,
	})
	if err != nil {
		t.Fatalf("create client %q: %v", email, err)
	}
	return client
}
