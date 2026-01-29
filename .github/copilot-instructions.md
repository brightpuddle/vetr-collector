# ACI vetR Collector - AI Coding Agent Instructions

## Project Overview

This is a Go-based data collector that queries Cisco ACI APIC controllers via REST API. It fetches configuration and operational data for health checks performed by Cisco Services. The tool produces a `aci-vetr-data.zip` file containing JSON responses from various ACI managed object classes.

**Key architectural components:**
- `cmd/collector/main.go` - Entry point with batch orchestration logic
- `pkg/aci/client.go` - HTTP client with automatic token refresh (every 480s)
- `pkg/cli/cli.go` - API fetching with retry logic and pagination for large datasets
- `pkg/req/reqs.json` - Embedded YAML defining ~100 ACI classes to query
- `pkg/archive/archive.go` - Thread-safe zip writer using mutex locks

## Critical Patterns

### Request Configuration
All API queries are defined in [pkg/req/reqs.json](pkg/req/reqs.json). This file is embedded at compile time (`//go:embed`) and parsed as YAML. Each entry specifies:
- `class`: ACI managed object class (e.g., `fvTenant`, `fvBD`)
- `query`: Optional query parameters (e.g., filters, subtree includes)

**When modifying queries:** Update `reqs.json`, then run `python make_script.py` to regenerate the `vetr-collector.sh` shell script alternative.

### Concurrency & Batching
The collector processes requests in parallel batches (default: 7 concurrent requests). See [cmd/collector/main.go#L63-L79](cmd/collector/main.go#L63-L79):
```go
for i := 0; i < len(reqs); i += args.BatchSize {
    var g errgroup.Group
    // Launch batch of requests in parallel
    for j := i; j < i+args.BatchSize && j < len(reqs); j++ {
        g.Go(func() error {
            return cli.Fetch(client, req, arc, cfg)
        })
    }
    err = g.Wait()
}
```

**Pagination:** Large datasets trigger automatic pagination ([cli.go#L109-L169](pkg/cli/cli.go#L109-L169)). When APIC returns "dataset is too big", the collector fetches data in pages (default: 1000 objects/page) and saves as separate JSON files (`class-0.json`, `class-1.json`, etc.).

### Token Management
The ACI client automatically refreshes authentication tokens when >480 seconds old ([client.go#L96-L98](pkg/aci/client.go#L96-L98)). This happens transparently during `client.Do()` calls unless `NoRefresh` modifier is used (only for login/refresh endpoints).

### Error Handling & Retries
Failed requests retry up to 3 times with 10-second delays ([cli.go#L67-L78](pkg/cli/cli.go#L67-L78)). Exception: "dataset is too big" errors immediately trigger pagination instead of retry.

## Development Workflow

### Building & Testing
```bash
# Run from source
go run ./cmd/collector/*.go

# Run tests (uses gock for HTTP mocking)
go test ./...

# Build release binaries (requires goreleaser)
./scripts/release
```

### Release Process
1. Tag version: `git tag v1.2.3`
2. Run `./scripts/release` - this:
   - Runs `python make_script.py` to generate `vetr-collector.sh`
   - Builds cross-platform binaries via goreleaser
   - Packages with README and LICENSE into zip archives

**Note:** `.goreleaser.yml` defines build targets: Windows/Linux/Darwin (arm64 for macOS). CGO is disabled for static binaries.

### Testing Patterns
Tests use [gock](https://github.com/h2non/gock) to mock HTTP responses. See [pkg/aci/client_test.go](pkg/aci/client_test.go):
```go
func testClient() Client {
    client, _ := NewClient(testHost, "usr", "pwd")
    gock.InterceptClient(client.HTTPClient)
    return client
}
```

Always call `defer gock.Off()` to clean up mocks after tests.

## Project-Specific Conventions

### Logging
Uses [zerolog](https://github.com/rs/zerolog) throughout. Log levels in [pkg/log/log.go](pkg/log/log.go):
- `log.Info()` - User-facing progress messages
- `log.Debug()` - Timing/diagnostic info (start/end times)
- `log.Warn()` - Retry attempts, non-fatal issues
- `log.Fatal()` - Unrecoverable errors (exits program)

### File Organization
- **Packages are thin:** Each `pkg/` subdirectory has 2-4 files (implementation + tests)
- **No internal pkg:** All packages are directly under `pkg/`
- **Single binary:** Only one cmd entry point at `cmd/collector/`

### CLI Argument Handling
Uses [go-arg](https://github.com/alexflint/go-arg) for structured CLI parsing. Arguments support environment variables (e.g., `ACI_URL`, `ACI_USERNAME`). Interactive prompts fill missing required values.

**Important:** Passwords with quotes are escaped ([cli.go#L45](pkg/cli/cli.go#L45)): `strings.ReplaceAll(cfg.Password, "\"", "\\\"")` to handle special characters in APIC passwords.

## External Dependencies

- **tidwall/gjson & sjson** - Fast JSON parsing/building without struct marshaling
- **golang.org/x/sync/errgroup** - Parallel error handling for batched requests
- **alexflint/go-arg** - CLI argument parsing with struct tags
- **rs/zerolog** - Structured logging
- **h2non/gock** - HTTP mocking for tests

## Common Gotchas

1. **Archive writes must be thread-safe:** Use `zipMux.Lock()` in [archive.go#L44](pkg/archive/archive.go#L44) since parallel goroutines write to the same zip file.

2. **URL normalization:** User input is stripped of `http://` and `https://` prefixes ([args.go#L63-L64](cmd/collector/args.go#L63-L64)), then `https://` is re-added in `aci.NewClient`.

3. **Version injection:** The `version` variable in [main.go](cmd/collector/main.go) is set via `-ldflags` during build: `-X main.version=$TAG`.

4. **Dual collection methods:** Binary collector (this codebase) and shell script (`vetr-collector.sh`) must stay in sync. Always run `make_script.py` after modifying `reqs.json`.
