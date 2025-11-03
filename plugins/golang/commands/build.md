# /build - Orchestrated code development with task tracking

Build production-grade Go code by orchestrating multiple expert agents with iterative validation until all success criteria are met.

## Usage

### Create new task

```bash
/build <description>
```

Creates a new TASK ID, orchestrates agents to implement the feature, and iterates until all criteria are met.

### Continue existing task

```bash
/build task <TASK_ID> <additional_context>
```

Adds context/constraints to existing task, re-runs validation with new requirements, and maintains task isolation.

## Examples

```bash
# Create new endpoint
/build ajoute un nouveau endpoint en prenant exemple sur l'existant
# → Output: TASK_ID: 20251011-abc123

# Add unit tests to existing task
/build task 20251011-abc123 ajoute les test unitaires

# Optimize performance
/build task 20251011-abc123 optimise les performance pprof

# Correct implementation
/build task 20251011-def456 tu aurais du ajouter tel valeur

# Achieve full coverage
/build task 20251011-abc123 gère tel partie jusqu'a avoir 100% de coverage
```

## Process Flow

### 1. Task Initialization

- Generate unique TASK_ID: `{YYYYMMDD}-{6-char-hash}`
- Create task context: `.tasks/<TASK_ID>.json`
- Parse requirements and constraints
- Initialize success criteria

### 2. Agent Orchestration

Based on requirements, orchestrate agents in order:

1. **DDD Architect** - Architecture and structure validation
2. **Go Expert** - Implementation with Go 1.25+ features
3. **Code Reviewer** - Quality validation
4. **Performance Optimizer** - pprof profiling and optimization

### 3. Iterative Development Loop

For each iteration (max 10):

```
┌─────────────────────────────────────────┐
│ 1. Implement/modify code                │
│ 2. Run validation checks                │
│ 3. Collect metrics                       │
│ 4. Verify success criteria              │
│ 5. If not met → adjust and retry        │
│    If met → mark as complete             │
└─────────────────────────────────────────┘
```

### 4. Success Criteria (Additive)

Each task validates against ALL criteria:

- ✅ **Compiles**: `go build` succeeds
- ✅ **Tests pass**: All unit and integration tests pass
- ✅ **Coverage**: Meets target (default 80%, can specify 100%)
- ✅ **No lint issues**: golangci-lint passes
- ✅ **Performance**: Benchmarks meet targets
- ✅ **Architecture**: DDD rules enforced
- ✅ **Review approved**: Code Reviewer validation

### 5. Task State Management

Each task maintains state in `.tasks/<TASK_ID>.json`:

```json
{
  "task_id": "20251011-abc123",
  "created_at": "2025-10-11T10:30:00Z",
  "updated_at": "2025-10-11T11:45:00Z",
  "description": "ajoute un nouveau endpoint en prenant exemple sur l'existant",
  "status": "in_progress",
  "current_iteration": 2,
  "max_iterations": 10,
  "files_changed": [
    "internal/api/handlers/user.go",
    "internal/api/handlers/user_test.go"
  ],
  "iterations": [
    {
      "iteration": 1,
      "timestamp": "2025-10-11T10:30:00Z",
      "context": "initial implementation",
      "agents_used": ["DDD Architect", "Go Expert"],
      "files_modified": ["internal/api/handlers/user.go"],
      "criteria_results": {
        "compiles": true,
        "tests_pass": true,
        "coverage_target": false,
        "lint_clean": true,
        "performance_ok": null,
        "architecture_valid": true,
        "review_approved": false
      },
      "metrics": {
        "lines_added": 127,
        "tests_added": 0,
        "coverage_pct": 0
      }
    },
    {
      "iteration": 2,
      "timestamp": "2025-10-11T11:45:00Z",
      "context": "ajoute les test unitaires",
      "agents_used": ["Go Expert", "Code Reviewer"],
      "files_modified": ["internal/api/handlers/user_test.go"],
      "criteria_results": {
        "compiles": true,
        "tests_pass": true,
        "coverage_target": true,
        "lint_clean": true,
        "performance_ok": null,
        "architecture_valid": true,
        "review_approved": false
      },
      "metrics": {
        "lines_added": 234,
        "tests_added": 15,
        "coverage_pct": 95.7
      }
    }
  ],
  "criteria": {
    "compiles": {"required": true, "met": true},
    "tests_pass": {"required": true, "met": true},
    "coverage_target": {"required": true, "met": true, "target": 80, "actual": 95.7},
    "lint_clean": {"required": true, "met": true},
    "performance_ok": {"required": false, "met": null},
    "architecture_valid": {"required": true, "met": true},
    "review_approved": {"required": true, "met": false}
  },
  "context_additions": [
    "ajoute les test unitaires",
    "optimise les performance pprof"
  ]
}
```

