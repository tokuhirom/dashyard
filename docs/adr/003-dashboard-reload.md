# ADR-003: Dashboard Reload - Manual Restart

## Status
Accepted

## Context
We need to decide how dashboard YAML definitions are loaded and refreshed. Options considered were hot reload via fsnotify vs manual restart.

## Decision
Dashboards are loaded once at startup. A restart is required to pick up changes.

## Rationale
- Simpler implementation with no file watching complexity
- Prevents unexpected dashboard changes in production
- Predictable behavior for operators

## Consequences
- Requires container restart to pick up dashboard changes
- Hot reload deferred to Phase 2 if needed
