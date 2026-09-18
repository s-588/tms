TMS exposes an HTML/HTMX API built with Templ. Most list endpoints return table partials for HTMX; page endpoints return full pages. The gRPC API is defined separately in `tms.proto`.

## Conventions

- Pagination uses `page` query parameter. Default is `1`.
- Most list endpoints accept `sort` and `order` query parameters.
  - `sort` is the column name.
  - `order` is usually `asc` or `desc`; if omitted, most resources default to `desc`.
- Clients are the exception: they accept only `sort`; if `sort` is present, order defaults to `desc`; otherwise `created_at asc`.
- Deletes are soft deletes unless noted otherwise.
- Bulk deletes use `POST /{resource}/delete` with repeated `selected_ids` form values.
- Forms render validation errors back into the same Templ form.
- Filters are parsed from URL query parameters.

---

## Clients

### Routes

| Method | Path | Handler | Description |
|---|---|---|---|
| GET | `/clients/page` | `GetClientsPage` | Full clients page |
| GET | `/clients` | `GetClients` | Clients table partial |
| GET | `/clients/{id}` | `GetClientHandler` | Client detail sheet |
| POST | `/clients` | `CreateClientHandler` | Create client |
| PUT | `/clients/{id}` | `UpdateClient` | Update client |
| DELETE | `/clients/{id}` | `DeleteClient` | Soft-delete client |
| POST | `/clients/delete` | `BulkDeleteClients` | Bulk soft-delete clients |
| GET | `/clients/{id}/orders` | `GetClientOrders` | Placeholder; currently not implemented |

### Filters

`GET /clients`

- `name`
- `email`
- `phone`
- `email_verified` — `true` or `false`
- `sort` — column name
- `page`

### Form Fields

`POST /clients`, `PUT /clients/{id}`

- `name`
- `email`
- `phone`

Validation:

- Name must be longer than 3 runes.
- Email must be a valid email format.
- Phone is parsed for Belarus (`BY`), normalized, and must be valid.
- Duplicate email and phone are rejected.
- Updating email resets `email_verified` to `false`.

### Bulk Delete

`POST /clients/delete`

Form:

- `selected_ids` — repeated client IDs

---

## Employees

### Routes

| Method | Path | Handler | Description |
|---|---|---|---|
| GET | `/employees/page` | `GetEmployeesPage` | Full employees page |
| GET | `/employees` | `GetEmployees` | Employees table partial |
| GET | `/employees/{id}` | `GetEmployeeHandler` | Employee detail sheet |
| POST | `/employees` | `CreateEmployeeHandler` | Create employee |
| PUT | `/employees/{id}` | `UpdateEmployeeHandler` | Update employee |
| DELETE | `/employees/{id}` | `DeleteEmployeeHandler` | Soft-delete employee |
| POST | `/employees/delete` | `BulkDeleteEmployeesHandler` | Bulk soft-delete employees |

### Filters

`GET /employees`

- `name`
- `job_title` — `driver`, `dispatcher`, `mechanic`, `logistics_manager`
- `status` — `available`, `assigned`, `unavailable`
- `salary_min`
- `salary_max`
- `sort` — default `employee_id`
- `order` — default `desc`
- `page`

### Form Fields

`POST /employees`, `PUT /employees/{id}`

- `name`
- `status`
- `job_title`
- `hire_date` — `YYYY-MM-DD`
- `salary`
- `license_issued` — `YYYY-MM-DD`
- `license_expiration` — `YYYY-MM-DD`

Validation:

- Name must be at least 2 characters.
- Status and job title must be valid.
- Salary must be a positive decimal.
- License expiration must be after license issue date.
- An employee with an expired license cannot be assigned.

### Bulk Delete

`POST /employees/delete`

Form:

- `selected_ids` — repeated employee IDs

---

## Orders

### Routes

| Method | Path | Handler | Description |
|---|---|---|---|
| GET | `/orders/page` | `GetOrdersPage` | Full orders page |
| GET | `/orders` | `GetOrders` | Orders table partial |
| GET | `/orders/{id}` | `GetOrderHandler` | Order detail sheet |
| POST | `/orders` | `CreateOrderHandler` | Create order |
| PUT | `/orders/{id}` | `UpdateOrderHandler` | Update order |
| DELETE | `/orders/{id}` | `DeleteOrderHandler` | Soft-delete order |
| POST | `/orders/delete` | `BulkDeleteOrdersHandler` | Bulk soft-delete orders |
| GET | `/orders/{id}/transports` | `GetOrderTransportsHandler` | Placeholder; currently not implemented |
| PUT | `/orders/{id}/transports` | `AssignOrderTransportsHandler` | Placeholder; currently not implemented |
| GET | `/orders/export` | `OrdersExport` | Export orders report to XLSX |
| GET | `/orders/{id}/contract` | `DownloadContract` | Download DOCX contract |
| GET | `/orders/{id}/certificate` | `DownloadAct` | Download DOCX act/certificate |

