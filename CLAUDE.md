# CLAUDE.md

Project context for Claude Code.

## What is this?

Dashyard is a lightweight Prometheus metrics dashboard. Go/Gin backend serves a React/TypeScript frontend as embedded static files. Dashboards are defined in YAML files.

## Build & Run

```bash
make build                       # Build frontend then backend
make test                        # Go tests
make gen-prompt                  # Generate LLM prompt from dummyprom metrics
make screenshots                 # Regenerate README screenshots (Docker)
cd frontend && npm run build     # Frontend build (TypeScript check + Vite)
```

Development (three terminals):
```bash
make dev-dummyprom               # Fake Prometheus on :9090
make dev-backend                 # Go server on :8080 (uses examples/kitchensink/config.yaml)
make dev-frontend                # Vite dev server on :5173
make dev-dummygithub             # Fake GitHub OAuth server on :5555
```

## Project Layout

- `main.go` -- Entry point, embeds `frontend/dist/` via `go:embed`
- `internal/config/` -- YAML config parsing with defaults
- `internal/dashboard/` -- Loads dashboard YAML files from a directory, builds nav tree
- `internal/handler/` -- Gin HTTP handlers (login, dashboards, query proxy, static)
- `internal/model/` -- Data structs: Dashboard, Row, Panel
- `internal/auth/` -- SHA-512 crypt password hashing, session middleware
- `internal/prometheus/` -- HTTP client for Prometheus `query_range` API
- `internal/server/` -- Gin router and middleware wiring
- `frontend/` -- React 19 + TypeScript + Vite + Chart.js
- `schemas/` -- JSON schemas for config.yaml and dashboard YAML
- `examples/kitchensink/` -- Kitchen-sink demo configs and dashboards
- `examples/real-world/` -- Real monitoring stack with gen-prompt workflow
- `cmd/dummyprom/` -- Fake Prometheus that generates synthetic metrics
- `cmd/dummyapp/` -- Fake web app exposing custom Prometheus metrics (HTTP RED, business KPIs)
- `cmd/dummygithub/` -- Fake GitHub OAuth server for local development

## Key Patterns

- Config fields flow from `config.yaml` -> `Config` struct -> handler -> API JSON response -> React props (e.g. `site_title`, `header_color`)
- Frontend is embedded in the binary at build time -- always build frontend before backend
- Dashboard files support subdirectories which become tree groups in the sidebar
- Panel types: `graph` (PromQL query) and `markdown` (rendered content)
- Tests are colocated: `foo.go` has `foo_test.go` in the same package

## Testing

```bash
go test ./...                    # All Go tests
go test ./internal/config/...    # Specific package
cd frontend && npm run build     # Type-checks and builds frontend
```

## Pre-commit Checklist

Before every commit, always run these checks locally and fix any issues:

```bash
golangci-lint run ./...          # Lint (errcheck, staticcheck, etc.)
go test ./...                    # All Go tests
cd frontend && npm run build     # Frontend type-check + build
```

If CI fails after pushing, fix issues and run all checks locally before pushing again.

## Git Workflow

- When starting a new task, always pull the latest main and create a new branch from it first (`git checkout main && git pull origin main && git checkout -b <branch>`)
- Create a separate branch and pull request for each feature or change
- Do not mix unrelated changes in a single branch/PR
- Include a link to the related issue in PR descriptions (e.g. `Fixes #123`)

## Go Module

`github.com/tokuhirom/dashyard` -- Go 1.25, Gin web framework, `gopkg.in/yaml.v3`
