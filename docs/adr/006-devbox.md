# ADR-006: Development Tooling - mise

## Status
Accepted (supersedes previous devbox + direnv decision)

## Context
We need consistent development environments across contributors. Previously we used devbox + direnv (Nix-based). We migrated to mise for simpler setup and faster environment activation.

## Decision
Use mise for development environment management with `.mise.toml` configuration.

## Rationale
- Ensures consistent Go and Node.js versions across developers
- Does not pollute the global system with project-specific tool versions
- Simpler installation and faster activation than Nix-based devbox
- Shell integration (`mise activate`) auto-activates when entering the project directory
- Simple setup: `mise install`

## Consequences
- Requires mise to be installed on developer machines
- Tool versions managed via `.mise.toml` in the project root
