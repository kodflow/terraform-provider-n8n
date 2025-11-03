# Go Code Standards - Quick Reference

## Core Rules

### Package Descriptor (Required)
Every `.go` file must start with:
```go
// Package <name> <one-line description>
//
// Purpose:
//   <What this package does>
//
// Responsibilities:
//   - <Responsibility 1>
//
// Features:
//   - <Feature 1>  (e.g., Metrics, Tracing, Database)
//
// Constraints:
//   - <Constraint 1>
//
package <name>
```

**Key RULES:**
- - No metrics/tracing WITHOUT `Features: Metrics/Tracing` declaration
- - Features must be explicitly declared to be used
- See PACKAGE_DESCRIPTOR.md for full spec

### Metrics (Key metrics)
- - Package Descriptor: Required on ALL files
- - Feature Declaration: NO telemetry without explicit declaration
- - Functions: < 35 lines (strict)
- - Cyclomatic complexity: < 10 (`gocyclo -over 9 .` should return zero)
- - Test coverage: 100% required
- - golangci-lint: Zero warnings
- - gosec: Zero security issues
- - Code duplication: < 3%

### File Structure (Mandatory)
```
package/
├── constants.go           # ALL constants
├── errors.go              # ALL errors
├── interfaces.go          # ALL interfaces (package xxx)
├── interfaces_test.go     # ALL mocks (package xxx_test)
├── user.go               # User struct + methods
├── user_config.go        # UserConfig struct
├── order.go              # Order struct + methods
└── service_test.go       # Tests (package xxx_test)
```

**Key - One File Per Struct (Required):**
- Each struct must have its own dedicated file
- Example: `user.go` for User, `user_config.go` for UserConfig
- constants.go for ALL constants
- errors.go for ALL error definitions
- Better organization, easier navigation, cleaner Git conflicts

**Key - Test Package Naming:**
```go
// - After - Black-box testing
package taskqueue_test

import "taskqueue"

//  Before - Do NOT use same package
package taskqueue
```

**Key - Test Files (Required):**
- Test files must use `package xxx_test` (black-box)
- Import the package under test
- **Zero benchmarks in committed code** (use temporarily for POC/optimization only)
- NO `*_helper.go` files outside tests
- NO `*_bench.go` files
- NO `Benchmark*` functions in committed `_test.go` files

**Key - Benchmark Policy:**
-  **Never commit benchmarks** to the repository
- - Write benchmarks TEMPORARILY for performance optimization work
- - Run benchmarks locally to validate improvements
- - Delete benchmarks before committing
- 📋 Document performance improvements in commit messages (e.g., "3x faster via sync.Pool")

**Example - Temporary benchmarks (DO NOT COMMIT):**
```go
// ⚠️ TEMPORARY - user_test.go (for local POC only)
package user_test

func TestUser(t *testing.T) { ... }  // - Commit this

//  DELETE before commit - temporary benchmark
func BenchmarkUserCreation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        NewUser("test")
    }
}
```

### Constructor Pattern (Mandatory)
```go
// Every struct must have:
type ServiceConfig struct {
    Dep1 Interface1  // dependencies
    Dep2 Interface2
    Val1 string      // configuration
}

func NewService(cfg ServiceConfig) (*Service, error) {
    if cfg.Dep1 == nil {
        return nil, errors.New("dep1 required")
    }
    return &Service{...}, nil
}

//  Avoid: svc := &Service{...}
// - Required: svc, err := NewService(cfg)
```

### Refactoring Rules
**If function > 35 lines OR complexity > 9:**
1. Extract validation → `validateX()`
2. Extract data access → `fetchX()`
3. Extract business logic → `processX()`
4. Extract side effects → `notifyX()`
5. Main orchestrates only (10-20 lines)

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

## 🎯 Testability Requirements

**ALL code must be 100% testable:**
- All dependencies injected via constructor
- All external I/O behind interfaces
- Time, rand, I/O abstracted for mocking
- NO global state or singletons
- Every error path tested

## 📋 Key Best Practices

### Error Handling
- Never ignore errors: `_` for errors is Avoid
- Always wrap errors: `fmt.Errorf("context: %w", err)`
- Use `errors.Is()` and `errors.As()`, not `==`

