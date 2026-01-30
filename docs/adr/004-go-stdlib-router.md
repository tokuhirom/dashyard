# ADR-004: Go Standard Library Router

## Status
Accepted

## Context
We need an HTTP router for the Go backend API. Options considered were gorilla/mux, chi, and Go 1.22+ enhanced http.ServeMux.

## Decision
Use Go 1.22+ `http.ServeMux` enhanced routing with method and path wildcard support.

## Rationale
- Eliminates need for third-party router dependencies
- Go 1.22+ ServeMux supports `GET /path`, `POST /path`, and `{wildcard}` patterns
- Sufficient for the API surface of this project
- Reduces dependency maintenance burden

## Consequences
- Limited to routing features available in the standard library
- No middleware chaining built-in (handled manually)
