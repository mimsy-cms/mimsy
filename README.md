<p align="center">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="docs/imgs/header.dark.png">
      <source media="(prefers-color-scheme: light)" srcset="docs/imgs/header.png">
      <img alt="Mimsy github header" src="docs/imgs/header.png">
    </picture>
</p>

# Mimsy CMS

> A headless CMS built for agencies using SvelteKit and TypeScript

Mimsy is a modern, developer-first content management system designed specifically for agencies working with SvelteKit. It prioritizes simplicity, type safety, and developer experience while providing the flexibility needed for complex content structures.

## Overview

Mimsy addresses common pain points with existing CMS solutions by:

- **Code as Source of Truth**: All schema changes happen in TypeScript files, ensuring version control and collaboration
- **Zero Downtime Deployments**: Uses pgroll for live migrations where old and new schemas coexist
- **TypeScript Native**: Type-safe resource definitions with full IntelliSense support
- **SvelteKit First**: Built specifically for teams already using SvelteKit

## Features

### Current
- ✅ Go-based REST API with dynamic resource adaptation
- ✅ PostgreSQL database with pgroll migrations
- ✅ Marketing landing page
- ✅ Kubernetes deployment with Flux GitOps
- ✅ Development environment with Nix flakes

### In Progress
- 🚧 TypeScript SDK for resource definitions
- 🚧 Admin panel with dynamic form generation
- 🚧 Authentication and authorization

### Planned
- 📋 Content versioning and drafts
- 📋 Media management
- 📋 Webhooks and real-time updates
- 📋 Multi-language support
- 📋 Plugin system

## How it works

<p align="center">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="docs/imgs/how_it_works.dark.png">
      <source media="(prefers-color-scheme: light)" srcset="docs/imgs/how_it_works.png">
      <img alt="How it works diagram" src="docs/imgs/how_it_works.png">
    </picture>
</p>

## Architecture

Mimsy follows a modular architecture with clear separation of concerns. For detailed architecture information, see [docs/architecture.md](./docs/architecture.md).

## Project Structure

```
mimsy/
├── api/          # Go backend service
├── web/          # SvelteKit admin panel
├── landing/      # Marketing website (SvelteKit)
├── sdk/          # TypeScript SDK
├── flux/         # Kubernetes deployment configs
└── docs/         # Documentation
```

### `/api` - Backend Service
The core API built with Go that handles all data operations, authentication, and serves the admin panel. Features include:
- RESTful endpoints that adapt based on resource definitions
- PostgreSQL integration with pgroll for migrations
- Built for simplicity and performance

### `/web` - Admin Panel
A SvelteKit-based admin interface for content management:
- Dynamic form generation based on resource schemas
- Real-time content preview
- User and permission management

### `/landing` - Marketing Site
The public-facing website built with:
- SvelteKit 2 and Svelte 5
- Tailwind CSS for styling
- Prism.js for syntax highlighting
- Full test coverage with Vitest and Playwright

### `/sdk` - TypeScript SDK
The developer SDK for defining collections and content blocks:
- Type-safe resource definitions
- Compile-time validation
- Full IntelliSense support

### `/flux` - Deployment
Kubernetes manifests managed by Flux for GitOps deployments:
- Automated image updates
- Separate environments for staging and production
- SSL/TLS termination with cert-manager

## Technology Stack

- **Backend**: Go
- **Frontend**: SvelteKit, Svelte 5, TypeScript
- **Database**: PostgreSQL with pgroll
- **Styling**: Tailwind CSS
- **Testing**: Vitest, Playwright
- **Deployment**: Kubernetes, Flux
- **Development**: Nix, Docker Compose

## Getting Started

### Prerequisites

- Nix (recommended) or manually install:
  - Go 1.24.2+
  - Node.js 22+
  - pnpm
  - Docker and Docker Compose
  - pgroll
  - just (command runner)

### Using Nix (Recommended)

```bash
# Enable flakes if not already enabled
nix develop

# All dependencies are now available
# You can start a dev environment, with all required dependencies with:
dev
```

### Manual Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/mimsy.git
cd mimsy
```

2. Install dependencies for each component:
```bash
# API dependencies
cd api && go mod download

# Web dependencies
cd ../web && pnpm install

# Landing page dependencies
cd ../landing && pnpm install

# SDK dependencies
cd ../sdk && pnpm install
```

## Development

### Running the API

```bash
cd api

# Copy environment variables
cp .env.example .env

# Start database
docker compose up -d

# Initialize pgroll
pgroll init --postgres-url "postgres://mimsy:mimsy@localhost?sslmode=disable" --schema mimsy

# Run migrations
pgroll start migrations/<name>.yaml --complete --postgres-url "postgres://mimsy:mimsy@localhost?sslmode=disable" --schema mimsy

# Start the API server
go run .
```

The API will be available at `http://localhost:8080`

### Running the Admin Panel

```bash
cd web
pnpm dev
```

The admin panel will be available at `http://localhost:5174`

### Running the Landing Page

```bash
cd landing
pnpm dev
```

The landing page will be available at `http://localhost:5173`

### Running Tests

```bash
# API tests
cd api && go test ./...

# Landing page tests
cd landing && pnpm test

# E2E tests
cd landing && pnpm test:e2e
```

## Deployment

Mimsy uses Kubernetes with Flux for GitOps-based deployments. See the [flux/](./flux/) directory for deployment configurations.

Basic deployment flow:
1. Push changes to the repository
2. Flux detects changes and updates the cluster
3. Images are automatically updated based on tags

## Documentation

All the documentation is available in the [docs/](./docs/) directory.

## Contributing

We welcome contributions! Please follow these guidelines available in the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Status

Mimsy is currently in early development (v0.1.0). The core architecture is established, and we're actively working on the SDK and admin panel. We welcome early adopters and contributors who want to help shape the future of content management for SvelteKit applications.

---

Built with ❤️ for the SvelteKit community