### Naming
- Packages: lowercase, single word, no underscores
- Interfaces: `-er` suffix (Reader, Writer, Formatter)
- Receivers: 1-2 char abbreviation (u *User, not this/self)
- No stutter: `user.Repository`, not `user.UserRepository`

### Concurrency
- Always run with `-race` flag
- Close channels from sender only
- Use context for cancellation: `ctx context.Context` as first param
- Never leak goroutines - provide exit mechanism

## ⚡ Performance Optimization (Required)

### 1. Constants for ALL Default Values
**RULE**: NO magic numbers. ALL default values in constants.

```go
//  Before - Magic numbers
func NewService(cfg Config) (*Service, error) {
    if cfg.Timeout == 0 {
        cfg.Timeout = 30 * time.Second  //  Magic number
    }
    buffer := make(chan Task, 100)  //  Magic number
    return &Service{...}, nil
}

// - After - Named constants
const (
    DefaultTimeout    = 30 * time.Second
    DefaultBufferSize = 100
)

func NewService(cfg Config) (*Service, error) {
    if cfg.Timeout == 0 {
        cfg.Timeout = DefaultTimeout  // - Named constant
    }
    buffer := make(chan Task, DefaultBufferSize)  // - Named constant
    return &Service{...}, nil
}
```

### 2. Bitwise Operations for Flags
**RULE**: Use bitwise flags for boolean options, not multiple bools.

**Memory savings**: 1 byte vs 8+ bytes

```go
//  Before - Multiple bools (uses ~4-8 bytes)
type Task struct {
    IsUrgent    bool  // 1 byte + padding
    IsRetryable bool  // 1 byte + padding
    LogMetrics  bool  // 1 byte + padding
}

// - After - Bitwise flags (uses 1 byte)
type Task struct {
    Flags uint8  // 1 byte total
}

const (
    TaskFlagNone      uint8 = 0
    TaskFlagUrgent    uint8 = 1 << 0  // 0001 = 1
    TaskFlagRetryable uint8 = 1 << 1  // 0010 = 2
    TaskFlagMetrics   uint8 = 1 << 2  // 0100 = 4
)

// Operations
func (t *Task) HasFlag(flag uint8) bool {
    return t.Flags&flag != 0  // Bitwise AND
}

func (t *Task) SetFlag(flag uint8) {
    t.Flags |= flag  // Bitwise OR
}

func (t *Task) ClearFlag(flag uint8) {
    t.Flags &^= flag  // Bitwise AND NOT
}

// Usage
task := &Task{Flags: TaskFlagRetryable | TaskFlagMetrics}
task.SetFlag(TaskFlagUrgent)
if task.HasFlag(TaskFlagUrgent) {
    // Handle urgent task
}
```

### 3. map[string]struct{} for Sets
**RULE**: Use `map[T]struct{}` for sets, not `map[T]bool`.

**Memory savings**: 0 bytes vs 1 byte per entry

```go
//  Before - map[string]bool (uses 1 byte per entry)
var validStatuses = map[string]bool{
    "pending":    true,  // +1 byte per entry
    "processing": true,
    "completed":  true,
}

func IsValid(status string) bool {
    return validStatuses[status]
}

// - After - map[string]struct{} (uses 0 bytes per entry)
var validStatuses = map[string]struct{}{
    "pending":    {},  // 0 bytes per entry
    "processing": {},
    "completed":  {},
}

func IsValid(status string) bool {
    _, exists := validStatuses[status]
    return exists
}

// - Also good for deduplication
func RemoveDuplicates(items []string) []string {
    seen := make(map[string]struct{}, len(items))
    result := make([]string, 0, len(items))

    for _, item := range items {
        if _, exists := seen[item]; !exists {
            seen[item] = struct{}{}
            result = append(result, item)
        }
    }
    return result
}
```

### 4. Struct Field Ordering for Memory Alignment
**RULE**: Order struct fields by size (largest to smallest).

**Memory savings**: 20-50% reduction in struct size

**Field Sizes (64-bit)**:
- Pointers, slices, maps: 8 bytes
- time.Time: 24 bytes (3 × int64)
- string: 16 bytes (pointer + length)
- int64, uint64, float64: 8 bytes
- int32, uint32, float32: 4 bytes
- int16, uint16: 2 bytes
- int8, uint8, bool: 1 byte

