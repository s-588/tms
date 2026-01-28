package http

import "github.com/s-588/tms/http/handlers"

func (s *Server) setHandlers() {
	s.setClientEndpoints()
	s.setEmployeeEndpoints()
	s.setFuelEndpoints()
	s.setOrderEndpoints()
	s.setPriceEndpoints()
	s.setTransportEndpoints()
	s.setSearchEndpoint()
}

func (s *Server) setClientEndpoints() {
	s.mux.HandleFunc("GET /clients", handlers.GetClientsHandler)
	s.mux.HandleFunc("GET /clients/new", handlers.NewClientPageHandler)
	s.mux.HandleFunc("GET /clients/{id}/edit", handlers.EditClientPageHandler)
	s.mux.HandleFunc("GET /clients/{id}", handlers.GetClientHandler)
	s.mux.HandleFunc("POST /clients", handlers.CreateClientHandler)
	s.mux.HandleFunc("DELETE /clients/{id}", handlers.DeleteClientHandler)
	s.mux.HandleFunc("PUT /clients/{id}", handlers.UpdateClientHandler)
	s.mux.HandleFunc("GET /clients/{id}/orders", handlers.GetClientOrdersHandler)
	s.mux.HandleFunc("PUT /clients/{id}/orders", handlers.AssignClientOrdersHandler)
}

func (s *Server) setEmployeeEndpoints() {
	s.mux.HandleFunc("GET /employees", handlers.GetEmployeesHandler)
	s.mux.HandleFunc("GET /employees/{id}", handlers.GetEmployeeHandler)
	s.mux.HandleFunc("GET /employees/new", handlers.NewEmployeePageHandler)
	s.mux.HandleFunc("GET /employees/{id}/edit", handlers.EditEmployeePageHandler)
	s.mux.HandleFunc("POST /employees", handlers.CreateEmployeeHandler)
	s.mux.HandleFunc("DELETE /employees/{id}", handlers.DeleteEmployeeHandler)
	s.mux.HandleFunc("PUT /employees/{id}", handlers.UpdateEmployeeHandler)
}

func (s *Server) setFuelEndpoints() {
	s.mux.HandleFunc("GET /fuels", handlers.GetFuelsHandler)
	s.mux.HandleFunc("GET /fuels/{id}", handlers.GetFuelHandler)
	s.mux.HandleFunc("GET /fuels/new", handlers.NewFuelPageHandler)
	s.mux.HandleFunc("GET /fuels/{id}/edit", handlers.EditFuelPageHandler)
	s.mux.HandleFunc("POST /fuels", handlers.CreateFuelHandler)
	s.mux.HandleFunc("DELETE /fuels/{id}", handlers.DeleteFuelHandler)
	s.mux.HandleFunc("PUT /fuels/{id}", handlers.UpdateFuelHandler)
}

func (s *Server) setOrderEndpoints() {
	s.mux.HandleFunc("GET /orders", handlers.GetOrdersHandler)
	s.mux.HandleFunc("GET /orders/{id}", handlers.GetOrderHandler)
	s.mux.HandleFunc("GET /orders/new", handlers.NewOrderPageHandler)
	s.mux.HandleFunc("GET /orders/{id}/edit", handlers.EditOrderPageHandler)
	s.mux.HandleFunc("POST /orders", handlers.CreateOrderHandler)
	s.mux.HandleFunc("DELETE /orders/{id}", handlers.DeleteOrderHandler)
	s.mux.HandleFunc("PUT /orders/{id}", handlers.UpdateOrderHandler)
	s.mux.HandleFunc("GET /orders/{id}/transports", handlers.GetOrderTransportsHandler)
	s.mux.HandleFunc("PUT /orders/{id}/transports", handlers.AssignOrderTransportsHandler)
}

func (s *Server) setPriceEndpoints() {
	s.mux.HandleFunc("GET /prices", handlers.GetPricesHandler)
	s.mux.HandleFunc("GET /prices/{id}", handlers.GetPriceHandler)
	s.mux.HandleFunc("GET /prices/new", handlers.NewPricePageHandler)
	s.mux.HandleFunc("GET /prices/{id}/edit", handlers.EditPricePageHandler)
	s.mux.HandleFunc("POST /prices", handlers.CreatePriceHandler)
	s.mux.HandleFunc("DELETE /prices/{id}", handlers.DeletePriceHandler)
	s.mux.HandleFunc("PUT /prices/{id}", handlers.UpdatePriceHandler)
}

func (s *Server) setTransportEndpoints() {
	s.mux.HandleFunc("GET /transports", handlers.GetTransportsHandler)
	s.mux.HandleFunc("GET /transports/{id}", handlers.GetTransportHandler)
	s.mux.HandleFunc("GET /transports/new", handlers.NewTransportPageHandler)
	s.mux.HandleFunc("GET /transports/{id}/edit", handlers.EditTransportPageHandler)
	s.mux.HandleFunc("POST /transports", handlers.CreateTransportHandler)
	s.mux.HandleFunc("DELETE /transports/{id}", handlers.DeleteTransportHandler)
	s.mux.HandleFunc("PUT /transports/{id}", handlers.UpdateTransportHandler)
}

func (s *Server) setSearchEndpoint() {
	s.mux.HandleFunc("GET /search/{query}", handlers.SearchHandler)
}
