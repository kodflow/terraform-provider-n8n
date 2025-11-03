# Go Plugin v2.0.1 - Advanced Go 1.23-1.25 Patterns

This plugin provides Go code standards enforcement with modern language features and performance patterns.

**v2.0.1**: Complete reference implementation with Go 1.23-1.25 advanced patterns (sync.Pool, sync.Once, sync.Map, iterators, context patterns). [See CHANGELOG](CHANGELOG.md)

## 📁 Files

### Core Documentation
- **`commands/review.md`** - Professional Go code review
- **`PACKAGE_DESCRIPTOR.md`** - Package descriptor specification (Required)
- **`GO_STANDARDS.md`** - Quick reference guide

### Reference Implementation
- **`reference-service/`** - **COMPLETE REFERENCE IMPLEMENTATION** with:
  - **Go 1.23-1.25 Advanced Patterns** (sync.Pool, sync.Once, sync.Map, iterators)
  - **Performance Optimizations** (atomic ops, bitwise flags, memory alignment)
  - **Concurrent worker pool** with goroutines
  - Channels and buffering
  - Context cancellation and graceful shutdown
  - Thread-safe concurrent access
  - 100% test coverage with race detection
  - Black-box testing with `package xxx_test`
  - All functions < 35 lines, complexity < 10
  - **15 implementation files** demonstrating ALL best practices

**👉 [See reference-service/README.md](reference-service/README.md) for detailed documentation of all patterns**

## Key Rules

### 1. Package Descriptor (NEW!)
**Every `.go` file must start with:**
```go
// Package <name> <description>
//
// Purpose:
//   <What it does>
//
// Responsibilities:
//   - <Responsibility 1>
//
// Features:
//   - <Feature 1>  (Metrics, Tracing, Database, etc.)
//
// Constraints:
//   - <Constraint 1>
//
package <name>
```

**Key points:**
- - No metrics/tracing without `Features: Metrics/Tracing`
- - Using telemetry without declaration = flagged
- - Features must be explicitly declared to be used

### 2. Code Metrics
- Functions: < 35 lines (strict)
- Complexity: < 10 (`gocyclo -over 9 .`)
- Coverage: 100% required

### 3. File Structure (Required: 1 File Per Struct)
```
package/
├── constants.go           # ALL constants
├── errors.go              # ALL errors
├── interfaces.go          # ALL interfaces
├── interfaces_test.go     # ALL mocks (package xxx_test)
├── user.go               # User struct + methods
├── user_test.go          # User tests
├── user_config.go        # UserConfig struct
├── order.go              # Order struct + methods
├── order_test.go         # Order tests
└── service.go            # Main service orchestration
```

**Key point**: Each struct must have its own dedicated file
- No `models.go` with multiple structs
- One file per struct (e.g., `user.go` for User struct)
- Better organization, clearer ownership, fewer Git conflicts

**Key point**: Test files must use `package xxx_test`:
```go
// - CORRECT
package taskqueue_test

import "taskqueue"

// - WRONG
package taskqueue  // Do NOT use same package
```

### 4. Constructor Pattern
```go
type ServiceConfig struct {
    Dep1 Interface1
    Val1 string
}

func NewService(cfg ServiceConfig) (*Service, error) {
    // Validate all required fields
    return &Service{...}, nil
}

// Avoid: &Service{...}
// Use: NewService(cfg)
```

## Available Features

Declare in Package Descriptor `Features:` section:

| Feature | Use Case | Required Import |
|---------|----------|-----------------|
| `Metrics` | OpenTelemetry metrics | `otel/metric` |
| `Tracing` | Distributed tracing | `otel/trace` |
| `Logging` | Structured logging | `log/slog` |
| `Database` | Database operations | DB interface |
| `Validation` | Input validation | Validator |
| `HTTP` | HTTP client/server | `net/http` |
| `Caching` | Cache operations | Cache interface |
| `RateLimiting` | Rate limiting | Rate limiter |
| `CircuitBreaker` | Circuit breaker | CB interface |
| `Retry` | Retry logic | Retry policy |
| `Authentication` | Auth/AuthZ | Auth provider |
| `gRPC` | gRPC services | `grpc` |
| `PubSub` | Message queues | Message broker |

**See PACKAGE_DESCRIPTOR.md for complete list**

## 📊 Quality Gates

```bash
# ALL must pass:
gocyclo -over 9 .
golangci-lint run
go vet ./...
staticcheck ./...
gosec ./...
go test -race ./...
go test -cover -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total  # Must be 100%
```

## 🚀 Usage

