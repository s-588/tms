package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
	"github.com/s-588/tms/internal/ui"
)

type CreateClientArgs struct {
	Name  string
	Email string
	Phone string
}

func (db DB) CreateClient(ctx context.Context, args CreateClientArgs) (models.Client, error) {
	arg := generated.CreateClientParams{
		Name:  args.Name,
		Email: args.Email,
		Phone: args.Phone,
	}
	genClient, err := db.queries.CreateClient(ctx, arg)
	if err != nil {
		return models.Client{}, parseClientError(err)
	}
	return convertGeneratedClientToModel(genClient), nil
}

func (db DB) GetClient(ctx context.Context, clientID int32) (models.Client, error) {
	genClient, err := db.queries.GetClient(ctx, clientID)
	if err != nil {
		return models.Client{}, parseClientError(err)
	}
	orders, err := db.queries.GetClientOrders(ctx, clientID)
	c := convertGeneratedClientToModel(genClient)
	for _, o := range orders {
		c.Orders = append(c.Orders, models.Order{
			OrderID:     o.OrderID,
			ClientID:    o.ClientID,
			TransportID: o.TransportID,
			EmployeeID:  o.EmployeeID,
			PriceID:     o.PriceID,
			Grade:       uint8(o.Grade),
			Distance:    o.Distance,
			Weight:      o.Weight,
			TotalPrice:  o.TotalPrice,
			Status:      models.OrderStatus(o.Status),
			NodeIDStart: o.NodeIDStart,
			NodeIDEnd:   o.NodeIDEnd,
			CreatedAt:   o.CreatedAt.Time,
			UpdatedAt:   o.UpdatedAt.Time,
			DeletedAt:   o.DeletedAt.Time,
		})
	}
	return c, nil
}

func (db DB) GetClientByEmail(ctx context.Context, email string) (models.Client, error) {
	genClient, err := db.queries.GetClientByEmail(ctx, email)
	if err != nil {
		return models.Client{}, err
	}
	return convertGeneratedClientToModel(genClient), nil
}

func (db DB) GetClients(ctx context.Context, page int32, filter models.ClientFilter) ([]models.Client, int32, error) {
	arg := generated.GetClientsParams{
		Page:                page,
		NameFilter:          ToStringPtr(filter.Name),
		EmailFilter:         ToStringPtr(filter.Email),
		PhoneFilter:         ToStringPtr(filter.Phone),
		EmailVerifiedFilter: ToBoolPtr(filter.EmailVerified),
		SortBy:              filter.GetSortBy(),
		SortOrder:           filter.GetSortOrder(),
	}
	rows, err := db.queries.GetClients(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var clients []models.Client
	var totalPages int32
	for _, row := range rows {
		totalPages = row.TotalCount
		clients = append(clients, convertGeneratedClientRowToModel(row))
	}
	return clients, totalPages, nil
}

type UpdateClientArgs struct {
	ClientID int32
	Name     string
	Email    string
	Phone    string
}

func (db DB) UpdateClient(ctx context.Context, args UpdateClientArgs) error {
	arg := generated.UpdateClientParams{
		ClientID: args.ClientID,
		Name:     args.Name,
		Email:    args.Email,
		Phone:    args.Phone,
	}
	return parseClientError(db.queries.UpdateClient(ctx, arg))
}

func (db DB) SoftDeleteClient(ctx context.Context, clientID int32) error {
	return db.queries.SoftDeleteClient(ctx, clientID)
}

func (db DB) HardDeleteClient(ctx context.Context, clientID int32) error {
	return db.queries.HardDeleteClient(ctx, clientID)
}

func (db DB) RestoreClient(ctx context.Context, clientID int32) error {
	return db.queries.RestoreClient(ctx, clientID)
}

func (db DB) BulkSoftDeleteClients(ctx context.Context, clientIDs []int32) error {
	return db.queries.BulkSoftDeleteClients(ctx, clientIDs)
}

func (db DB) BulkHardDeleteClients(ctx context.Context, clientIDs []int32) error {
	return db.queries.BulkHardDeleteClients(ctx, clientIDs)
}

func (db DB) VerifyClientEmail(ctx context.Context, emailToken string) error {
	return db.queries.VerifyClientEmail(ctx, &emailToken)
}

func (db DB) SetEmailVerificationToken(ctx context.Context, clientID int32, token string, expiration time.Time) error {
	arg := generated.SetEmailVerificationTokenParams{
		EmailToken:           token,
		EmailTokenExpiration: ToPgTypeTimestamptzFromTime(expiration),
		ClientID:             clientID,
	}
	return db.queries.SetEmailVerificationToken(ctx, arg)
}

// conversion helpers
func convertGeneratedClientToModel(c generated.Client) models.Client {
	return models.Client{
		ClientID:             c.ClientID,
		Name:                 c.Name,
		Email:                c.Email,
		EmailVerified:        c.EmailVerified,
		EmailToken:           fromStringPtr(c.EmailToken),
		EmailTokenExpiration: fromPgTimestamptz(c.EmailTokenExpiration),
		Phone:                c.Phone,
		Score:                uint8(c.Score),
		CreatedAt:            fromPgTimestamptz(c.CreatedAt),
		UpdatedAt:            fromPgTimestamptz(c.UpdatedAt),
		DeletedAt:            fromPgTimestamptz(c.DeletedAt),
	}
}

func convertGeneratedClientRowToModel(row generated.GetClientsRow) models.Client {
	return models.Client{
		ClientID:             row.ClientID,
		Name:                 row.Name,
		Email:                row.Email,
		EmailVerified:        row.EmailVerified,
		EmailToken:           fromStringPtr(row.EmailToken),
		EmailTokenExpiration: fromPgTimestamptz(row.EmailTokenExpiration),
		Phone:                row.Phone,
		Score:                uint8(row.Score),
		CreatedAt:            fromPgTimestamptz(row.CreatedAt),
		UpdatedAt:            fromPgTimestamptz(row.UpdatedAt),
		DeletedAt:            fromPgTimestamptz(row.DeletedAt),
	}
}

func (db DB) ListClients(ctx context.Context) ([]ui.ListItem, error) {
	rows, err := db.queries.ListClients(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]ui.ListItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, ui.ListItem{
			ID:   r.ClientID,
			Name: r.Name,
		})
	}
	return items, nil
}

func parseClientError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case "23505": // unique constaint violation
			switch pgErr.ConstraintName {
			case "clients_email_key":
				return ErrDuplicateEmail
			case "clients_phone_key":
				return ErrDuplicatePhone
			case "23514": // check constaint violation
				switch pgErr.ConstraintName {
				case "clients_phone_check":
					return ErrIncorrectPhone
				}
			}
			return fmt.Errorf("uknown error: %w", err)
		}
	}
	return err
}

func (db DB) CountClientsOrders(ctx context.Context, id int32) (total int64, canceled int64, err error) {
	rows, err := db.queries.CountClientOrders(ctx, id)
	if err != nil {
		return 0, 0, parseClientError(err)
	}
	return rows.TotalOrders, rows.CancelledOrders, nil
}