### Filters

`GET /orders`

- `client_id`
- `transport_id`
- `employee_id`
- `price_id`
- `distance_min`
- `distance_max`
- `weight_min`
- `weight_max`
- `price_min`
- `price_max`
- `grade_min` — `0..5`
- `grade_max` — `0..5`
- `status` — `pending`, `assigned`, `in_progress`, `completed`, `cancelled`
- `sort` — default `order_id`
- `order` — default `desc`
- `page`

### Create Form

`POST /orders`

Fields:

- `client_id`
- `transport_id`
- `employee_id`
- `price_id`
- `weight`
- `status`
- `node_start`
- `node_end`

Server-side logic:

- Calculates distance between `node_start` and `node_end`.
- Fetches transport and checks `payload_capacity >= weight`.
- Fetches employee and requires:
  - `status == available`
  - `job_title == driver`
- Calculates total price using the TMS pricing engine.
- Sets `grade = 0`.

### Update Form

`PUT /orders/{id}`

Fields:

- `client_id`
- `transport_id`
- `employee_id`
- `price_id`
- `weight`
- `status`
- `node_start`
- `node_end`

Current HTTP update behavior:

- Preserves existing `grade`.
- Preserves existing `distance`.
- Preserves existing `total_price`.
- Does not recalculate distance or price in the HTTP handler.

### Export

`GET /orders/export`

Query parameters:

- Same filters as `GET /orders`
- `start` — `YYYY-MM-DD`
- `end` — `YYYY-MM-DD`

Returns:

- `orders_report.xlsx`

Report includes:

- Summary totals.
- Orders table.
- Orders by month.
- Orders by status.
- Top clients by revenue.

### Documents

`GET /orders/{id}/contract`

Returns a DOCX contract.

`GET /orders/{id}/certificate`

Returns a DOCX act/certificate.

---

## Prices

### Routes

| Method | Path | Handler | Description |
|---|---|---|---|
| GET | `/prices/page` | `GetPricesPage` | Full prices page |
| GET | `/prices` | `GetPrices` | Prices table partial |
| GET | `/prices/{id}` | `GetPriceHandler` | Price detail sheet |
| GET | `/prices/new` | `NewPricePageHandler` | New price form |
| POST | `/prices` | `CreatePriceHandler` | Create price |
| PUT | `/prices/{id}` | `UpdatePriceHandler` | Update price |
| DELETE | `/prices/{id}` | `DeletePriceHandler` | Soft-delete price |
| POST | `/prices/delete` | `BulkDeletePricesHandler` | Bulk soft-delete prices |

### Filters

`GET /prices`

- `cargo_type`
- `weight_min`
- `weight_max`
- `distance_min`
- `distance_max`
- `sort` — default `price_id`
- `order` — default `desc`
- `page`

### Form Fields

`POST /prices`, `PUT /prices/{id}`

- `cargo_type`
- `weight` — positive decimal
- `distance` — positive decimal

Duplicate price configurations are rejected.

### Bulk Delete

`POST /prices/delete`

Form:

- `selected_ids` — repeated price IDs

---

## Transports

### Routes

| Method | Path | Handler | Description |
|---|---|---|---|
| GET | `/transports/page` | `GetTransportsPage` | Full transports page |
| GET | `/transports` | `GetTransports` | Transports table partial |
| GET | `/transports/{id}` | `GetTransportHandler` | Transport detail sheet |
| GET | `/transports/new` | `NewTransportPageHandler` | New transport form |
| GET | `/transports/{id}/edit` | `EditTransportPageHandler` | Edit transport form |
| POST | `/transports` | `CreateTransportHandler` | Create transport |
| PUT | `/transports/{id}` | `UpdateTransportHandler` | Update transport |
| DELETE | `/transports/{id}` | `DeleteTransportHandler` | Soft-delete transport |

Bulk delete is currently commented out in the router.

### Filters

`GET /transports`

- `model`
- `license_plate`
- `payload_min`
- `payload_max`
- `fuel_min`
- `fuel_max`
- `sort` — default `transport_id`
- `order` — default `desc`
- `page`

### Form Fields

`POST /transports`, `PUT /transports/{id}`

- `model`
- `license_plate`
- `payload_capacity` — positive integer
- `fuel_consumption` — positive integer

Duplicate license plates are rejected.

---

## Inspections

### Routes