## Implementation

When this command is invoked, execute the following workflow:

### Phase 1: Task Setup

```bash
# Generate TASK_ID
TASK_ID=$(date +%Y%m%d)-$(openssl rand -hex 3)

# Create task directory if needed
mkdir -p .tasks

# Initialize task file
cat > .tasks/${TASK_ID}.json <<EOF
{
  "task_id": "${TASK_ID}",
  "created_at": "$(date -Iseconds)",
  "description": "$USER_INPUT",
  "status": "initialized",
  "iterations": []
}
EOF

# Output TASK_ID to user
echo "🚀 TASK ID: ${TASK_ID}"
echo "📝 Description: $USER_INPUT"
echo ""
```

### Phase 2: Agent Coordination

#### Step 1: Architecture Validation (DDD Architect)

```
You are the DDD Architect agent working on TASK ${TASK_ID}.

Task Description: ${DESCRIPTION}
Previous Context: ${CONTEXT_ADDITIONS}

Validate and ensure:
1. Domain structure follows DDD principles
2. Layer separation is maintained
3. Naming conventions are correct
4. Interfaces are properly defined
5. Dependencies point inward

If implementing new code:
- Determine which domain/layer it belongs to
- Define interfaces first
- Ensure proper package structure
- Document architectural decisions

Output format:
- Architecture plan
- Files to create/modify
- Interface definitions
- Any violations or concerns
```

#### Step 2: Implementation (Go Expert)

```
You are the Go Expert agent working on TASK ${TASK_ID}.

Task Description: ${DESCRIPTION}
Architecture Plan: ${ARCHITECTURE_OUTPUT}
Previous Context: ${CONTEXT_ADDITIONS}

Implement the feature following:
1. Use Go 1.25+ features where appropriate
2. Follow the architecture plan exactly
3. Implement all error handling
4. Add comprehensive documentation
5. Use examples from existing codebase as reference

Write production-grade code with:
- Proper error handling
- Context propagation
- Logging where appropriate
- Type safety
- Idiomatic Go patterns

Output: Complete implementation
```

#### Step 3: Testing (Go Expert)

```
You are the Go Expert agent working on TASK ${TASK_ID}.

Previous Implementation: ${IMPLEMENTATION}
Coverage Target: ${COVERAGE_TARGET}%

Write comprehensive tests:
1. Unit tests for all functions
2. Table-driven tests for multiple cases
3. Edge case testing
4. Error path testing
5. Integration tests if needed

Achieve ${COVERAGE_TARGET}% coverage minimum.

Output: Complete test suite
```

#### Step 4: Validation & Metrics

```bash
# Compile check
go build ./...
COMPILE_STATUS=$?

# Run tests
go test -v -race -coverprofile=coverage.out ./...
TEST_STATUS=$?

# Calculate coverage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

# Lint check
golangci-lint run
LINT_STATUS=$?

# Update task metrics
jq ".iterations += [{
  \"iteration\": $CURRENT_ITERATION,
  \"criteria_results\": {
    \"compiles\": $([ $COMPILE_STATUS -eq 0 ] && echo true || echo false),
    \"tests_pass\": $([ $TEST_STATUS -eq 0 ] && echo true || echo false),
    \"coverage_target\": $(awk -v c=$COVERAGE -v t=$TARGET 'BEGIN {print (c >= t ? "true" : "false")}'),
    \"lint_clean\": $([ $LINT_STATUS -eq 0 ] && echo true || echo false)
  },
  \"metrics\": {
    \"coverage_pct\": $COVERAGE
  }
}]" .tasks/${TASK_ID}.json > .tasks/${TASK_ID}.json.tmp
mv .tasks/${TASK_ID}.json.tmp .tasks/${TASK_ID}.json
```

