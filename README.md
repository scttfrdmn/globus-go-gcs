# Globus Connect Server Go Port

[![Go Report Card](https://goreportcard.com/badge/github.com/scttfrdmn/globus-go-gcs)](https://goreportcard.com/report/github.com/scttfrdmn/globus-go-gcs)
[![CI](https://github.com/scttfrdmn/globus-go-gcs/workflows/CI/badge.svg)](https://github.com/scttfrdmn/globus-go-gcs/actions)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/scttfrdmn/globus-go-gcs)](go.mod)

A complete Go port of Globus Connect Server v5 CLI - a drop-in replacement for the Python `globus-connect-server` command with 100% feature parity.

> **STATUS**: ✅ **v1.0 Complete** - All 83 commands implemented | 🚧 **v2.0 In Planning** - Security & Performance
>
> - **v1.0**: Feature parity achieved (January 2025)
> - **v2.0**: HIPAA compliance + 2-5x performance improvements (Planning phase)
> - See [ROADMAP.md](ROADMAP.md) for v2.0 plans

---

## ⚠️ IMPORTANT: v2.0 Breaking Change

**If you plan to upgrade from v1.x to v2.0**, please read **[BREAKING_CHANGES_V2.md](./BREAKING_CHANGES_V2.md)** first.

### What's Changing

For **security compliance (HIPAA/PHI)**, v2.0 removes the ability to pass secrets as command-line arguments.

**Affected commands** (5):
- `user-credential s3-keys add/update`
- `user-credential activescale-create`
- `oidc create/update`

**New secure methods**:
- ✅ Interactive prompt (recommended)
- ✅ Stdin pipe with `--secret-stdin`
- ✅ Environment variable with `--secret-env`

### Why This Change Is Necessary

```bash
# v1.x - INSECURE (visible in process list and shell history)
globus-connect-server user-credential s3-keys add \
  --secret-access-key wJalrXUtnFEMI/K7MDENG...  # ❌ SECURITY RISK

# Any user can see:
$ ps aux | grep secret
user  12345  ... --secret-access-key wJalrXUtnFEMI/K7MDENG...

# v2.0 - SECURE (not visible)
echo "wJalrXUtnFEMI/K7MDENG..." | \
  globus-connect-server user-credential s3-keys add \
    --access-key-id AKIA... \
    --secret-stdin  # ✅ SECURE
```

**Security violations in v1.x**:
- ❌ NIST 800-53 IA-5(7): Prohibits embedded unprotected passwords
- ❌ HIPAA Security Rule § 164.312(a)(2)(iv): Requires encryption of credentials
- ❌ Secrets visible in: process listings, shell history, system logs

**Read the complete migration guide**: [BREAKING_CHANGES_V2.md](./BREAKING_CHANGES_V2.md)

---

## Overview

**globus-go-gcs** is a command-line tool for managing Globus Connect Server v5 endpoints, written in Go. It provides complete feature parity with the official Python `globus-connect-server` CLI while offering the benefits of a compiled, single-binary distribution.

### Why Go?

- ✅ **Single Binary**: No Python runtime or dependencies required
- ✅ **Cross-Platform**: Easy ARM64 support, works on all major platforms
- ✅ **Performance**: Faster startup time, lower memory footprint
- ✅ **Maintainability**: Strong typing, better tooling, clear dependency management
- ✅ **Code Quality**: Idiomatic Go practices, A+ Go Report Card grade

### Project Goals

1. **100% Feature Parity**: All commands, flags, and config files from Python version
2. **Drop-in Replacement**: Same command structure, compatible configuration files
3. **Config Compatibility**: Reads and writes same config files as Python CLI
4. **Better UX**: Improved error messages, tab completion, interactive prompts
5. **Idiomatic Go**: Follow Go best practices, achieve **A+ Go Report Card grade**

## Features (Roadmap)

### Phase 1: Proof of Concept (In Development)
- [x] Project setup and infrastructure
- [ ] OAuth2 authentication (login/logout/whoami)
- [ ] Token storage compatible with Python CLI
- [ ] Configuration file management
- [ ] Basic endpoint commands (`endpoint show`)
- [ ] JSON output support

### Phase 2: Core Management
- [ ] Complete endpoint lifecycle management
- [ ] Data transfer node management
- [ ] Role-based access control (RBAC)

### Phase 3: Collections & Storage
- [ ] Collection management (mapped & guest)
- [ ] Storage gateway management
- [ ] Storage connector support (POSIX, S3, GCS, Azure)

### Phase 4: Advanced Features
- [ ] OIDC server management
- [ ] Sharing policies
- [ ] Authentication policies
- [ ] User credential management
- [ ] Audit log search
- [ ] Self-diagnostic utilities

### Phase 5: Polish & Release
- [ ] Multi-platform binaries
- [ ] Complete documentation
- [ ] Migration guides
- [ ] v1.0.0 release

See [PROJECT_PLAN.md](PROJECT_PLAN.md) for detailed implementation plan.

## Installation

### From Source (Development)

```bash
# Clone repository
git clone https://github.com/scttfrdmn/globus-go-gcs.git
cd globus-go-gcs

# Build
go build -o globus-connect-server ./cmd/globus-connect-server

# Install
sudo mv globus-connect-server /usr/local/bin/
```

### Binary Releases (Coming Soon)

Pre-built binaries for Linux, macOS, and Windows will be available once Phase 1 is complete.

## Quick Start

```bash
# Authenticate
globus-connect-server login

# Show endpoint information
globus-connect-server endpoint show

# List collections
globus-connect-server collection list

# Get help
globus-connect-server --help
globus-connect-server endpoint --help
```

## Command Structure

The CLI follows the same structure as the Python version:

```bash
# Authentication & Session
globus-connect-server login
globus-connect-server logout
globus-connect-server whoami
globus-connect-server session show

# Endpoint Management
globus-connect-server endpoint setup <name>
globus-connect-server endpoint show
globus-connect-server endpoint update <field> <value>
globus-connect-server endpoint role list

# Node Management
globus-connect-server node setup <name>
globus-connect-server node list
globus-connect-server node show <id>

# Collection Management
globus-connect-server collection create <storage-gateway-id> <path>
globus-connect-server collection list
globus-connect-server collection show <id>

# Storage Gateway Management
globus-connect-server storage-gateway create <type> <name>
globus-connect-server storage-gateway list
globus-connect-server storage-gateway show <id>

# ... and more
```

See [docs/COMMAND_REFERENCE.md](docs/COMMAND_REFERENCE.md) (coming soon) for complete command documentation.

## Configuration

The Go CLI uses the same configuration files as the Python version:

- `~/.globus-connect-server/config.json` - Main configuration
- `~/.globus-connect-server/tokens.json` - OAuth tokens
- `/var/lib/globus-connect-server/info.json` - Local endpoint state

This ensures compatibility - you can switch between Python and Go CLIs seamlessly.

## Documentation

- [PROJECT_PLAN.md](PROJECT_PLAN.md) - Complete project plan and roadmap
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines (coming soon)
- [docs/MIGRATION_FROM_PYTHON.md](docs/MIGRATION_FROM_PYTHON.md) - Migration guide (coming soon)

## Architecture

```
globus-go-gcs/
├── cmd/
│   └── globus-connect-server/    # CLI entry point
├── pkg/
│   ├── client/gcs/                # GCS Manager API client
│   ├── config/                    # Configuration management
│   ├── models/                    # Data structures
│   └── output/                    # Output formatting
└── internal/
    └── commands/                  # Command implementations
```

### GCS Manager API Client

Unlike other Globus services, GCS Manager API runs on individual endpoints (not a central API). We implement a custom GCS Manager API client based on the Python SDK v4 `GCSClient`:

- Connects to GCS endpoints by FQDN (e.g., `https://gcs.university.edu/api`)
- Manages endpoints, collections, storage gateways, roles, user credentials
- Built on top of [globus-go-sdk](https://github.com/scttfrdmn/globus-go-sdk) v3

## Development

### Prerequisites

- Go 1.21 or higher
- Access to a Globus Connect Server endpoint for testing

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run linter
make lint
```

### Testing

```bash
# Unit tests
go test ./...

# Integration tests (requires GCS endpoint)
go test -tags=integration ./...

# With coverage
go test -cover ./...
```

## Code Quality Standards

This project is committed to **idiomatic Go practices** and achieving an **A+ grade on [Go Report Card](https://goreportcard.com/)**.

### Quality Checklist

Before every commit, ensure:

```bash
# 1. Format code
make fmt

# 2. Run go vet
make vet

# 3. Run linter
make lint

# 4. Run tests
make test

# Or run all checks at once
make verify
```

### Go Report Card Criteria

We maintain **100% compliance** with:

- ✅ **gofmt**: All code formatted with `gofmt`
- ✅ **go vet**: Zero issues from `go vet`
- ✅ **golint/staticcheck**: Zero lint warnings
- ✅ **gocyclo**: Cyclomatic complexity < 15 per function
- ✅ **ineffassign**: No ineffectual assignments
- ✅ **misspell**: No spelling errors

### Idiomatic Go Practices

We follow established Go conventions:

- Clear, focused packages with single responsibility
- Accept interfaces, return structs
- Context-first parameter ordering
- Comprehensive error handling with error wrapping
- Table-driven tests with good coverage
- Complete documentation for all exported symbols

See [CODE_STANDARDS.md](CODE_STANDARDS.md) for detailed guidelines and examples.

### CI/CD Quality Gates

Our CI pipeline enforces:

- Multi-OS testing (Ubuntu, macOS)
- Multi-Go version support (1.21, 1.22, 1.23)
- Comprehensive linting with golangci-lint (30+ linters enabled)
- 80%+ test coverage requirement
- Zero security issues (gosec)

## Project Management

This project follows the **lens-style project management methodology** with:

- **4 User Personas**: System Admin, Data Manager, Research PI, IT Manager
- **6 Development Phases**: Setup → POC → Core → Collections → Advanced → Release
- **GitHub Project Board**: Track progress with custom fields (persona, phase, priority)
- **147 Labels**: Organized by type, priority, area, persona, phase, status
- **Issue Templates**: Feature requests, bug reports with persona fields
- **Milestones**: One per phase with clear success criteria

See [.github/PROJECT_BOARD_SETUP.md](.github/PROJECT_BOARD_SETUP.md) (coming soon) for details.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) (coming soon) for guidelines.

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.

## Disclaimer

This is an independent, community-developed project and is not officially affiliated with, endorsed by, or supported by Globus, the University of Chicago, or their affiliated organizations.

## Acknowledgments

- Built on [globus-go-sdk](https://github.com/scttfrdmn/globus-go-sdk)
- Inspired by [globus-go-cli](https://github.com/scttfrdmn/globus-go-cli)
- Project management patterns from [lens](https://github.com/scttfrdmn/lens)
- Based on official Globus Connect Server v5 Python CLI

---

**Project Status**: Phase 0 - Project Setup
**Next Milestone**: [Phase 1] Proof of Concept - Authentication Working
**Last Updated**: October 2025
