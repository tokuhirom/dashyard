# ADR-006: Development Tooling - Devbox + direnv

## Status
Accepted

## Context
We need consistent development environments across contributors. Options considered were system-level installs, Docker-based dev environments, and Nix-based tooling.

## Decision
Use devbox for reproducible development environments, integrated with direnv for automatic shell activation.

## Rationale
- Ensures consistent Go and Node.js versions across developers
- Does not pollute the global system with project-specific tool versions
- direnv auto-activates the devbox environment when entering the project directory
- Simple setup: `devbox init && devbox add go nodejs`

## Consequences
- Requires devbox and direnv to be installed on developer machines
- Nix store usage for tool management
