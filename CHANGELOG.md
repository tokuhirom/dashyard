# Changelog

## [v0.4.1](https://github.com/tokuhirom/dashyard/compare/v0.4.0...v0.4.1) - 2026-02-01
- Add Docker-based E2E test runner by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/117

## [v0.3.1](https://github.com/tokuhirom/dashyard/compare/v0.3.0...v0.3.1) - 2026-02-01
- Add dummyapp and reorganize examples directory by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/112
- Remove unnecessary host port mappings and local dev targets by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/114
- Add support for multiple named datasources by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/115

## [v0.3.0](https://github.com/tokuhirom/dashyard/compare/v0.2.1...v0.3.0) - 2026-02-01
- Reorganize examples directory and add CI validation by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/110

## [v0.2.1](https://github.com/tokuhirom/dashyard/compare/v0.2.0...v0.2.1) - 2026-02-01
- Bump actions/setup-node from 4 to 6 by @dependabot[bot] in https://github.com/tokuhirom/dashyard/pull/99
- Add OAuth/OIDC support with Goth and gorilla/sessions by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/106
- Add GitHub OAuth documentation and screenshot to README by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/107
- Add --dashboards-dir CLI flag, remove dashboards.dir from config files by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/108
- Run validate on example config and dashboards in CI by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/109

## [v0.2.0](https://github.com/tokuhirom/dashyard/compare/v0.1.0...v0.2.0) - 2026-01-31
- Add /ready endpoint by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/90
- Add frontend unit tests with Vitest by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/91
- Run Docker container as non-root user by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/93
- Add Dependabot for dependency scanning by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/94
- Expose Dashyard's own metrics via /metrics endpoint by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/95
- Bump actions/setup-go from 5 to 6 by @dependabot[bot] in https://github.com/tokuhirom/dashyard/pull/97
- Bump actions/upload-artifact from 4 to 6 by @dependabot[bot] in https://github.com/tokuhirom/dashyard/pull/96
- Bump Songmu/tagpr from 1.11.1 to 1.14.0 by @dependabot[bot] in https://github.com/tokuhirom/dashyard/pull/101
- Bump @vitejs/plugin-react from 4.7.0 to 5.1.2 in /frontend by @dependabot[bot] in https://github.com/tokuhirom/dashyard/pull/100
- Bump react-markdown from 9.1.0 to 10.1.0 in /frontend by @dependabot[bot] in https://github.com/tokuhirom/dashyard/pull/102
- Bump vite from 6.4.1 to 7.3.1 in /frontend by @dependabot[bot] in https://github.com/tokuhirom/dashyard/pull/103
- Bump actions/checkout from 4 to 6 by @dependabot[bot] in https://github.com/tokuhirom/dashyard/pull/98

## [v0.1.0](https://github.com/tokuhirom/dashyard/compare/v0.0.13...v0.1.0) - 2026-01-31
- Add oxlint for frontend linting by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/78
- Add logarithmic Y-axis scale support by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/82
- Rewrite README with AI-native positioning by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/83
- Add Dashyard logo by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/84

## [v0.0.13](https://github.com/tokuhirom/dashyard/compare/v0.0.12...v0.0.13) - 2026-01-31
- Add gen-prompt subcommand for LLM dashboard generation by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/61
- Replace dummyprom gen-prompt with real monitoring stack by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/63
- Remove pie and doughnut chart types by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/68
- Add seconds unit for duration/latency metrics by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/69
- Add PromQL division-by-zero guard guidance by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/70
- Add full_width option for markdown panels by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/71
- Add make screenshots to CLAUDE.md by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/72
- Fix gen-prompt: force rebuild Docker image before run by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/73
- Add make gen-prompt-up to start full monitoring stack by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/74
- Fix formatBytes showing undefined for fractional values by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/75
- add auto generated dashboards by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/77

