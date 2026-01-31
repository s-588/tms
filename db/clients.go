package db

import (
	"context"

	"github.com/s-588/tms/db/generated"
	"github.com/s-588/tms/models"
)

func (db DB) CreateClient(ctx context.Context, name, email, phone string, emailVerified bool) (models.Client, error) {
	arg := generated.CreateClientParams{
		Name:          name,
		Email:         email,
		Phone:         phone,
		EmailVerified: emailVerified,
	}

	genClient, err := db.queries.CreateClient(ctx, arg)
	if err != nil {
		return models.Client{}, err
	}

	return convertGeneratedClientToModel(genClient), nil
}

func (db DB) DeleteClient(ctx context.Context, clientID int) error {
	return db.queries.DeleteClient(ctx, int32(clientID))
}

func (db DB) DeleteClientOrderAssignments(ctx context.Context, clientID int) error {
	return db.queries.DeleteClientOrderAssignments(ctx, int32(clientID))
}

func (db DB) GetClientByID(ctx context.Context, clientID int) (models.Client, error) {
	genClient, err := db.queries.GetClientByclient_id(ctx, int32(clientID))
	if err != nil {
		return models.Client{}, err
	}

	return convertGeneratedClientToModel(genClient), nil
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
		order := models.Order{
			OrderID:    row.OrderID,
			Distance:   row.Distance,
			Weight:     row.Weight,
			TotalPrice: row.TotalPrice.String(),
			Status:     row.Status,
			CreatedAt:  row.CreatedAt,
			UpdatedAt:  row.UpdatedAt,
			DeletedAt:  row.DeletedAt,
		}
		orders = append(orders, order)
	}

	return orders, totalCount, nil
}

func (db DB) GetClientsPaginated(ctx context.Context, limit, offset int) ([]models.Client, int64, error) {
	arg := generated.GetClientsPaginatedParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	rows, err := db.queries.GetClientsPaginated(ctx, arg)
	if err != nil {
		return nil, 0, err
	}

	var clients []models.Client
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		client := convertGeneratedClientPaginatedRowToModel(row)
		clients = append(clients, client)
	}

	return clients, totalCount, nil
}

func (db DB) InsertClientOrderAssignments(ctx context.Context, clientID int, orderIDs []int) error {
	arg := generated.InsertClientOrderAssignmentsParams{
		ClientID: int32(clientID),
		Column2:  convertIntSliceToInt32(orderIDs),
	}

	return db.queries.InsertClientOrderAssignments(ctx, arg)
}

func (db DB) UpdateClient(ctx context.Context, clientID int, name, email, phone string) error {
	arg := generated.UpdateClientParams{
		ClientID: int32(clientID),
		Name:     name,
		Email:    email,
		Phone:    phone,
	}

	return db.queries.UpdateClient(ctx, arg)
}

func (db DB) VerifyClientEmail(ctx context.Context, emailToken string) error {
	return db.queries.VerifyClientEmail(ctx, &emailToken)
}

func convertGeneratedClientToModel(genClient generated.Client) models.Client {
	return models.Client{
		ClientID:             genClient.ClientID,
		Name:                 genClient.Name,
		Email:                genClient.Email,
		EmailVerified:        genClient.EmailVerified,
		EmailToken:           genClient.EmailToken,
		EmailTokenExpiration: genClient.EmailTokenExpiration,
		Phone:                genClient.Phone,
		CreatedAt:            genClient.CreatedAt,
		UpdatedAt:            genClient.UpdatedAt,
		DeletedAt:            genClient.DeletedAt,
	}
}

func convertGeneratedClientPaginatedRowToModel(row generated.GetClientsPaginatedRow) models.Client {
	return models.Client{
		ClientID:             row.ClientID,
		Name:                 row.Name,
		Email:                row.Email,
		EmailVerified:        row.EmailVerified,
		EmailToken:           row.EmailToken,
		EmailTokenExpiration: row.EmailTokenExpiration,
		Phone:                row.Phone,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
		DeletedAt:            row.DeletedAt,
	}
}
