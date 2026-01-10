# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Quick Commands

### Development
- **Start database**: `docker compose up db`
- **Run server**: `go run ./cmd` (from project root)
- **Run tests**: `go test ./...`
- **Run single test**: `go test -v ./tests/controllers -run TestPokemon` (replace TestPokemon with test name)

### Database Migrations
- **Apply migrations (dev)**: `atlas migrate apply --env dev`
- **Apply migrations (test)**: `atlas migrate apply --env test`
- **Generate new migration**: `atlas migrate diff migration_name --env dev`

### Setup
1. `go mod tidy` - Download dependencies
2. Create `.env` from `.env.sample` and `.env.test` from `.env.test.sample`
3. Install atlasgo: `curl -sSf https://atlasgo.sh | sh`
4. Run migrations for your environment

## Architecture Overview

### High-Level Flow
1. **Entry point** (`cmd/main.go`) calls `Run()` in `run.go`
2. **Initialization** (`run.go`) sets up:
   - Database connection (singleton)
   - Background job queue (River workers)
   - HTTP server with middleware stack
   - Graceful shutdown handling
3. **Routes** (`routes.go`) defines all API endpoints in a single file
4. **Request flow**: Routes → Controllers → Models (GORM) → Database

### Key Components

#### Server Setup (`run.go`)
- Initializes database, workers queue, and HTTP server
- Implements graceful shutdown with 10-second timeout
- Middleware stack applied in order: errors → db → workers → auth → logging
- Database is seeded with default roles and admin user

#### Database Layer (`config/database.go` and `models/`)
- PostgreSQL with GORM ORM
- Singleton database connection using `sync.Once`
- Custom slog-based logger with slow query threshold (200ms)
- Models use a generic `DbModel` interface for CRUD operations
- All models support lifecycle hooks: `BeforeCreate`, `AfterCreate`, `BeforeUpdate`, etc. (context-aware versions)

#### Controllers (`controllers/`)
- **Generic CRUD pattern**: `GetAll[T]()`, `Get[T]()`, `Create[T]()`, `Update[T]()`, `Delete[T]()` use Go generics
- Each model must implement `DbModel` interface (includes `GetAllQuery`, `GetQuery`, `UpdateQuery`, `DeleteQuery`)
- Return `BuildResponse()` to serialize models for API responses
- Auth controllers (`Register`, `Login`) handle JWT token generation
- Custom controller (`GetRanking()`) for domain-specific queries

#### Models (`models/`)
- GORM models with validation via `go-playground/validator`
- Implement required `DbModel` interface methods for query filtering and transformation
- Filters applied per request in `GetAllQuery()` and `GetQuery()` methods
- Support pagination via query params (handled in controller)

#### Middlewares (`middlewares/`)
- **Auth** (`auth.go`): JWT verification and role-based access control (AdminOnly, UserOnly, OrganizerOnly)
- **DB** (`db.go`): Injects GORM instance into request context
- **Workers** (`workers.go`): Injects River job queue client into context
- **Errors** (`errors.go`): Catches panics and handler errors, returns JSON error responses
- **Logging** (`logging.middleware`): Logs request/response metadata
- **OnlyCurrentUser** (`only_current_user.go`): Filters queries to current user's data

#### Background Jobs (River Queue)
- **Setup** (`config/workers.go`): Two River clients - one for inserting jobs (with GORM), one for executing
- **Worker registration** (`run.go::appWorkers()`): Central place to register new workers
- **Workers** (`workers/`): Extend `river.WorkerDefaults[ArgsType]`; implement `Work(ctx, job)` method
- **Job args** (`workerargs/`): Define job payload structures

#### Request/Response Handling (`utils/`)
- `HandlerWithError`: Custom handler type that returns errors; `ErrorsMiddleware` catches them
- `Encode()`: JSON encoder that handles both success and error responses
- `ErrorsResponse`: Standard error format with `ErrorType` and `Details`
- Pagination: Built from query params; response includes `data`, `total_count`, `page`, `limit`

### Database Migrations
- **Tool**: Atlas with GORM provider (reads models from `models/` package)
- **Config**: `atlas.hcl` - defines dev, test, prod environments
- **Location**: Migrations stored in `migrations/` directory
- **Approach**: Schema-driven (models → migrations), not migration-driven

## Code Patterns

### Adding a New Model
1. Create model struct in `models/model_name.go` with GORM tags
2. Implement `DbModel` interface:
   - `GetID()`: return `m.ID`
   - Lifecycle hooks (can be empty if no custom logic needed)
   - `GetAllQuery()`: apply filters based on request
   - `GetQuery()`, `UpdateQuery()`, `DeleteQuery()`: same pattern
   - `BuildResponse()`: return struct or map for JSON serialization
3. Add validation in `Validate()` using validator tags
4. Add routes in `routes.go` using generic CRUD functions
5. Generate migration: `atlas migrate diff model_name --env dev`

### Adding a New Background Job
1. Define job args in `workerargs/job_name.go`
2. Create worker in `workers/job_name.go` extending `river.WorkerDefaults[ArgsType]`
3. Register in `run.go::appWorkers()` with `river.AddWorker()`
4. Enqueue from controller/handler: `insertClient.InsertTx(tx, jobArgs)`

### Adding a New Route
1. Add to `routes.go` with appropriate middleware chain
2. Role-based access: wrap with `m.AdminOnly()`, `m.UserOnly()`, or `m.OrganizerOnly()`
3. Error handling: wrap with `m.ErrorsMiddleware()`
4. If querying specific user's data, use `m.OnlyCurrentUserMiddleware()`

## Testing

### Test Setup
- Tests run the full server via `TestMain` in `tests/controllers/setup_server_test.go`
- Admin token automatically obtained during setup
- Database state isolated but shared across tests (use cleanup in individual tests if needed)

### Running Tests
```bash
# All tests
go test ./...

# Specific test file
go test -v ./tests/controllers -run TestPokemon

# Single test function
go test -v ./tests/controllers -run TestPokemons/Get
```

### Test Helpers
- `setup_server_test.go::TestMain()`: Server initialization and auth setup
- `helpers_test.go`: Utility functions for making HTTP requests in tests (e.g., `Get()`, `Post()`, `Put()`)

## Important Concepts

### Context Usage
- Request context available via `models.Db(r)` to get GORM instance
- Background jobs access database via context: `ctx.Value(utils.DbContextKey)` returns GORM instance
- Graceful shutdown uses signal context for clean termination

### Error Handling
- Handlers return errors to `ErrorsMiddleware` which converts to JSON
- Validation errors mapped automatically to error responses
- Panics caught and returned as 500 errors

### Authentication
- JWT tokens issued on login, stored in Authorization header
- Roles determine endpoint access; user roles loaded on token verification
- Current user ID available to filter queries per user

### Query Filtering
- Implemented per-model in `GetAllQuery()` and `GetQuery()` methods
- Common filters: role-based (admin sees all, user sees own), soft delete checks
- Pagination applied by controller after filtering

## Environment Variables
See `.env.sample` for required variables. Key ones:
- `PORT`: Server port
- `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`: Database config
- `DB_URL`: Full PostgreSQL connection string for River queue
- `ADMIN_EMAIL`, `ADMIN_PASSWORD`, `ADMIN_USERNAME`: Initial admin seed
- `APP_ENV`: Set to `prod` to skip loading .env file