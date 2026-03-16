# Blue Gopher
A simple REST API made using `net/http` and `sqlite3`.

## Running the project
Build the project from the root folder.
```bash
go build -o blue-gopher ./cmd 
./blue-gopher
```

## Architecture
The project uses Layered Architecture with a clear responsability separation and business logic isolation.

```
├── cmd
├── internal
│   ├── customerrors
│   ├── database
│   │   └── migrations
│   ├── domain
│   ├── http
│   │   ├── handlers
│   │   └── routers
│   ├── middleware
│   ├── repositories
│   └── services
├── pkg
│   └── config
└── test
```

### `cmd/`
Application entrypoint. It should contain:
- `main.go`
- Server initialization
- Dependency setup
- Manual Dependency Injection

### `internal/`
Private Application code (can't be imported by other modules).

### `internal/domain/`
Contains business entities and core domain models. This layer does not depends on HTTP or database implementations.

### `internal/repositories/`
Responsible for data access. It is responsible for interaction with SQL or an ORM with each table.

### `internal/services/`
Contains business logic. It is responsible for:
- Validations
- Business rules
- Orchestrating repository calls

### `internal/http/`
Transport layer

#### `handlers/`
- Receive HTTP requests
- Decode and validate input
- Call services
- Return JSON responses

#### `routers/`
- Route registration
- Mapping endpoint to handlers
- Middleware calls

### `internal/middleware/`
HTTP middleware such as:
- Logging
- Authentication
- CORS

### `internal/database`
Handles:
- Database initialization
- Database connection setup
- Migration execution

#### `migrations/`
Versioned `.sql` files.

### `internal/customerrors/`
Application specific errors.

### `pkg/config/`
Reusable configuration layer. Handles:
- Environment variables
- Configuration structs

### `test/`
Contains:
- Unit tests


## Request Lifecycle
```
HTTP Request
     ↓
Router
     ↓
Middleware
     ↓
Handler
     ↓
Service
     ↓
Repository
     ↓
Database
```
