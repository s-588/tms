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
	v, ok := val.([]ListItem)
	if !ok {
		return []ListItem{}
	}
	return v
}

func GetEmployeesFromContext(ctx context.Context) []ListItem {
	val := ctx.Value(EmployeesKey)
	if val == nil {
		return []ListItem{}
	}
	v, ok := val.([]ListItem)
	if !ok {
		return []ListItem{}
	}
	return v
}

func GetTransportsFromContext(ctx context.Context) []ListItem {
	val := ctx.Value(TransportsKey)
	if val == nil {
		return []ListItem{}
	}
	v, ok := val.([]ListItem)
	if !ok {
		return []ListItem{}
	}
	return v
}

func GetPricesFromContext(ctx context.Context) []ListItem {
	val := ctx.Value(PricesKey)
	if val == nil {
		return []ListItem{}
	}
	v, ok := val.([]ListItem)
	if !ok {
		return []ListItem{}
	}
	return v
}

func GetNodesFromContext(ctx context.Context) []ListItem {
	val := ctx.Value(NodesKey)
	if val == nil {
		return []ListItem{}
	}
	v, ok := val.([]ListItem)
	if !ok {
		return []ListItem{}
	}
	return v
}

func GetFormFromContext(ctx context.Context) Form {
	val := ctx.Value(FormKey)
	if val == nil {
		return Form{}
	}
	v, ok := val.(Form)
	if !ok {
		return Form{}
	}
	return v
}

func GetFilterFromContext(ctx context.Context) models.OrderFilter {
	val := ctx.Value(FilterKey)
	if val == nil {
		return models.OrderFilter{}
	}
	v, ok := val.(models.OrderFilter)
	if !ok {
		return models.OrderFilter{}
	}
	return v
}