```go
//  Before - Random ordering (~56 bytes with padding)
type User struct {
    ID        string    // 16 bytes
    Active    bool      // 1 byte + 7 padding
    CreatedAt time.Time // 24 bytes
    Age       int32     // 4 bytes + 4 padding
    Email     string    // 16 bytes
}

// - After - Ordered by size (~48 bytes)
type User struct {
    // 24-byte fields first
    CreatedAt time.Time // 24 bytes

    // 16-byte strings
    ID    string // 16 bytes
    Email string // 16 bytes

    // 4-byte fields
    Age int32 // 4 bytes

    // 1-byte fields last (packed together)
    Active bool // 1 byte
}
// Total: 48 bytes (17% smaller!)

// - Real example - Optimal ordering
type WorkerConfig struct {
    // 8-byte aligned (pointers/interfaces) first
    Repository TaskRepository   // 8 bytes
    Executor   TaskExecutor     // 8 bytes
    Publisher  MessagePublisher // 8 bytes
    Logger     *slog.Logger     // 8 bytes

    // 8-byte time.Duration
    ShutdownTimeout time.Duration // 8 bytes
    ProcessTimeout  time.Duration // 8 bytes

    // 8-byte int (on 64-bit)
    WorkerCount int // 8 bytes
    BufferSize  int // 8 bytes
}
```

### 5. chan struct{} for Signals
**RULE**: Use `chan struct{}` for signaling, not `chan bool`.

```go
//  Before
done := make(chan bool)
done <- true

// - After
done := make(chan struct{})
close(done)  // Or: done <- struct{}{}
```

### Documentation
- Every exported symbol must have godoc
- Start with symbol name: `// UserRepository manages...`
- Document thread-safety guarantees
- Document panics if method can panic

## 🚫 Common Violations (Key issues)

1. Missing Package Descriptor
2. Undeclared features (telemetry without Features declaration)
3. **Wrong test package** (using `package xxx` instead of `package xxx_test`)
4. **Multiple structs in one file** (must be 1 file per struct)
5. **Committed benchmarks** (Zero benchmarks in repo - temporary use only)
6. **Separate benchmark files** (`*_bench.go`)
7. Function > 35 lines
8. Cyclomatic complexity > 9
9. Coverage < 100%
10. Missing constructor
11. Missing Config for services
12. Wrong file structure
13. Ignored errors (`_`)
14. golangci-lint warnings
15. Security vulnerabilities
16. Direct struct literals
17. **Magic numbers** (using literals instead of constants)
18. **Multiple bools as flags** (should use bitwise uint8)
19. **map[T]bool for sets** (should use map[T]struct{})
20. **Unordered struct fields** (not ordered by size)

## - Success Checklist

Before submitting code:

**Architecture & Structure:**
- [ ] Package Descriptor on EVERY .go file
- [ ] Features explicitly declared (Metrics, Tracing, etc.)
- [ ] NO telemetry imports without Features declaration
- [ ] **1 file per struct** (user.go, user_config.go, etc.)
- [ ] constants.go exists with ALL constants
- [ ] errors.go exists with ALL errors
- [ ] interfaces.go exists with ALL interfaces
- [ ] interfaces_test.go exists with mocks (package xxx_test)
- [ ] Test files use `package xxx_test` (black-box)

**Code Quality:**
- [ ] All functions < 35 lines
- [ ] All functions gocyclo < 10
- [ ] 100% test coverage
- [ ] NO ignored errors
- [ ] golangci-lint passes
- [ ] gosec passes
- [ ] go test -race passes

**Design Patterns:**
- [ ] Every struct has NewXXXX() constructor
- [ ] Services have XXXXConfig struct
- [ ] All dependencies injected

**Performance Optimizations:**
- [ ] **NO magic numbers** - all defaults in constants
- [ ] **Bitwise flags** used instead of multiple bools
- [ ] **map[T]struct{}** used for sets (not map[T]bool)
- [ ] **Struct fields ordered by size** (largest first)
- [ ] **chan struct{}** used for signals (not chan bool)
- [ ] Pre-allocated slices with capacity
- [ ] strings.Builder in loops
- [ ] strconv instead of fmt.Sprintf for conversions

---
