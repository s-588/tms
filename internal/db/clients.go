package db

import (
	"context"
	"time"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
)

func (db DB) CreateClient(ctx context.Context, name, email, phone string) (models.Client, error) {
	arg := generated.CreateClientParams{
		Name:  name,
		Email: email,
		Phone: phone,
	}
	genClient, err := db.queries.CreateClient(ctx, arg)
	if err != nil {
		return models.Client{}, err
	}
	return convertGeneratedClientToModel(genClient), nil
}

func (db DB) GetClient(ctx context.Context, clientID int) (models.Client, error) {
	genClient, err := db.queries.GetClient(ctx, int32(clientID))
	if err != nil {
		return models.Client{}, err
	}
	return convertGeneratedClientToModel(genClient), nil
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
		Page:                               page,
		NameFilter:          ToStringPtr(filter.Name),
		EmailFilter:         ToStringPtr(filter.Email),
		PhoneFilter:         ToStringPtr(filter.Phone),
		EmailVerifiedFilter: ToBoolPtr(filter.EmailVerified),
		ScoreMinFilter:      ToInt16PtrFromInt(filter.ScoreMin),
		ScoreMaxFilter:      ToInt16PtrFromInt(filter.ScoreMax),
		CreatedFromFilter:   ToPgTypeTimestamptz(filter.CreatedFrom),
		CreatedToFilter:     ToPgTypeTimestamptz(filter.CreatedTo),
		UpdatedFromFilter:   ToPgTypeTimestamptz(filter.UpdatedFrom),
		UpdatedToFilter:     ToPgTypeTimestamptz(filter.UpdatedTo),
		SortOrder:           ToStringPtr(filter.SortOrder),
		SortBy:              ToStringPtr(filter.SortBy),
	}
	rows, err := db.queries.GetClients(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var clients []models.Client
	var totalPage int32
	for _, row := range rows {
		totalPage = row.TotalCount
		clients = append(clients, convertGeneratedClientRowToModel(row))
	}
	return clients, totalPage, nil
}

func (db DB) UpdateClient(ctx context.Context, clientID int,
	name, email, phone models.Optional[string],
) error {
	arg := generated.UpdateClientParams{
		ClientID: int32(clientID),
		Name:     ToStringPtr(name),
		Email:    ToStringPtr(email),
		Phone:    ToStringPtr(phone),
	}
	return db.queries.UpdateClient(ctx, arg)
}

func (db DB) SoftDeleteClient(ctx context.Context, clientID int) error {
	return db.queries.SoftDeleteClient(ctx, int32(clientID))
}

func (db DB) HardDeleteClient(ctx context.Context, clientID int) error {
	return db.queries.HardDeleteClient(ctx, int32(clientID))
}

func (db DB) RestoreClient(ctx context.Context, clientID int) error {
	return db.queries.RestoreClient(ctx, int32(clientID))
}

func (db DB) BulkSoftDeleteClients(ctx context.Context, clientIDs []int) error {
	return db.queries.BulkSoftDeleteClients(ctx, convertIntSliceToInt32(clientIDs))
}

func (db DB) BulkHardDeleteClients(ctx context.Context, clientIDs []int) error {
	return db.queries.BulkHardDeleteClients(ctx, convertIntSliceToInt32(clientIDs))
}

func (db DB) VerifyClientEmail(ctx context.Context, emailToken string) error {
	return db.queries.VerifyClientEmail(ctx, &emailToken)
}

func (db DB) SetEmailVerificationToken(ctx context.Context, clientID int, token string, expiration time.Time) error {
	arg := generated.SetEmailVerificationTokenParams{
		EmailToken:           token,
		EmailTokenExpiration: ToPgTypeTimestamptzFromTime(expiration),
		ClientID:             int32(clientID),
	}
	return db.queries.SetEmailVerificationToken(ctx, arg)
}

func (db DB) GetClientOrders(ctx context.Context, clientID, limit, offset int) ([]models.Order, int64, error) {
	arg := generated.GetClientOrdersParams{
		ClientID: int32(clientID),
		Limit:    int32(limit),
		Offset:   int32(offset),
	}
	rows, err := db.queries.GetClientOrders(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var orders []models.Order
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		orders = append(orders, convertGetClientOrdersRowToModel(row))
	}
	return orders, totalCount, nil
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

func convertGetClientOrdersRowToModel(row generated.GetClientOrdersRow) models.Order {
	return models.Order{
		OrderID:     row.OrderID,
		ClientID:    row.ClientID,
		TransportID: row.TransportID,
		EmployeeID:  row.EmployeeID,
		Grade:       uint8(row.Grade),
		Distance:    row.Distance,
		Weight:      row.Weight,
		TotalPrice:  row.TotalPrice,
		PriceID:     row.PriceID,
		Status:      models.OrderStatus(row.Status),
		NodeIDStart: row.NodeIDStart,
		NodeIDEnd:   row.NodeIDEnd,
		CreatedAt:   fromPgTimestamptz(row.CreatedAt),
		UpdatedAt:   fromPgTimestamptz(row.UpdatedAt),
		DeletedAt:   fromPgTimestamptz(row.DeletedAt),
	}
}