#### Step 5: Performance Analysis (if requested)

```
You are the Performance Optimizer agent working on TASK ${TASK_ID}.

Implementation: ${IMPLEMENTATION}
Context: ${CONTEXT_ADDITIONS}

Analyze and optimize:
1. Profile with pprof (CPU, memory, allocations)
2. Identify bottlenecks
3. Optimize hot paths
4. Run benchmarks
5. Compare before/after metrics

Target improvements:
- Reduce allocations
- Improve throughput
- Optimize critical paths

Output: Performance report and optimizations
```

#### Step 6: Final Review (Code Reviewer)

```
You are the Code Reviewer agent working on TASK ${TASK_ID}.

Complete implementation: ${ALL_CODE}
Test suite: ${ALL_TESTS}
Metrics: ${METRICS}

Perform thorough review:
1. All standards enforced
2. Best practices followed
3. Documentation complete
4. Tests comprehensive
5. Code is production-ready

Approval criteria:
- All previous criteria met
- Code quality exceptional
- No technical debt
- Ready for production

Output: APPROVED or list of issues to fix
```

### Phase 3: Iteration Decision

```bash
# Check if all required criteria are met
ALL_MET=$(jq -r '.criteria | to_entries | map(select(.value.required == true and .value.met == false)) | length == 0' .tasks/${TASK_ID}.json)

if [ "$ALL_MET" = "true" ]; then
  # Update status to complete
  jq '.status = "completed"' .tasks/${TASK_ID}.json > .tasks/${TASK_ID}.json.tmp
  mv .tasks/${TASK_ID}.json.tmp .tasks/${TASK_ID}.json

  echo "✅ TASK COMPLETE: ${TASK_ID}"
  echo "All criteria met in $CURRENT_ITERATION iterations"
else
  # Increment iteration and retry
  CURRENT_ITERATION=$((CURRENT_ITERATION + 1))

  if [ $CURRENT_ITERATION -gt $MAX_ITERATIONS ]; then
    echo "❌ Max iterations reached for TASK ${TASK_ID}"
    echo "Manual intervention required"
    jq '.status = "blocked"' .tasks/${TASK_ID}.json > .tasks/${TASK_ID}.json.tmp
    mv .tasks/${TASK_ID}.json.tmp .tasks/${TASK_ID}.json
  else
    echo "🔄 Iteration $CURRENT_ITERATION/$MAX_ITERATIONS"
    echo "Retrying with adjusted approach..."
    # Go back to Phase 2
  fi
fi
```

## Output Format

### During Execution

```
🚀 TASK ID: 20251011-abc123
📝 Description: ajoute un nouveau endpoint en prenant exemple sur l'existant

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔄 Iteration 1/10

  🏗️  Phase 1: Architecture Validation (DDD Architect)
      ✅ Domain structure validated
      ✅ Layer separation correct
      ✅ Interfaces defined
      📝 Files to modify:
         - internal/domain/user/repository.go (interface)
         - internal/api/handlers/user.go (implementation)
         - internal/infrastructure/postgres/user_repo.go (adapter)

  💻 Phase 2: Implementation (Go Expert)
      ✅ Code implemented with Go 1.25 features
      ✅ Error handling complete
      ✅ Documentation added
      📊 Lines added: 127

  🧪 Phase 3: Testing (Go Expert)
      ⏳ Tests in progress...
      ✅ Unit tests: 15 tests written
      ✅ Table-driven tests for edge cases
      📊 Coverage: 0% → 95.7%

  ✓ Validation Results:
      ✅ Compiles: SUCCESS
      ✅ Tests pass: 15/15 PASSED
      ✅ Coverage: 95.7% (target: 80%)
      ✅ Lint: NO ISSUES
      ⏳ Architecture: VALID
      ⏳ Review: PENDING

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔄 Iteration 2/10

  👀 Phase 6: Final Review (Code Reviewer)
      ⏳ Reviewing implementation...
      ✅ Standards enforced
      ✅ Documentation complete
      ✅ Tests comprehensive
      ✅ Production-ready

      APPROVED ✅

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ TASK COMPLETE: 20251011-abc123

📊 Summary:
   Iterations: 2/10
   Files changed: 3
   Lines added: 361
   Tests added: 15
   Coverage: 95.7%
   Duration: 3m 42s

📁 Files modified:
   - internal/domain/user/repository.go
   - internal/api/handlers/user.go
   - internal/api/handlers/user_test.go

💡 Continue with:
   /build task 20251011-abc123 <additional requirements>
```

