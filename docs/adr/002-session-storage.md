# ADR-002: Session Storage - Cookie Only

## Status
Accepted

## Context
We need session management for authenticated API access. The options considered were cookie-only (stateless) vs cookie + server-side session storage.

## Decision
Use stateless sessions with HMAC-SHA256 signed cookies.

## Rationale
- Stateless design is simpler to implement and maintain
- No server-side storage needed (no database or memory store)
- Scales easily for single-instance deployment
- Cookie contains signed JSON payload with user ID and expiry

## Consequences
- Cannot revoke individual sessions without changing the signing secret
- Acceptable for MVP; server restart clears all sessions by generating a new secret
- Session data size limited by cookie size constraints
