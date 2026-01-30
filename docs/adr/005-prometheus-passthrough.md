# ADR-005: Prometheus Response Pass-through

## Status
Accepted

## Context
We need to proxy Prometheus query_range API responses to the frontend. Options were to parse and transform the response on the backend, or pass it through as-is.

## Decision
Stream raw Prometheus JSON responses to the frontend without server-side parsing.

## Rationale
- Simpler backend implementation with no response transformation logic
- Faster with no deserialization/reserialization overhead
- Forward-compatible with Prometheus API changes
- Frontend already needs to understand Prometheus response format for rendering

## Consequences
- Frontend must handle Prometheus response format directly
- No server-side validation of Prometheus responses
