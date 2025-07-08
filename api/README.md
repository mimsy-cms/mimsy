# Mimsy API

The mimsy API is a Go-based backend service that manages the lifecycle of mimsy collections. It provides endpoints for creating, updating, and deleting collections, as well as managing their items.

## Getting Started

### Prerequisites

- Go 1.24.2 or later
- Docker
- Docker Compose
- [air](https://github.com/air-verse/air)
- [pgroll](https://pgroll.com)

## Installation

Create the `.env` file based on the `.env.example` file:

```bash
cp .env.example .env
```

Start the compose services:

```bash
docker compose up -d
```

Initialize pgroll:

```bash
pgroll init --postgres-url "postgres://mimsy:mimsy@localhost?sslmode=disable" --schema mimsy
```

Run a migration:

```bash
pgroll start migrations/<name>.yaml --complete --postgres-url "postgres://mimsy:mimsy@localhost?sslmode=disable" --schema mimsy
```
