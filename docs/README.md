# Transport Management System (TMS)

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![PostGIS](https://img.shields.io/badge/PostGIS-3.4-4169E1?style=for-the-badge&logo=postgresql)](https://postgis.net/)
[![sqlc](https://img.shields.io/badge/sqlc-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://sqlc.dev/)
[![Goose](https://img.shields.io/badge/Goose-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://github.com/pressly/goose)
[![Templ](https://img.shields.io/badge/Templ-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://templ.guide/)
[![templUI](https://img.shields.io/badge/templUI-blue?style=for-the-badge)](https://templui.io/)
[![HTMX](https://img.shields.io/badge/HTMX-3366CC?style=for-the-badge&logo=htmx&logoColor=white)](https://htmx.org/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white)](https://tailwindcss.com/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Docker Compose](https://img.shields.io/badge/Docker_Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://docs.docker.com/compose/)
[![Taskfile](https://img.shields.io/badge/Taskfile-00ADD8?style=for-the-badge&logo=task&logoColor=white)](https://taskfile.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)](LICENSE)
---

## About The Project

TMS helps logistics companies manage clients, vehicles, employees, orders, and geospatial data. It’s designed to be simple to run and extend, using a modern Go stack with strong typing and minimal JavaScript.

Core features:
- Type-safety provided by sqlc and Templ from DB to UI.
- PostGIS calculates distances between nodes, used for order pricing.
- Order cost is computed from distance, weight, fuel consumption, and client loyalty discount.
- Generating documents of contracts and acts are rendered from DOCX templates; reports export to Excel with charts.
- Soft delete & auditing. Every table has `deleted_at` triggers, allowing safe data restoration.

### Built With

| Category            | Technology |
|---------------------|------------|
| Language            | Go 1.26+ |
| Database            | PostgreSQL 16 with PostGIS 3.4 |
| DB code generation  | sqlc |
| Migrations          | Goose |
| HTML templating     | Templ + templUI |
| Frontend interactivity | HTMX |
| Styling             | Tailwind CSS |
| Containerization    | Docker + Docker Compose |
| Task automation     | Taskfile |

---

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.26+
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)

#### For development environment

- [Task](https://taskfile.dev/)

Task's command `install-tools` can install all other dependencies. For installing TailwindCSS `npm` is required.
```bash
task install-tools
```

- [sqlc](https://sqlc.dev/)
- [Templ](https://templ.guide)
- [Goose](https://github.com/pressly/goose)
- [TailwindCSS](https://tailwindcss.com/)


### Quick Start with Docker Compose

Clone the repository

```bash
git clone https://github.com/yourusername/tms.git
cd tms
```

Start the containers.
```bash
docker compose up -d
```

Or `task run` if Task is installed.

App can be accesed on port `8080` by default.

### For development

Start Templ and Tailwind CSS servers for frontend development
```bash
task dev
```
To generate files for sqlc and Templ use `task generate`

To start containers
```bash
task run     
```

### Internal

Both app and Postgres uses .env file that contain all environment variables.

[Database docs](db.md).
[Endpoint docs](endpoints.md).

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.
