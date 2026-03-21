<!-- updated: 2026-02-27T13:00:00Z -->

# Shared Hook Utilities

## Purpose

Common utilities and configuration shared across lifecycle hooks.

## Files

| File           | Description            |
| -------------- | ---------------------- |
| `utils.sh`     | Bash utility functions |
| `.env.example` | Environment template   |

## utils.sh Functions

Source with: `source "$(dirname "$0")/../shared/utils.sh"`

Common utilities for:

- Logging and output formatting
- Version checking
- Tool installation helpers
- Environment validation

## Conventions

- Keep utilities generic and reusable
- Document function parameters
- Use consistent error handling
