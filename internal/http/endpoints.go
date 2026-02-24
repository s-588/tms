package http

func (s *Server) setHandlers() {
	s.setClientEndpoints()
	s.setEmployeeEndpoints()
	s.setOrderEndpoints()
	s.setPriceEndpoints()
	s.setTransportEndpoints()
	s.setSearchEndpoint()
}

func (s *Server) setClientEndpoints() {
	s.mux.HandleFunc("GET /clients", s.Handler.GetClientsPage)
	s.mux.HandleFunc("GET /clients/tage", s.Handler.GetClients)
	s.mux.HandleFunc("GET /clients/new", s.Handler.NewClientPageHandler)
	s.mux.HandleFunc("GET /clients/{id}/edit", s.Handler.EditClientPageHandler)
	s.mux.HandleFunc("GET /clients/{id}", s.Handler.GetClientHandler)
	s.mux.HandleFunc("POST /clients", s.Handler.CreateClientHandler)
	s.mux.HandleFunc("DELETE /clients/{id}", s.Handler.DeleteClient)
	s.mux.HandleFunc("PUT /clients/{id}", s.Handler.UpdateClient)
	s.mux.HandleFunc("GET /clients/{id}/orders", s.Handler.GetClientOrders)
}

func (s *Server) setEmployeeEndpoints() {
	s.mux.HandleFunc("GET /employees", s.Handler.GetEmployeesHandler)
	// s.mux.HandleFunc("GET /employees/{id}", s.Handler.GetEmployeeHandler)
	s.mux.HandleFunc("GET /employees/new", s.Handler.NewEmployeePageHandler)
	s.mux.HandleFunc("GET /employees/{id}/edit", s.Handler.EditEmployeePageHandler)
	s.mux.HandleFunc("POST /employees", s.Handler.CreateEmployeeHandler)
	s.mux.HandleFunc("DELETE /employees/{id}", s.Handler.DeleteEmployeeHandler)
	s.mux.HandleFunc("PUT /employees/{id}", s.Handler.UpdateEmployeeHandler)
}

func (s *Server) setOrderEndpoints() {
	s.mux.HandleFunc("GET /orders", s.Handler.GetOrdersHandler)
	s.mux.HandleFunc("GET /orders/{id}", s.Handler.GetOrderHandler)
	s.mux.HandleFunc("GET /orders/new", s.Handler.NewOrderPageHandler)
	s.mux.HandleFunc("GET /orders/{id}/edit", s.Handler.EditOrderPageHandler)
	s.mux.HandleFunc("POST /orders", s.Handler.CreateOrderHandler)
	s.mux.HandleFunc("DELETE /orders/{id}", s.Handler.DeleteOrderHandler)
	s.mux.HandleFunc("PUT /orders/{id}", s.Handler.UpdateOrderHandler)
	s.mux.HandleFunc("GET /orders/{id}/transports", s.Handler.GetOrderTransportsHandler)
	s.mux.HandleFunc("PUT /orders/{id}/transports", s.Handler.AssignOrderTransportsHandler)
}

func (s *Server) setPriceEndpoints() {
	s.mux.HandleFunc("GET /prices", s.Handler.GetPricesHandler)
	s.mux.HandleFunc("GET /prices/{id}", s.Handler.GetPriceHandler)
	s.mux.HandleFunc("GET /prices/new", s.Handler.NewPricePageHandler)
	s.mux.HandleFunc("GET /prices/{id}/edit", s.Handler.EditPricePageHandler)
	s.mux.HandleFunc("POST /prices", s.Handler.CreatePriceHandler)
	s.mux.HandleFunc("DELETE /prices/{id}", s.Handler.DeletePriceHandler)
	s.mux.HandleFunc("PUT /prices/{id}", s.Handler.UpdatePriceHandler)
}

func (s *Server) setTransportEndpoints() {
	s.mux.HandleFunc("GET /transports", s.Handler.GetTransportsHandler)
	s.mux.HandleFunc("GET /transports/{id}", s.Handler.GetTransportHandler)
	s.mux.HandleFunc("GET /transports/new", s.Handler.NewTransportPageHandler)
	s.mux.HandleFunc("GET /transports/{id}/edit", s.Handler.EditTransportPageHandler)
	s.mux.HandleFunc("POST /transports", s.Handler.CreateTransportHandler)
	s.mux.HandleFunc("DELETE /transports/{id}", s.Handler.DeleteTransportHandler)
	s.mux.HandleFunc("PUT /transports/{id}", s.Handler.UpdateTransportHandler)
}

func (s *Server) setSearchEndpoint() {
	s.mux.HandleFunc("GET /search/{query}", s.Handler.SearchHandler)
}
