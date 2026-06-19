# Booking Event API

A small project offering a streamlined event booking API complete with user authentication and event management capabilities.

## What it does

- register and login users
- manage events with create, update, delete, list, and get-by-id
- protect event creation/update/delete with authentication
- keep refresh tokens in cookies
- expose Swagger docs for the API

## Tech stack

- Go 1.26
- Gin Web Framework
- PostgreSQL
- Swagger docs via swaggo
- YAML configuration

## Setup

1. Create a `config.yaml` file in the project root with your PostgreSQL and app settings:

```yaml
app:
  env: development
  port: "8080"

db:
  host: localhost
  port: "5432"
  user: your_db_user
  password: your_db_password
  name: your_db_name
```

2. Make sure PostgreSQL is running and the database exists.

   Alternatively, use Docker Compose if you prefer:

```bash
docker compose up -d
```

3. Install dependencies:

```bash
go mod tidy
```

4. Optionally, generate Swagger docs if you want the docs files updated:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/main.go -o docs
```

## Run

Use Go directly:

```bash
go run cmd/main.go
```

Or use `air` for live reload during development if you have it installed:

```bash
air
```

## Makefile commands

The project includes a simple `Makefile` for common tasks:

```bash
make migrate-up     # run database migrations
make migrate-down   # roll back the last migration
make gen-docs       # regenerate Swagger docs
```

The server starts on port `8080` by default.

## API docs

Open the Swagger UI in your browser:

```text
http://localhost:8080/api/docs/index.html
```
