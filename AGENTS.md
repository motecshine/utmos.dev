# Repository Guidelines
Use this guide to contribute to the DJI docs and Go reference code.

## Project Structure & Module Organization
- `docs/dji/00.dji-wpml` holds the WPML spec Markdown; keep the numeric prefixes when adding pages.
- `docs/dji/wpml` contains Go code for WPML mission serialization/validation with tests in the same package.
- `docs/dji/protocol` stores Go structs for Pilot-to-Cloud/Dock payloads, split by domain (`aircraft`, `wayline`, `file`, etc.).
- `docs/dji/10.overview`, `20.quick-start`, `30.feature-set`, and `50.debug` are user-facing docs; preserve numbering for navigation.
- `.specify/templates` and `.specify/scripts` house contributor templates and helper scripts—reuse them over ad-hoc files.

## Build, Test, and Development Commands
- Go modules are absent; run with `GO111MODULE=off` or add a `go.mod` before formal CI usage.
- `cd docs/dji/wpml && GO111MODULE=off go test ./...` runs the WPML suite.
- `cd docs/dji/protocol && GO111MODULE=off go test ./...` checks protocol payload helpers.
- `gofmt -w $(find docs/dji -name '*.go')` keeps Go sources consistent.
- Markdown has no build step; preview locally and keep numbering prefixes intact.

## Coding Style & Naming Conventions
- Standard Go style: gofmt formatting, PascalCase exports, camelCase fields, `_test.go` for tests.
- Keep method/topic strings and struct fields aligned with DJI SDK naming; avoid ad-hoc renames.
- Place new Go packages under `docs/dji/<domain>` to mirror the spec; colocate helpers and tests.
- Docs: start with a single `#`, favor short paragraphs/tables, and follow the numeric filename pattern.

## Testing Guidelines
- Extend existing table-driven tests in `*_test.go` before creating new suites.
- Reuse fixtures in `docs/dji/wpml/test_helpers.go` for mission/waypoint scenarios.
- Run `go test ./... -cover` in each touched package; keep `final_coverage_test.go` happy by preserving coverage.
- For docs, validate JSON/YAML examples with `jq`/`yamllint` when possible and keep request/response samples runnable.

## Commit & Pull Request Guidelines
- With no history yet, default to Conventional Commits (e.g., `docs: clarify pilot mqtt topics`, `fix: guard nil waypoint payload`).
- Subjects stay imperative and ≤72 characters; include scopes like `wpml`, `protocol/wayline`, or `docs`.
- PRs should explain intent, list test commands run, and flag spec/doc updates; link issues/tasks when available.
- For protocol/doc changes, add before/after snippets or payload diffs to speed review.

## Active Technologies
- Go 1.21 + EMQX MQTT broker, RabbitMQ, chi (REST) + gorilla/websocket, InfluxDB client, OpenTelemetry SDK, Prometheus exporter, Loki/Tempo clients (001-feature-project-setup)
- InfluxDB for telemetry/events; RabbitMQ for queues; PostgreSQL for registry/command status/tenancy metadata (001-feature-project-setup)

## Recent Changes
- 001-feature-project-setup: Added Go 1.21 + EMQX MQTT broker, RabbitMQ, chi (REST) + gorilla/websocket, InfluxDB client, OpenTelemetry SDK, Prometheus exporter, Loki/Tempo clients
