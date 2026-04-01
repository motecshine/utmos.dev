# Repository Guidelines

## Project Structure & Module Organization
`cmd/` contains service entrypoints: `iot-api`, `iot-ws`, `iot-uplink`, `iot-downlink`, `iot-gateway`, and `dji-adapter`. Core service logic lives in `internal/` by domain (`api/`, `gateway/`, `uplink/`, `ws/`, `shared/`). Reusable packages live in `pkg/` (`adapter/`, `metrics/`, `rabbitmq/`, `tracer/`, `models/`). Contracts and public API definitions live in `api/v1/`. Environment configs are in `configs/`, Docker assets in `deployments/docker/`, architecture and protocol docs in `docs/`, and higher-level feature specs in `specs/`. Tests are split across `tests/integration`, `tests/contract`, and `tests/mocks`.

## Build, Test, and Development Commands
Use `make` targets as the default workflow:

- `make build` builds the five `iot-*` binaries into `bin/`.
- `make test` runs all Go tests with race detection and coverage output.
- `make test-short` skips longer integration paths.
- `make test-integration` runs only `tests/integration/...`.
- `make coverage` enforces the 80% coverage gate.
- `make lint` runs `golangci-lint`.
- `make fmt` applies `gofmt -s -w .`.
- `make run-api` or `make run-uplink` starts an individual service locally.
- `make docker-compose-up` starts local dependencies.

## Coding Style & Naming Conventions
Target Go first; follow Uber Go style and standard Go formatting. Let `gofmt` control indentation and spacing. Keep package names lowercase, exported identifiers in `CamelCase`, and test files suffixed with `_test.go`. Prefer descriptive domain-oriented names such as `service_handler.go` or `device_commands.go`. CI enforces `errcheck`, `govet`, `gosec`, `errorlint`, `forcetypeassert`, and `revive`; add exported comments and avoid unchecked type assertions.

## Testing Guidelines
Write tests alongside the code when possible and keep cross-service flows in `tests/integration`. Contract checks belong in `tests/contract` and should stay aligned with [api/v1/openapi.yaml](/Users/laputalaputa/Desktop/code/utmos.dev/api/v1/openapi.yaml). The repo expectation is test-first development with at least 80% coverage. Use `TestXxx` names, `testify` assertions, and `make test` before opening a PR.

## Commit & Pull Request Guidelines
Recent history uses Conventional Commits, for example `feat(observability): ...`, `refactor(lint): ...`, and `docs(004): ...`. Keep the type lowercase, add a focused scope when useful, and make each commit buildable. PRs should summarize behavior changes, list impacted services, note config or schema changes, and include command evidence such as `make test`, `make lint`, and `make coverage`. If you change API or protocol behavior, update the matching docs and specs in the same PR.