### Run Review
```bash
/review                  # Review changed files
/review path/to/file.go  # Review specific file
/review --full           # Full codebase review
```

### Example: Basic Service (No Telemetry)

```go
// Package userservice provides user management
//
// Purpose:
//   User CRUD operations
//
// Responsibilities:
//   - User creation and validation
//
// Features:
//   - Database
//   - Validation
//   - Logging
//
package userservice

// - NO telemetry imports - CLEAN
```

### Example: Service WITH Telemetry

```go
// Package userservice provides user management with observability
//
// Purpose:
//   User CRUD operations with full observability
//
// Responsibilities:
//   - User creation and validation
//   - Metrics collection
//
// Features:
//   - Metrics        // - Explicitly declared
//   - Tracing        // - Explicitly declared
//   - Database
//
package userservice

import (
    "go.opentelemetry.io/otel/metric"  // - OK - Metrics declared
    "go.opentelemetry.io/otel/trace"   // - OK - Tracing declared
)
```

## Common Violations

1. Missing Package Descriptor
2. Undeclared telemetry usage
3. Function > 35 lines
4. Complexity > 9
5. Coverage < 100%
6. Missing constructor
7. Wrong file structure
8. Ignored errors

## Success Checklist

Before submitting:
- [ ] Package Descriptor on every .go file
- [ ] Features explicitly declared
- [ ] No telemetry without declaration
- [ ] All functions < 35 lines
- [ ] All functions complexity < 10
- [ ] 100% test coverage
- [ ] interfaces.go with ALL interfaces
- [ ] interfaces_test.go with ALL mocks
- [ ] Every struct has NewXXXX()
- [ ] Services have XXXXConfig
- [ ] NO ignored errors
- [ ] golangci-lint passes
- [ ] gosec passes
- [ ] go test -race passes

## 📚 Reference Implementation

See `reference-service/` directory for COMPLETE, PRODUCTION-READY example:

### What It Demonstrates

**Concurrency:**
- Worker pool with multiple goroutines
- Buffered channels for task distribution
- Context-based cancellation
- Graceful shutdown with WaitGroup
- Thread-safe concurrent access (Mutex)
- Non-blocking channel operations

**Testing:**
- Black-box testing (`package xxx_test`) ✅
- Table-driven tests
- Concurrent test execution (`t.Parallel()`)
- Test helpers with `t.Helper()`
- Thread-safe mocks
- 100% code coverage
- Race detection (`go test -race`)

**Design:**
- **1 file per struct** (Required)
- Constructor with Config struct
- Dependency injection via interfaces
- Builder pattern for test data
- Repository pattern
- All functions < 35 lines
- All complexity < 10

### Structure (1:1 File per Struct)

**Implementation:**
1. **`constants.go`** - ALL constants + bitwise flags
2. **`errors.go`** - ALL error definitions
3. **`interfaces.go`** - ALL interface definitions
4. **`task.go`** - Task struct + methods
5. **`task_status.go`** - TaskStatus type + validation
6. **`task_request.go`** - CreateTaskRequest struct
7. **`task_result.go`** - TaskResult struct
8. **`worker_config.go`** - WorkerConfig struct
9. **`worker.go`** - Concurrent worker orchestration

**Tests (1:1 mapping):**
1. **`interfaces_test.go`** - ALL mocks (package taskqueue_test)
2. **`constants_test.go`** - Constants validation
3. **`errors_test.go`** - Error messages tests
4. **`task_test.go`** - Task entity tests
5. **`task_status_test.go`** - Status validation tests
6. **`task_request_test.go`** - Request validation tests
7. **`task_result_test.go`** - Result tests
8. **`worker_config_test.go`** - Config tests
9. **`worker_test.go`** - Worker integration tests

**Documentation:**
- **`README.md`** - Service documentation
- **`STRUCTURE.md`** - File organization guide

Run:
```bash
cd reference-service
go test -race -cover ./...
gocyclo -over 9 .
```

Expected:
- - 100% coverage
- - Zero race conditions
- - Zero complexity violations

## Review Process

1. **Automated Checks** - Tools run first
2. **Package Descriptor** - Verify declaration vs usage
3. **Structural Compliance** - File structure check
4. **Manual Checks** - Comprehensive review
5. **Testability** - Coverage verification
6. **Verdict** - Approved or Rejected with fixes

## 🎓 Learning Resources

- Read PACKAGE_DESCRIPTOR.md for feature system
- Study examples/ for patterns
- Check GO_STANDARDS.md for quick ref
- Review commands/review.md for full checklist

