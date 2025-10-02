# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go-based REST API that fetches book cover images from Goodreads. Supports searching by book title/author or ISBN-13. Uses memcached for caching and is deployed via Kubernetes with Helm charts.

## Common Commands

### Development
```bash
# Run the application locally
go run cmd/api/main.go

# Run tests
go test ./internal/handler/...

# Build the binary
go build -o bin/server ./cmd/api
```

### Docker
```bash
# Build and run with Docker Compose
docker compose up --build

# Build Docker image manually
docker build -t bookcover-api .
```

### Testing
```bash
# Run all tests (handler, scraper, service)
make test

# Run tests with verbose output
make test-verbose

# Run with coverage
make test-coverage

# Generate HTML coverage report
make test-coverage-html

# Run a specific test
go test -v ./internal/handler -run TestBookcoverSearch_CacheHit
```

## Architecture

### Project Structure
```
bookcover-api/
├── cmd/api/              # Application entry point
│   └── main.go
├── internal/             # Private application code
│   ├── cache/           # Memcached client wrapper
│   │   ├── cache.go
│   │   └── cache_test.go (no tests yet)
│   ├── config/          # Configuration constants
│   │   └── config.go
│   ├── handler/         # HTTP request handlers
│   │   ├── home.go
│   │   ├── bookcover.go
│   │   └── bookcover_test.go (95.7% coverage)
│   ├── middleware/      # HTTP middleware
│   │   └── middleware.go
│   ├── scraper/         # Web scraping logic (Goodreads)
│   │   ├── goodreads.go
│   │   └── goodreads_test.go (52.1% coverage)
│   ├── service/         # Business logic layer
│   │   ├── bookcover.go
│   │   └── bookcover_test.go (91.9% coverage)
│   └── server/          # Server setup and routing
│       └── server.go
├── pkg/                 # Potentially reusable packages
│   └── response/        # HTTP response helpers
│       └── response.go
├── mocks/               # Test mocks
│   └── cache.go
├── helm/                # Kubernetes Helm charts
├── Dockerfile
├── compose.yaml         # Docker Compose with memcached
└── Makefile
```

### Key Components

**Application Flow:**
```
HTTP Request → Server → Middleware → Handler → Service → Scraper
                                                    ↓
                                                  Cache
```

**Entry Point (cmd/api/main.go)**
- Minimal main function that calls `server.Start()`

**Server Setup (internal/server/server.go)**
- Initializes all dependencies (cache, scraper, service, handlers)
- Configures routes with middleware chains
- Three endpoints:
  - `GET /` - Home page (returns 404)
  - `GET /bookcover?book_title=X&author_name=Y` - Search by title/author
  - `GET /bookcover/{isbn}` - Search by ISBN-13

**Middleware (internal/middleware/middleware.go)**
- Custom middleware system using function composition
- Middlewares applied in reverse order via `Chain()` function
- Available: HttpMethod, JsonHeaderMiddleware, CorsHeaderMiddleware

**Handlers (internal/handler/)**
- `BookcoverHandler`: Processes HTTP requests, validates input
- Delegates business logic to service layer
- Uses `pkg/response` for consistent JSON responses
- Tests live alongside implementation (`bookcover_test.go`)

**Service Layer (internal/service/bookcover.go)**
- `BookcoverService`: Core business logic
- Handles caching strategy (check cache → fetch → cache result)
- Abstracts scraper implementation via interface
- Normalizes input (spaces to `+`, lowercase for cache keys)

**Scraper (internal/scraper/goodreads.go)**
- `Goodreads` struct: Web scraping implementation
- `FetchByTitleAuthor`: Searches and parses Goodreads HTML
- `FetchByISBN`: Direct ISBN lookup
- Uses goquery for HTML parsing
- Regex to extract high-quality image URLs

**Cache Layer (internal/cache/cache.go)**
- Singleton pattern for memcached client
- `CacheClient` interface for testability
- Mock implementation in `mocks/cache.go`
- Environment variable: `MEMCACHED_HOST` (defaults to port 11211)

### Testing Strategy
- Tests live next to implementation (`*_test.go` files)
- **Overall coverage: 75.0%**
  - Handler: 95.7% ✅
  - Service: 91.9% ✅
  - Scraper: 52.1% (HTTP methods not tested, parsing logic 85-92%)
- Uses `httptest` for HTTP request/response mocking
- Mock scraper and cache from `mocks/` package
- Test coverage: cache hits, cache misses, validation, error cases, HTML parsing

### Deployment
- **Docker Compose**: Runs API + memcached locally (`docker compose up`)
- **Jenkins CI/CD**: Runs tests on commits, builds Docker images
- **Docker**: Multi-stage build (Go 1.23 build stage, Alpine runtime)
- **Kubernetes**: Helm chart deploys API and memcached pods
- **ArgoCD**: Handles deployment automation
- **Harbor**: Container registry for images

## Environment Variables
- `MEMCACHED_HOST` - Memcached server hostname (required for caching)
- `GOOGLE_BOOKS_API_KEY` - Defined as constant but not currently used in code

## Design Decisions

### Why `internal/` directory?
The `internal/` directory is a Go convention that prevents external projects from importing these packages. Since this is a standalone API service (not a reusable library), using `internal/` is the correct architectural choice.

### Layered Architecture
The project follows a clean layered architecture:
- **Handler layer**: HTTP request/response handling, validation
- **Service layer**: Business logic, caching strategy, orchestration
- **Scraper layer**: External data fetching (Goodreads web scraping)
- **Infrastructure**: Cache, middleware, configuration

Each layer depends only on abstractions (interfaces) from lower layers, enabling:
- Easy testing with mocks
- Swappable implementations (e.g., different scrapers)
- Clear separation of concerns
