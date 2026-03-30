package ui

import (
	"context"

	"github.com/s-588/tms/cmd/models"
)

func GetClientsFromContext(ctx context.Context) []ListItem {
	val := ctx.Value(ClientsKey)
	if val == nil {
		return []ListItem{}
	}
	return val.([]ListItem)
}

func GetEmployeesFromContext(ctx context.Context) []ListItem {
	val := ctx.Value(EmployeesKey)
	if val == nil {
		return []ListItem{}
	}
	return val.([]ListItem)
}

func GetTransportsFromContext(ctx context.Context) []ListItem {
	val := ctx.Value(TransportsKey)
	if val == nil {
		return []ListItem{}
	}
	return val.([]ListItem)
}

func GetPricesFromContext(ctx context.Context) []ListItem {
	val := ctx.Value(PricesKey)
	if val == nil {
		return []ListItem{}
	}
	return val.([]ListItem)
}

func GetNodesFromContext(ctx context.Context) []ListItem {
	val := ctx.Value(NodesKey)
	if val == nil {
		return []ListItem{}
	}
	return val.([]ListItem)
}

func GetFormFromContext(ctx context.Context) Form {
	val := ctx.Value(FormKey)
	if val == nil {
		return Form{}
	}
	return val.(Form)
}

func GetFilterFromContext(ctx context.Context) models.OrderFilter {
	val := ctx.Value(FilterKey)
	if val == nil {
		return models.OrderFilter{}
	}
	return val.(models.OrderFilter)
}