| Method | Path | Handler | Description |
|---|---|---|---|
| GET | `/inspections/page` | `GetInspectionsPage` | Full inspections page |
| GET | `/inspections` | `GetInspections` | Inspections table partial |
| GET | `/inspections/{id}` | `GetInspectionHandler` | Inspection detail sheet |
| POST | `/inspections` | `CreateInspectionHandler` | Create inspection |
| PUT | `/inspections/{id}` | `UpdateInspectionHandler` | Update inspection |
| DELETE | `/inspections/{id}` | `DeleteInspectionHandler` | Soft-delete inspection |
| POST | `/inspections/delete` | `BulkDeleteInspectionsHandler` | Bulk soft-delete inspections |

### Filters

`GET /inspections`

- `transport_id`
- `status` — `ready`, `repair`, `overdue`
- `inspection_from` — `YYYY-MM-DD`
- `inspection_to` — `YYYY-MM-DD`
- `expiration_from` — `YYYY-MM-DD`
- `expiration_to` — `YYYY-MM-DD`
- `sort` — default `inspection_id`
- `order` — default `desc`
- `page`

### Form Fields

`POST /inspections`, `PUT /inspections/{id}`

- `transport_id`
- `status`
- `inspection_date` — `YYYY-MM-DD`
- `inspection_expiration` — `YYYY-MM-DD`

Expiration must be after inspection date.

### Bulk Delete

`POST /inspections/delete`

Form:

- `selected_ids` — repeated inspection IDs

---

## Insurances

### Routes

| Method | Path | Handler | Description |
|---|---|---|---|
| GET | `/insurances/page` | `GetInsurancesPage` | Full insurances page |
| GET | `/insurances` | `GetInsurances` | Insurances table partial |
| GET | `/insurances/{id}` | `GetInsuranceHandler` | Insurance detail sheet |
| POST | `/insurances` | `CreateInsuranceHandler` | Create insurance |
| PUT | `/insurances/{id}` | `UpdateInsuranceHandler` | Update insurance |
| DELETE | `/insurances/{id}` | `DeleteInsuranceHandler` | Soft-delete insurance |
| POST | `/insurances/delete` | `BulkDeleteInsurancesHandler` | Bulk soft-delete insurances |

### Filters

`GET /insurances`

- `transport_id`
- `insurance_from` — `YYYY-MM-DD`
- `insurance_to` — `YYYY-MM-DD`
- `expiration_from` — `YYYY-MM-DD`
- `expiration_to` — `YYYY-MM-DD`
- `payment_min`
- `payment_max`
- `coverage_min`
- `coverage_max`
- `sort` — default `insurance_id`
- `order` — default `desc`
- `page`

### Form Fields

`POST /insurances`, `PUT /insurances/{id}`

- `transport_id`
- `insurance_date` — `YYYY-MM-DD`
- `insurance_expiration` — `YYYY-MM-DD`
- `payment` — positive decimal
- `coverage` — positive decimal

Expiration must be after insurance date.

### Bulk Delete

`POST /insurances/delete`

Form:

- `selected_ids` — repeated insurance IDs

---

## Nodes

### Routes

| Method | Path | Handler | Description |
|---|---|---|---|
| GET | `/nodes/page` | `GetNodesPage` | Full nodes page |
| GET | `/nodes` | `GetNodes` | Nodes table partial |
| GET | `/nodes/{id}` | `GetNodeHandler` | Node detail sheet |
| POST | `/nodes` | `CreateNodeHandler` | Create node |
| PUT | `/nodes/{id}` | `UpdateNodeHandler` | Update node |
| DELETE | `/nodes/{id}` | `DeleteNodeHandler` | Soft-delete node |
| POST | `/nodes/delete` | `BulkDeleteNodesHandler` | Bulk soft-delete nodes |

### Filters

`GET /nodes`

- `name`
- `sort` — default `node_id`
- `order` — default `desc`
- `page`

### Form Fields

`POST /nodes`, `PUT /nodes/{id}`

- `name` — optional on create, required on update
- `address` — required
- `x` — float coordinate
- `y` — float coordinate

Duplicate node addresses are rejected.

### Bulk Delete

`POST /nodes/delete`

Form:

- `selected_ids` — repeated node IDs

---

## Search

### Route

| Method | Path | Handler | Description |
|---|---|---|---|
| GET | `/search/{query}` | `SearchHandler` | Placeholder; currently not implemented |

---

## Static and Index

| Method | Path | Description |
|---|---|---|
| GET | `/` | Index page |
| GET | `/static/...` | Static files |

---

## gRPC

The gRPC API is defined in `tms.proto`. It exposes services for:

- `ClientService`
- `EmployeeService`
- `TransportService`
- `OrderService`
- `PriceService`
- `NodeService`
- `InspectionService`
- `InsuranceService`

See `tms.proto` for request/response messages, filters, pagination, and enum values.
