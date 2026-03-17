# AGENTS.md - Coding Guidelines for Agentic Agents

## Project Overview

**seckill** - A Go-based flash sale (seckill) e-commerce system with microservices architecture.

- **Module**: `seckil` (Go 1.25.7)
- **Architecture**: API Gateway (`api/`) + Backend Services (`server/`)
- **Data Layer**: MySQL + Redis + GORM
- **Service Discovery**: Consul + Nacos

## Build Commands

```bash
# Build all modules
go build ./...

# Build specific service
go build ./api/basic/cmd/
go build ./server/basic/cmd/

# Tidy dependencies
go mod tidy

# Download dependencies
go mod download
```

## Test Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestFunctionName ./path/to/package

# Run tests in specific package
go test ./server/models/

# Run tests with verbose output
go test -v ./...
```

## Lint Commands

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# With golangci-lint (if installed)
golangci-lint run
```

## Code Style Guidelines

### Package Structure
- `api/` - API Gateway and HTTP handlers
- `server/` - Backend microservices
- `middleware/` - Shared middleware components
- `proto/` - Protocol Buffer definitions (gRPC)

### Naming Conventions
- **Packages**: lowercase, single word (e.g., `models`, `config`, `router`)
- **Files**: snake_case.go (e.g., `model.go`, `config.go`)
- **Types**: PascalCase (e.g., `Member`, `Product`, `Order`)
- **Interfaces**: PascalCase with -er suffix (e.g., `Reader`, `Writer`)
- **Functions**: PascalCase for exported, camelCase for internal
- **Variables**: camelCase (e.g., `createdAt`, `totalAmount`)
- **Constants**: PascalCase or camelCase

### Imports
```go
// Group imports: stdlib -> third-party -> local
import (
    "time"
    
    "gorm.io/gorm"
    
    "seckil/server/models"
)
```

### Error Handling
- Always check errors explicitly
- Return errors up the call stack
- Use `fmt.Errorf()` with context: `fmt.Errorf("failed to create order: %w", err)`
- Log errors at appropriate level before returning if needed

### Database Models (GORM)
- Use struct tags for GORM configuration
- Add comments in Chinese for fields
- Use pointer types for nullable fields (`*time.Time`)
- Include standard fields: `CreatedAt`, `UpdatedAt`
- Define table comments with `gorm:"comment:'...'"`

Example:
```go
type Product struct {
    ID          uint      `gorm:"primaryKey;comment:'商品ID'"`
    Name        string    `gorm:"type:varchar(200);not null;index;comment:'商品名称'"`
    Price       float64   `gorm:"type:decimal(10,2);not null;comment:'价格'"`
    CreatedAt   time.Time `gorm:"comment:'创建时间'"`
    UpdatedAt   time.Time `gorm:"comment:'更新时间'"`
}
```

### Configuration
- Store in `config.yaml` at project root
- Use `config/` package for configuration structs
- Support environment-specific overrides

### Comments
- Use Chinese comments for business logic and model fields
- Use Go-style comments: `//` for single line, `/* */` for multi-line
- Document all exported types and functions

## Testing Guidelines

- Test files: `*_test.go`
- Test functions: `TestXxx(t *testing.T)`
- Table-driven tests preferred
- Use `t.Parallel()` for parallelizable tests
- Mock external dependencies (DB, Redis, HTTP)

## Common Patterns

### Service Initialization
```go
// init/init.go - Service bootstrap
// config/global.go - Global configuration
// config/config.go - Configuration loading
```

### Model Definitions
- Place in `server/models/`
- Group related models in single file
- Use Chinese comments for field documentation

## External Dependencies

Key libraries used:
- GORM - ORM framework
- Consul - Service discovery
- Nacos - Configuration management
- Redis - Caching layer
