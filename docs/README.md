[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![PostGIS](https://img.shields.io/badge/PostGIS-3.4-4169E1?style=for-the-badge&logo=postgresql)](https://postgis.net/)
[![sqlc](https://img.shields.io/badge/sqlc-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://sqlc.dev/)
[![Templ](https://img.shields.io/badge/Templ-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://templ.guide/)
[![HTMX](https://img.shields.io/badge/HTMX-3366CC?style=for-the-badge&logo=htmx&logoColor=white)](https://htmx.org/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white)](https://tailwindcss.com/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)](LICENSE)
# Transport Management System (TMS)

A logistics management system for clients, vehicles, employees, orders, and geospatial data. Built with a modern type-safe Go stack and minimal JavaScript.

## Key Features

- End-to-end type safety (sqlc + Templ)
- PostGIS distance calculations for order pricing
- Dynamic order cost based on distance, weight, fuel, and client loyalty
- Contract/act generation from DOCX templates
- Excel reports with charts
- Soft-delete + auditing on every table

## Built With

| Category              | Technology                  |
|-----------------------|-----------------------------|
| Language              | Go 1.26+                    |
| Database              | PostgreSQL 16 + PostGIS 3.4 |
| DB tooling            | sqlc, Goose                 |
| Frontend              | Templ + templUI, HTMX, Tailwind CSS |
| Containerization      | Docker + Docker Compose     |
| Task runner           | Taskfile                    |

## Getting Started

### Prerequisites

- Go 1.26+
- Docker & Docker Compose
- Task (recommended for development)

```bash
task install-tools   # installs sqlc, Templ, Goose, Tailwind, etc.
```

### Quick Start

```bash
git clone https://github.com/yourusername/tms.git
cd tms
docker compose up -d
# or: task run
```

App is available at `http://localhost:8080`.

### Development

```bash
task dev        # Templ + Tailwind watchers
task generate   # regenerate sqlc + Templ
task run        # start containers
```

Environment variables are loaded from `.env` for both the app and Postgres.

See [Database docs](db.md) and [Endpoint docs](endpoints.md) for details.

## License

Distributed under the MIT License. See [LICENSE](LICENSE).