## [v0.0.12](https://github.com/tokuhirom/dashyard/compare/v0.0.11...v0.0.12) - 2026-01-31
- Add auto-refresh interval selector for dashboards by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/53
- Add issue linking convention to CLAUDE.md by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/56
- Add stacked chart support for graph panels by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/55
- Add trusted proxy support for X-Forwarded-For by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/57
- Remove IP ACL feature, keep only trusted proxy support by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/58
- Add validate subcommand by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/59
- Make screenshot script stable with networkidle waits by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/60

## [v0.0.11](https://github.com/tokuhirom/dashyard/compare/v0.0.10...v0.0.11) - 2026-01-31
- Make dummyprom output deterministic for stable screenshots by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/49
- Fix Go version requirement in README by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/51

## [v0.0.10](https://github.com/tokuhirom/dashyard/compare/v0.0.9...v0.0.10) - 2026-01-31
- Add Docker Compose workflow for automated screenshots by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/46
- Hide variable selectors for repeat-only variables by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/43
- Add absolute time range support by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/48

## [v0.0.9](https://github.com/tokuhirom/dashyard/compare/v0.0.8...v0.0.9) - 2026-01-31
- Support threshold lines on graph panels by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/41

## [v0.0.8](https://github.com/tokuhirom/dashyard/compare/v0.0.7...v0.0.8) - 2026-01-31
- Move host/port from config file to CLI flags by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/32
- Limit percent y-axis range to 0–100 by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/33
- Document chart_type and unit options in README by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/35
- Support y_min and y_max on graph panels by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/36

## [v0.0.7](https://github.com/tokuhirom/dashyard/compare/v0.0.6...v0.0.7) - 2026-01-31
- Add golangci-lint and fix errcheck violations by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/27
- Update README with feature screenshots by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/29
- Fix CI lint job: build frontend before linting by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/30
- Add installation guide to README by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/31

## [v0.0.6](https://github.com/tokuhirom/dashyard/compare/v0.0.5...v0.0.6) - 2026-01-31
- Add dashboard template variables and repeat rows by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/25

## [v0.0.5](https://github.com/tokuhirom/dashyard/compare/v0.0.4...v0.0.5) - 2026-01-31
- Add hot-reload for dashboard files using fsnotify by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/20

## [v0.0.4](https://github.com/tokuhirom/dashyard/compare/v0.0.3...v0.0.4) - 2026-01-31
- Add serve and mkpasswd subcommands by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/18

## [v0.0.3](https://github.com/tokuhirom/dashyard/compare/v0.0.2...v0.0.3) - 2026-01-31
- Fix goreleaser before hooks YAML format by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/16
- Add branch creation step to git workflow by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/15
- Add under-development caution to README by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/9
- Replace flag with kong for CLI parsing by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/13

## [v0.0.2](https://github.com/tokuhirom/dashyard/compare/v0.0.1...v0.0.2) - 2026-01-31
- Add git workflow guidelines to CLAUDE.md by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/7
- Add E2E testing section to README by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/8
- Add MIT license by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/11
- Add GoReleaser setup for automated releases by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/12

## [v0.0.2](https://github.com/tokuhirom/dashyard/compare/v0.0.1...v0.0.2) - 2026-01-31
- Add git workflow guidelines to CLAUDE.md by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/7
- Add E2E testing section to README by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/8
- Add MIT license by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/11
- Add GoReleaser setup for automated releases by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/12

## [v0.0.1](https://github.com/tokuhirom/dashyard/commits/v0.0.1) - 2026-01-31
- Bump golang.org/x/crypto from 0.40.0 to 0.45.0 by @dependabot[bot] in https://github.com/tokuhirom/dashyard/pull/2
- Bump github.com/quic-go/quic-go from 0.54.0 to 0.57.0 by @dependabot[bot] in https://github.com/tokuhirom/dashyard/pull/1
- Add Playwright E2E tests with CI by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/4
- Add tagpr workflow for automated release tagging by @tokuhirom in https://github.com/tokuhirom/dashyard/pull/5
