package ui

type contextKey string

const (
	ClientsKey    contextKey = "clients"
	EmployeesKey  contextKey = "employees"
	TransportsKey contextKey = "transports"
	PricesKey     contextKey = "prices"
	NodesKey      contextKey = "nodes"
	FormKey       contextKey = "form"
	FilterKey     contextKey = "filter"
)

type FormField struct {
	Value string
	Err   error
}

type Form map[string]FormField

type ListItem struct {
	ID int32
	Name     string
}

type List map[string][]ListItem