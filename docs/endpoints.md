# Table specific

## Clients

### API Endpoints
 * Get client list **GET /clients?limit={limit}&offset={offset}**. Default limit is 50. Returns Templ table with users.
 * Get client **GET /clients/{id}** - returns Templ user page.
 * **GET /clients/new** - Returns Templ page with form to create a new client
 * **GET /clients/{id}/edit** - Returns Templ page with form to edit existing client
 * Add client **POST /clients** - add client.
 * Delete client **DELETE /clients/{id}** - delete user. It set clients delete_date to current time.
 * Update client **PUT /clients/{id}** - update users, field might be:
    - **name** - update user name
    - **email** - update user email and set email_verified to **false**
    - **phone** - update user phone
* Verify clients email **/verify?token={token}** - set users **email_verified** to **true**
### Junction with orders
* Get all orders **GET /clients/{id}/orders** - return Templ table with orders.
* Assign client to order **PUT /clients/{id}**. This endpoint accepts all assigned orders, this means that user with each update request must provide all necessary connections each time. If order was assigned in previous requests but it is not included in current - it will be unassigned.

## Employees

### API Endpoints
 * Get employees list **GET /employees?limit={limit}&offset={offset}**. Default limit is 50. Returns Templ table with employees.
 * Get employee **GET /employees/{id}** - returns Templ employee page.
 * **GET /employees/new** - Returns Templ page with form to create a new employee
 * **GET /employees/{id}/edit** - Returns Templ page with form to edit existing employee
 * Add employee **POST /employees** - add employee.
 * Delete employee **DELETE /employees/{id}** - delete employee. 
 * Update employee **PUT /employees/{id}?{field}** - update employees, field might be:
    - **name** - update employees name

## Fuel types

### API Endpoints
 * Get fuel types list **GET /fuels?limit={limit}&offset={offset}**. Default limit is 50. Returns Templ table with fuels.
 * Get fuel **GET /fuels/{id}** - returns Templ fuel page.
 * **GET /fuels/new** - Returns Templ page with form to create a new fuel type
 * **GET /fuels/{id}/edit** - Returns Templ page with form to edit existing fuel type
 * Add fuel **POST /fuels** - add fuel.
 * Delete fuel **DELETE /fuels/{id}** - delete fuel.
 * Update fuel **PUT /fuels/{id}?{field}** - update fuel, field might be:
    - **name** - update fuel name
    - **supplier** - update fuel supplier field
    - **cost** - update fuel cost per litter

## Orders

### API Endpoints
 * Get orders list **GET /orders?limit={limit}&offset={offset}**. Default limit is 50. Returns Templ table with orders.
 * Get order **GET /orders/{id}** - returns Templ order page.
 * **GET /orders/new** - Returns Templ page with form to create a new order
 * **GET /orders/{id}/edit** - Returns Templ page with form to edit existing order
 * Add order **POST /orders** - add order.
 * Delete order **DELETE /orders/{id}** - delete order.
 * Update order **PUT /orders/{id}?{field}** - update order, field might be:
    - **distance** - update order distance
    - **weight** - update order weight
    - **price** - update total price. **You can but you shouldn't do this because total_price is calculable field in database.**
    - **owner** - client who make this order.
### Junction with transports
* Get all transports **GET /orders/{id}/transports** - return Templ table with transports.
 * Assign or unassign transport to order **PUT /orders/{id}**. This endpoint accepts all assigned transports, this means that user with each update request must provide all necessary connections each time. If order was assigned in previous requests but it is not included in current - it will be unassigned.

## Prices

### API Endpoints
 * Get price-list **GET /prices?limit={limit}&offset={offset}**. Default limit is 50. Returns Templ table with prices.
 * Get price **GET /prices/{id}** - returns Templ price page.
 * **GET /prices/new** - Returns Templ page with form to create a new price
 * **GET /prices/{id}/edit** - Returns Templ page with form to edit existing price
 * Add price **POST /prices** - add price.
 * Delete price **DELETE /prices/{id}** - delete price.
 * Update price **PUT /price/{id}?{fields}** - update price, field might be:
    - **cargo_type** - update type of cargo
    - **cost** - update coeffiecient of total price increase.
    - **weight** - update weight coeffiecient.
    - **distance** - update distance coeffiecient.

## Transports

### API Endpoints
 * Get list of transports **GET /transports?limit={limit}&offset={offset}**. Default limit is 50. Returns Templ table with transports.
 * Get transport **GET /transports/{id}** - returns Templ transport page.
 * **GET /transports/new** - Returns Templ page with form to create a new transport
 * **GET /transports/{id}/edit** - Returns Templ page with form to edit existing transport
 * Add transport **POST /transports** - add transport.
 * Delete transport **DELETE /transports/{id}** - delete transport.
 * Update transport **PUT /transports/{id}?{fields}** - update transports, field might be:
    - **employee_id** - update the employee who assigned to this vehicle.
    - **model** - update model of the vehicle.
    - **license_plate** - update license plate of the vehicle.
    - **payload_capacity** - update payload capacity of the vehicle.
    - **fuel** - update fuel that this vehicle uses.
    - **fuel_consumption** - update fuel consuption per 100 km of this vehicle.

## Search
 * Search across all tables **GET /search/{query}?limit={limit}&offset={offset}**