### For Continued Tasks

```
🔄 TASK ID: 20251011-abc123 (continuation)
📝 New context: ajoute les test unitaires

Previous state:
  ✅ Compiles
  ✅ Tests pass (15 tests)
  ✅ Coverage: 95.7%
  ✅ Lint clean

New requirements: ajoute les test unitaires

[... continue execution ...]
```

## Task Isolation

Each TASK_ID maintains complete isolation:

1. **Separate context**: Each task has its own JSON file
2. **Independent criteria**: Success criteria don't interfere
3. **File tracking**: Changes tracked per task
4. **Additive constraints**: New context adds to existing, never replaces
5. **No cross-talk**: Task XXXX changes don't affect task YYYY

## Task Management Commands

### List all tasks

```bash
# List active tasks
find .tasks -name "*.json" -type f -exec sh -c 'echo "$(basename {} .json): $(jq -r .description {})"' \;

# List by status
jq -r 'select(.status == "in_progress") | .task_id + ": " + .description' .tasks/*.json
```

### View task details

```bash
# Full details
cat .tasks/20251011-abc123.json | jq .

# Summary only
jq '{task_id, description, status, iterations: (.iterations | length), criteria}' .tasks/20251011-abc123.json
```

### Archive completed tasks

```bash
mkdir -p .tasks/archive
mv .tasks/20251011-abc123.json .tasks/archive/
```

## Configuration

Default settings (can be overridden):

- **Max iterations**: 10
- **Coverage target**: 80%
- **Required criteria**: compiles, tests_pass, coverage_target, lint_clean, architecture_valid, review_approved
- **Optional criteria**: performance_ok

Override with:

```bash
/build --max-iter=15 --coverage=100 <description>
```

## Error Handling

### If iteration fails

1. Capture error details in task JSON
2. Add error context for next iteration
3. Adjust agent approach
4. Retry with new strategy

### If max iterations reached

1. Mark status as "blocked"
2. Report blocking issues
3. Preserve complete state
4. Suggest manual intervention
5. Allow resume with additional context

### If criteria impossible

1. Mark criteria as "not_applicable"
2. Document reasoning
3. Require explicit override
4. Continue with remaining criteria

## Integration

### With Makefile

Use Makefile commands for:

- `make build` - Compilation
- `make test` - Testing
- `make lint` - Linting
- `make bench` - Benchmarking

### With Codacy

Validate against Codacy metrics:

- Code quality grade
- Issues count
- Coverage metrics
- Security scanning

### With Git

Track task progress:

```bash
# Commit per iteration
git add .
git commit -m "feat: ${DESCRIPTION} [TASK: ${TASK_ID}] [ITER: ${ITERATION}]"
```

## Best Practices

1. **Clear descriptions**: Be specific and actionable
2. **Incremental additions**: Add one requirement at a time
3. **Test early**: Add tests before optimization
4. **Performance last**: Optimize after functionality is complete
5. **Isolated tasks**: One concern per TASK_ID
6. **Reference existing code**: Use "en prenant exemple sur l'existant"
7. **Explicit targets**: Specify coverage/performance targets if different from defaults
8. **Review task state**: Check `.tasks/<TASK_ID>.json` before adding context

## Notes

- Task files stored in `.tasks/` (gitignored, local only)
- TASK_ID format: `YYYYMMDD-{6-hex}` (e.g., `20251011-abc123`)
- All agents work with shared task context
- Success criteria are additive and never removed
- Each iteration builds on previous work
- Task isolation prevents context pollution
- Maximum 10 iterations prevents infinite loops
- All metrics tracked for transparency
