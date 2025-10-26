# Globus Connect Server Go Port - Project Plan

## Executive Summary

**Project**: globus-go-gcs - Complete Go port of Globus Connect Server v5 CLI
**Goal**: Drop-in replacement for Python `globus-connect-server` CLI with 100% feature parity
**Approach**: Phased implementation using lens-style project management methodology
**Timeline**: 6-9 months for full parity, 2-3 months for Phase 1 proof-of-concept

## Project Context

### Upstream Projects
- **Python GCS CLI**: Official Globus Connect Server v5 command-line tool
- **globus-go-sdk**: Go SDK for Globus APIs (adjacent project, provides API client foundation)
- **globus-go-cli**: Go CLI for Globus services (adjacent project, provides CLI patterns)
- **lens**: Reference project for GitHub project management methodology

### Key Requirements
1. **Complete Feature Parity**: All commands, flags, config files from Python version
2. **Drop-in Replacement**: Same command structure, compatible configuration
3. **Phased Implementation**: Proof-of-concept → Core features → Advanced features → Full parity
4. **Lens-Style Project Management**: Persona-driven, milestone-tracked, issue-based planning

---

## Architecture Overview

### Project Structure (Following globus-go-sdk and globus-go-cli patterns)

```
globus-go-gcs/
├── cmd/
│   └── globus-connect-server/    # Main CLI entry point
│       └── main.go
├── internal/
│   └── commands/                  # Command implementations
│       ├── endpoint/              # Endpoint management commands
│       ├── node/                  # Data transfer node commands
│       ├── collection/            # Collection management
│       ├── storage_gateway/       # Storage gateway commands
│       ├── auth/                  # Authentication commands (login/logout)
│       ├── session/               # Session management
│       ├── oidc/                  # OIDC server management
│       ├── sharing_policy/        # Sharing policies
│       ├── auth_policy/           # Auth policies
│       ├── user_credentials/      # User credential management
│       ├── audit/                 # Audit log search
│       └── self_diagnostic/       # Diagnostic utilities
├── pkg/
│   ├── client/                    # GCS Manager API client (extends globus-go-sdk)
│   │   ├── endpoint.go
│   │   ├── node.go
│   │   ├── collection.go
│   │   ├── storage_gateway.go
│   │   └── ...
│   ├── config/                    # Configuration management
│   │   ├── config.go              # Read/write ~/.globus-connect-server/
│   │   └── session.go             # Session token management
│   ├── models/                    # Data structures
│   │   ├── endpoint.go
│   │   ├── collection.go
│   │   └── ...
│   └── output/                    # Output formatting (text, JSON, CSV)
│       ├── formatter.go
│       └── table.go
├── .github/
│   ├── ISSUE_TEMPLATE/            # Issue templates with persona fields
│   │   ├── feature_request.yml
│   │   ├── bug_report.yml
│   │   ├── technical_debt.yml
│   │   └── documentation.yml
│   ├── workflows/                 # CI/CD workflows
│   │   ├── ci.yml
│   │   ├── release.yml
│   │   └── labels.yml
│   ├── labels.yml                 # Label definitions (persona, phase, area)
│   ├── pull_request_template.md   # PR template with persona impact
│   ├── PROJECT_BOARD_SETUP.md     # Project board documentation
│   └── GITHUB_ISSUES_SUMMARY.md   # Issue tracking summary
├── docs/
│   ├── USER_REQUIREMENTS.md       # Requirements with success metrics
│   ├── DESIGN_PRINCIPLES.md       # Architectural decisions
│   ├── USER_SCENARIOS/            # Persona walkthroughs
│   │   ├── 01_SYSTEM_ADMIN_WALKTHROUGH.md
│   │   ├── 02_DATA_MANAGER_WALKTHROUGH.md
│   │   ├── 03_RESEARCH_PI_WALKTHROUGH.md
│   │   └── 04_IT_MANAGER_WALKTHROUGH.md
│   ├── MIGRATION_FROM_PYTHON.md   # Python → Go migration guide
│   └── COMMAND_REFERENCE.md       # Complete command documentation
├── go.mod
├── go.sum
├── README.md
├── ROADMAP.md                     # Phased development roadmap
├── CHANGELOG.md
└── LICENSE

```

### Technology Stack

- **CLI Framework**: [Cobra](https://github.com/spf13/cobra) (same as globus-go-cli)
- **Configuration**: [Viper](https://github.com/spf13/viper) for config management
- **Base SDK**: globus-go-sdk v3.65.0-1+ (current stable)
  - Provides: Auth client, OAuth2, token management, transport layer
  - **Note**: GCS Manager API client NOT included - we'll implement it
- **GCS Manager API Client**: Custom implementation in `pkg/client/gcs/`
  - Based on Python SDK v4 GCSClient patterns
  - Connects to individual GCS endpoints by FQDN (not central API)
  - Will be contributed back to globus-go-sdk after proven
- **OAuth2**: golang.org/x/oauth2 for authentication
- **Output Formatting**: text/template for tables, encoding/json for JSON
- **Testing**: Go standard testing + testify for assertions

### SDK Version Strategy

**Current Approach: Build on v3, Prepare for v4**

1. **Phase 1-4**: Use globus-go-sdk **v3.65.0-1** (stable)
   - Implement GCS Manager API client in this project
   - CLI works with proven, stable foundation
   - Faster time to market

2. **Phase 5**: Prepare for globus-go-sdk **v4** migration
   - globus-go-sdk v4 targeting December 2025
   - Create migration path when v4 stable
   - Use Go module versioning (v2 of globus-go-gcs)

3. **Future**: Contribute back to globus-go-sdk
   - Once GCS Manager client proven in production
   - Add as `pkg/services/gcs/` in globus-go-sdk v4+
   - Benefit entire Globus Go community

**Why Not Wait for v4?**
- v4 won't be stable until Q4 2025 (delays project by 6+ months)
- GCS Manager API client not in v4 plan anyway
- v3 is mature, production-ready, well-documented
- Can migrate to v4 later using module versioning

## GCS Manager API Client Implementation

### Overview

Unlike other Globus services (Transfer, Auth, Groups, etc.) that use centralized APIs, **GCS Manager API runs on individual GCS endpoints**. This means:

- **No central API**: Each GCS endpoint has its own API at `https://<gcs-endpoint-fqdn>/api`
- **Endpoint-specific**: Must connect to specific FQDN (e.g., `gcs.university.edu`)
- **Python SDK v4 has it**: `GCSClient` class in globus-sdk-python
- **Go SDK missing**: Not in globus-go-sdk v3 or v4 plan - we'll implement it

### API Client Architecture

We'll implement the GCS Manager API client in `pkg/client/gcs/`:

```go
pkg/client/gcs/
├── client.go              # GCS Manager API base client
├── endpoint.go            # Endpoint operations
├── collection.go          # Collection CRUD
├── storage_gateway.go     # Storage gateway CRUD
├── role.go                # Role-based access control
├── user_credential.go     # User credential management
├── models.go              # Data structures
└── scopes.go              # Scope helpers
```

### GCS Manager API Methods

Based on Python SDK v4 GCSClient, we need to implement:

#### Endpoint Management
```go
// Get GCS info (unauthenticated)
func (c *Client) GetGCSInfo(ctx context.Context) (*GCSInfo, error)

// Get endpoint details
func (c *Client) GetEndpoint(ctx context.Context) (*Endpoint, error)

// Update endpoint
func (c *Client) UpdateEndpoint(ctx context.Context, update *EndpointUpdate) (*Endpoint, error)
```

#### Collection Management
```go
// List collections with optional filters
func (c *Client) ListCollections(ctx context.Context, opts *ListCollectionsOptions) (*CollectionList, error)

// Get collection by ID
func (c *Client) GetCollection(ctx context.Context, collectionID string) (*Collection, error)

// Create mapped or guest collection
func (c *Client) CreateCollection(ctx context.Context, collection *Collection) (*Collection, error)

// Update collection
func (c *Client) UpdateCollection(ctx context.Context, collectionID string, update *CollectionUpdate) (*Collection, error)

// Delete collection
func (c *Client) DeleteCollection(ctx context.Context, collectionID string) error
```

#### Storage Gateway Management
```go
// List storage gateways
func (c *Client) ListStorageGateways(ctx context.Context) (*StorageGatewayList, error)

// Get storage gateway by ID
func (c *Client) GetStorageGateway(ctx context.Context, gatewayID string) (*StorageGateway, error)

// Create storage gateway (POSIX, S3, etc.)
func (c *Client) CreateStorageGateway(ctx context.Context, gateway *StorageGateway) (*StorageGateway, error)

// Update storage gateway
func (c *Client) UpdateStorageGateway(ctx context.Context, gatewayID string, update *StorageGatewayUpdate) (*StorageGateway, error)

// Delete storage gateway
func (c *Client) DeleteStorageGateway(ctx context.Context, gatewayID string) error
```

#### Role Management (RBAC)
```go
// List roles for endpoint or collection
func (c *Client) ListRoles(ctx context.Context, opts *ListRolesOptions) (*RoleList, error)

// Get role by ID
func (c *Client) GetRole(ctx context.Context, roleID string) (*Role, error)

// Create role (assign to principal)
func (c *Client) CreateRole(ctx context.Context, role *Role) (*Role, error)

// Delete role
func (c *Client) DeleteRole(ctx context.Context, roleID string) error
```

#### User Credential Management
```go
// List user credentials
func (c *Client) ListUserCredentials(ctx context.Context) (*UserCredentialList, error)

// Get user credential by ID
func (c *Client) GetUserCredential(ctx context.Context, credentialID string) (*UserCredential, error)

// Create user credential (S3, OAuth, Activescale, etc.)
func (c *Client) CreateUserCredential(ctx context.Context, credential *UserCredential) (*UserCredential, error)

// Update user credential
func (c *Client) UpdateUserCredential(ctx context.Context, credentialID string, update *UserCredentialUpdate) (*UserCredential, error)

// Delete user credential
func (c *Client) DeleteUserCredential(ctx context.Context, credentialID string) error
```

#### Scope Helpers
```go
// Get GCS endpoint scopes
func GetGCSEndpointScopes(endpointID string) []string

// Get GCS collection scopes
func GetGCSCollectionScopes(endpointID, collectionID string) []string
```

### Client Initialization

```go
package gcs

import (
    "github.com/scttfrdmn/globus-go-sdk/v3/pkg/core/auth"
    "github.com/scttfrdmn/globus-go-sdk/v3/pkg/core/transport"
)

// Client represents a GCS Manager API client
type Client struct {
    baseURL    string // https://<gcs-endpoint-fqdn>/api
    authorizer auth.Authorizer
    transport  *transport.HTTPTransport
}

// NewClient creates a new GCS Manager API client
func NewClient(gcsAddress string, authorizer auth.Authorizer, opts ...ClientOption) (*Client, error) {
    // Parse address (FQDN or full URL)
    baseURL := parseGCSAddress(gcsAddress) // Add /api if needed

    client := &Client{
        baseURL:    baseURL,
        authorizer: authorizer,
        transport:  transport.NewHTTPTransport(),
    }

    // Apply options
    for _, opt := range opts {
        opt(client)
    }

    return client, nil
}

// parseGCSAddress converts FQDN or URL to proper API base URL
func parseGCSAddress(address string) string {
    // If it's just a FQDN like "gcs.example.edu"
    if !strings.HasPrefix(address, "http") {
        address = "https://" + address
    }

    // If URL doesn't end with /api, append it
    if !strings.HasSuffix(address, "/api") {
        address = strings.TrimSuffix(address, "/") + "/api"
    }

    return address
}
```

### Example Usage

```go
import (
    "github.com/scttfrdmn/globus-go-gcs/pkg/client/gcs"
    "github.com/scttfrdmn/globus-go-sdk/v3/pkg/services/auth"
)

// Get access token from globus-go-sdk
authorizer := auth.NewStaticTokenAuthorizer("your-access-token")

// Create GCS client for specific endpoint
gcsClient, err := gcs.NewClient("gcs.university.edu", authorizer)
if err != nil {
    log.Fatal(err)
}

// Get endpoint info
endpoint, err := gcsClient.GetEndpoint(ctx)
if err != nil {
    log.Fatal(err)
}

// List collections
collections, err := gcsClient.ListCollections(ctx, nil)
if err != nil {
    log.Fatal(err)
}
```

### Implementation Priority

1. **Phase 1**: Core client + endpoint operations (proof of concept)
2. **Phase 2**: Storage gateway + collection management
3. **Phase 3**: Role management (RBAC)
4. **Phase 4**: User credentials + advanced features

---

## Persona Definitions

### Primary Personas (GCS Users)

#### 1. System Administrator
**Profile:**
- Name: Alex Chen
- Role: Linux Systems Administrator at mid-size university
- Technical Level: 5/5 (expert)
- Responsibility: Deploy and maintain GCS endpoints for research groups

**Pain Points:**
- Manual Python CLI installation across multiple servers
- Dependency conflicts with Python environments
- Need consistent deployment across CentOS, Ubuntu, Rocky Linux
- Config file management across multiple endpoints

**Success Metrics:**
- Single binary deployment (no Python dependencies)
- Config file 100% compatible with Python version
- Cross-platform builds (Linux x86_64, ARM64)
- Installation time < 5 minutes (vs 20 minutes with Python)

#### 2. Data Manager
**Profile:**
- Name: Dr. Sarah Johnson
- Role: Research Data Manager for large research project
- Technical Level: 3/5 (proficient)
- Responsibility: Manage collections, sharing policies, access controls

**Pain Points:**
- Complex collection setup commands
- Remembering exact syntax for storage gateway configs
- Managing multiple mapped collections
- Setting up guest collections for collaborators

**Success Metrics:**
- Interactive wizards for complex operations
- Clear error messages with actionable suggestions
- Consistent command structure (easier to remember)
- Tab completion for all commands

#### 3. Research PI
**Profile:**
- Name: Prof. Michael Rodriguez
- Role: Principal Investigator managing lab data sharing
- Technical Level: 2/5 (basic command line skills)
- Responsibility: Share datasets with collaborators, manage access

**Pain Points:**
- CLI intimidating, prefers GUI but needs automation
- Forgets commands between uses
- Needs to delegate management to grad students
- Compliance requirements for sensitive data

**Success Metrics:**
- Self-documenting commands (built-in help)
- Audit logging for compliance
- Delegation with role-based access
- Examples in help text for common tasks

#### 4. IT Manager
**Profile:**
- Name: Jennifer Martinez
- Role: Research Computing IT Manager
- Technical Level: 4/5 (very proficient)
- Responsibility: Oversee 20+ GCS endpoints across institution

**Pain Points:**
- Need automation for monitoring endpoints
- JSON output for scripting
- Bulk operations across multiple endpoints
- Cost tracking and resource management

**Success Metrics:**
- JSON output for all commands
- Scriptable without interactive prompts
- Batch operations support
- Consistent exit codes for automation

---

## Command Structure (Python GCS CLI Parity)

### Core Commands (Must-Have for MVP)

```bash
# Authentication & Session
globus-connect-server login                          # Interactive OAuth login
globus-connect-server logout                         # Clear local tokens
globus-connect-server whoami                         # Show current user
globus-connect-server session show                   # Display session info
globus-connect-server session update                 # Update session consent

# Endpoint Management
globus-connect-server endpoint setup <name>          # Create new endpoint
globus-connect-server endpoint show                  # Show endpoint details
globus-connect-server endpoint update <field> <value># Update endpoint config
globus-connect-server endpoint upgrade               # Upgrade endpoint
globus-connect-server endpoint role list             # List endpoint roles
globus-connect-server endpoint role create           # Assign role to user

# Node Management
globus-connect-server node setup <name>              # Create data transfer node
globus-connect-server node show <node-id>            # Show node details
globus-connect-server node list                      # List all nodes
globus-connect-server node enable <node-id>          # Enable node
globus-connect-server node disable <node-id>         # Disable node
globus-connect-server node cleanup <node-id>         # Remove node

# Collection Management
globus-connect-server collection create <storage-gateway-id> <path> # Create mapped collection
globus-connect-server collection list                # List collections
globus-connect-server collection show <id>           # Show collection details
globus-connect-server collection update <id> <field> <value> # Update collection
globus-connect-server collection delete <id>         # Delete collection

# Storage Gateway Management
globus-connect-server storage-gateway create <type> <name> # Create storage gateway
globus-connect-server storage-gateway list           # List storage gateways
globus-connect-server storage-gateway show <id>      # Show gateway details
globus-connect-server storage-gateway update <id> <field> <value> # Update gateway
globus-connect-server storage-gateway delete <id>    # Delete gateway
```

### Advanced Commands (Phase 2+)

```bash
# OIDC Server Management (added in v5.4.14)
globus-connect-server oidc enable                    # Enable OIDC server
globus-connect-server oidc show                      # Show OIDC config
globus-connect-server oidc disable                   # Disable OIDC server

# Sharing Policies (added in v5.4.17)
globus-connect-server sharing-policy create          # Create sharing policy
globus-connect-server sharing-policy list            # List policies
globus-connect-server sharing-policy show <id>       # Show policy details
globus-connect-server sharing-policy update <id>     # Update policy
globus-connect-server sharing-policy delete <id>     # Delete policy

# Auth Policies (added in v5.4.57)
globus-connect-server auth-policy create             # Create auth policy
globus-connect-server auth-policy list               # List policies
globus-connect-server auth-policy show <id>          # Show policy
globus-connect-server auth-policy update <id>        # Update policy
globus-connect-server auth-policy delete <id>        # Delete policy

# User Credentials (added in v5.4.51)
globus-connect-server user-credentials create <type> # Create credential (S3, OAuth, etc.)
globus-connect-server user-credentials list          # List credentials
globus-connect-server user-credentials show <id>     # Show credential
globus-connect-server user-credentials update <id>   # Update credential
globus-connect-server user-credentials delete <id>   # Delete credential

# Audit Logs (High Assurance)
globus-connect-server audit search <query>           # Search audit logs
globus-connect-server audit export <query>           # Export audit logs

# Diagnostics
globus-connect-server self-diagnostic run            # Run diagnostic tests
globus-connect-server self-diagnostic show           # Show diagnostic results
```

### Configuration Files (Must Be Compatible)

The Go CLI must read and write the same config files as the Python CLI:

- `~/.globus-connect-server/config.json` - Main configuration
- `~/.globus-connect-server/tokens.json` - OAuth tokens
- `/var/lib/globus-connect-server/info.json` - Local endpoint state

---

## Phased Implementation Plan

### Phase 0: Project Setup & Infrastructure (Weeks 1-2)

**Goal**: Establish project structure, GitHub project management, CI/CD

#### Deliverables:
- [ ] GitHub repository created with README, LICENSE
- [ ] Project structure following globus-go-sdk/cli patterns
- [ ] go.mod initialized with dependencies (Cobra, Viper, globus-go-sdk)
- [ ] Makefile for build, test, install
- [ ] CI/CD workflows (GitHub Actions)
  - [ ] Test on push/PR
  - [ ] Lint with golangci-lint
  - [ ] Build for multiple platforms (Linux x86_64, ARM64, macOS)
- [ ] GitHub Project Board created with custom fields:
  - [ ] Persona (single select: System Admin, Data Manager, Research PI, IT Manager)
  - [ ] Phase (single select: Phase 0, Phase 1, Phase 2, Phase 3, Phase 4)
  - [ ] Estimate (number: story points or days)
  - [ ] Priority (single select: Critical, High, Medium, Low)
- [ ] Label system configured (147 labels following lens pattern):
  - [ ] Type labels (bug, enhancement, documentation, technical-debt)
  - [ ] Priority labels (critical, high, medium, low)
  - [ ] Area labels (cli, config, auth, endpoint, collection, etc.)
  - [ ] Persona labels (persona: system-admin, persona: data-manager, etc.)
  - [ ] Phase labels (phase: 0-setup, phase: 1-mvp, phase: 2-core, etc.)
  - [ ] Status labels (triage, ready, in-progress, in-review, blocked)
- [ ] Issue templates created:
  - [ ] Feature request with persona and phase fields
  - [ ] Bug report with persona field
  - [ ] Technical debt template
  - [ ] Documentation template
- [ ] PR template with persona impact section

**Milestone**: [Phase 0] Project Infrastructure Complete

#### Issues to Create:
1. #1 - [Phase 0] Initialize Go module and project structure
2. #2 - [Phase 0] Set up GitHub Actions CI/CD pipeline
3. #3 - [Phase 0] Configure GitHub Project Board with custom fields
4. #4 - [Phase 0] Create issue and PR templates
5. #5 - [Phase 0] Write PROJECT_PLAN.md and ROADMAP.md
6. #6 - [Phase 0] Create label system following lens pattern

---

### Phase 1: Proof of Concept - Core Authentication & Basic Commands (Weeks 3-6)

**Goal**: Demonstrate feasibility with working authentication and 1-2 core commands

#### Deliverables:
- [ ] OAuth2 authentication flow working
  - [ ] `globus-connect-server login` (interactive browser flow)
  - [ ] `globus-connect-server logout`
  - [ ] `globus-connect-server whoami`
  - [ ] `globus-connect-server session show`
- [ ] Token storage compatible with Python CLI
  - [ ] Read existing `~/.globus-connect-server/tokens.json`
  - [ ] Write tokens in same format
- [ ] Config management
  - [ ] Read `~/.globus-connect-server/config.json`
  - [ ] Write config in same format
- [ ] GCS Manager API client (pkg/client/)
  - [ ] Base client with authentication
  - [ ] Error handling
  - [ ] Retry logic with exponential backoff
- [ ] Basic endpoint commands (proof of concept)
  - [ ] `globus-connect-server endpoint show`
  - [ ] `globus-connect-server endpoint list` (if multi-endpoint)
- [ ] Output formatting
  - [ ] Text output (human-readable tables)
  - [ ] JSON output (--format=json flag)
- [ ] Comprehensive tests
  - [ ] Unit tests for all packages
  - [ ] Integration tests with mock API server
- [ ] Documentation
  - [ ] README with installation and quick start
  - [ ] MIGRATION_FROM_PYTHON.md guide

**Milestone**: [Phase 1] Proof of Concept - Authentication Working

#### Issues to Create:
7. #7 - [Phase 1][System Admin] Implement OAuth2 login flow
8. #8 - [Phase 1][System Admin] Implement token storage compatible with Python CLI
9. #9 - [Phase 1][System Admin] Implement config file management
10. #10 - [Phase 1][All Personas] Create GCS Manager API base client
11. #11 - [Phase 1][System Admin] Implement `endpoint show` command
12. #12 - [Phase 1][IT Manager] Add JSON output formatting
13. #13 - [Phase 1][System Admin] Write integration tests for auth flow
14. #14 - [Phase 1][All Personas] Write migration guide from Python CLI

**Success Criteria:**
- Can authenticate using OAuth2 flow
- Can read tokens from Python CLI installation
- Can execute `endpoint show` successfully
- All commands output valid JSON with --format=json
- 80%+ test coverage

---

### Phase 2: Core Endpoint & Node Management (Weeks 7-12)

**Goal**: Complete endpoint and node lifecycle management (setup, update, delete)

#### Deliverables:
- [ ] Complete endpoint management
  - [ ] `endpoint setup` - Create new endpoint
  - [ ] `endpoint update` - Update endpoint configuration
  - [ ] `endpoint upgrade` - Upgrade endpoint version
  - [ ] `endpoint role list` - List roles
  - [ ] `endpoint role create` - Assign roles
  - [ ] `endpoint role delete` - Remove roles
- [ ] Complete node management
  - [ ] `node setup` - Create data transfer node
  - [ ] `node list` - List all nodes
  - [ ] `node show` - Show node details
  - [ ] `node update` - Update node configuration
  - [ ] `node enable` - Enable node
  - [ ] `node disable` - Disable node
  - [ ] `node cleanup` - Remove node
  - [ ] `node secrets` - Manage node secrets
- [ ] API client extensions
  - [ ] pkg/client/endpoint.go - Endpoint operations
  - [ ] pkg/client/node.go - Node operations
  - [ ] pkg/models/endpoint.go - Data structures
  - [ ] pkg/models/node.go - Data structures
- [ ] Interactive prompts for complex operations
  - [ ] Confirm destructive operations (delete, cleanup)
  - [ ] Validation of inputs
  - [ ] Helpful error messages
- [ ] Enhanced testing
  - [ ] Unit tests for all commands
  - [ ] Integration tests against staging API
  - [ ] End-to-end tests with real endpoint creation

**Milestone**: [Phase 2] Endpoint & Node Management Complete

#### Issues to Create:
15. #15 - [Phase 2][System Admin] Implement `endpoint setup` command
16. #16 - [Phase 2][System Admin] Implement `endpoint update` command
17. #17 - [Phase 2][System Admin] Implement endpoint role management commands
18. #18 - [Phase 2][System Admin] Implement `node setup` command
19. #19 - [Phase 2][System Admin] Implement complete node lifecycle commands
20. #20 - [Phase 2][System Admin] Add interactive prompts for destructive operations
21. #21 - [Phase 2][Data Manager] Add validation and helpful error messages
22. #22 - [Phase 2][System Admin] Write end-to-end tests for endpoint creation

**Success Criteria:**
- Can create endpoint from scratch
- Can create and manage data transfer nodes
- Can assign roles to users
- Config files remain compatible with Python CLI
- 80%+ test coverage maintained

---

### Phase 3: Collection & Storage Gateway Management (Weeks 13-18)

**Goal**: Complete collection and storage gateway operations

#### Deliverables:
- [ ] Storage gateway management
  - [ ] `storage-gateway create` - Create gateway (POSIX, S3, etc.)
  - [ ] `storage-gateway list` - List all gateways
  - [ ] `storage-gateway show` - Show gateway details
  - [ ] `storage-gateway update` - Update gateway config
  - [ ] `storage-gateway delete` - Delete gateway
- [ ] Collection management
  - [ ] `collection create` - Create mapped collection
  - [ ] `collection list` - List all collections
  - [ ] `collection show` - Show collection details
  - [ ] `collection update` - Update collection settings
  - [ ] `collection delete` - Delete collection
- [ ] Storage connector support
  - [ ] POSIX filesystem connector
  - [ ] S3 connector
  - [ ] Google Cloud Storage connector
  - [ ] Azure Blob connector
  - [ ] Additional connectors as documented in Python SDK
- [ ] Guest collection support
  - [ ] Create guest collections
  - [ ] Manage guest collection permissions
- [ ] API client extensions
  - [ ] pkg/client/storage_gateway.go
  - [ ] pkg/client/collection.go
  - [ ] pkg/models/storage_gateway.go
  - [ ] pkg/models/collection.go

**Milestone**: [Phase 3] Collection & Storage Management Complete

#### Issues to Create:
23. #23 - [Phase 3][Data Manager] Implement storage gateway commands
24. #24 - [Phase 3][Data Manager] Implement collection commands
25. #25 - [Phase 3][System Admin] Add POSIX storage connector support
26. #26 - [Phase 3][System Admin] Add S3 storage connector support
27. #27 - [Phase 3][System Admin] Add cloud storage connector support (GCS, Azure)
28. #28 - [Phase 3][Research PI] Implement guest collection management
29. #29 - [Phase 3][Data Manager] Add collection sharing permissions

**Success Criteria:**
- Can create and manage storage gateways
- Can create mapped and guest collections
- Supports all major storage connector types
- Collection access control working
- 80%+ test coverage maintained

---

### Phase 4: Advanced Features (Weeks 19-24)

**Goal**: Add advanced features introduced in GCS v5.4+

#### Deliverables:
- [ ] OIDC server management (v5.4.14+)
  - [ ] `oidc enable` - Enable OIDC server
  - [ ] `oidc show` - Show OIDC configuration
  - [ ] `oidc disable` - Disable OIDC server
- [ ] Sharing policies (v5.4.17+)
  - [ ] `sharing-policy create` - Create policy
  - [ ] `sharing-policy list` - List policies
  - [ ] `sharing-policy show` - Show policy
  - [ ] `sharing-policy update` - Update policy
  - [ ] `sharing-policy delete` - Delete policy
- [ ] Auth policies (v5.4.57+)
  - [ ] `auth-policy create` - Create policy
  - [ ] `auth-policy list` - List policies
  - [ ] `auth-policy show` - Show policy
  - [ ] `auth-policy update` - Update policy
  - [ ] `auth-policy delete` - Delete policy
- [ ] User credentials (v5.4.51+)
  - [ ] `user-credentials create` - Create credential (S3, OAuth, Activescale)
  - [ ] `user-credentials list` - List credentials
  - [ ] `user-credentials show` - Show credential
  - [ ] `user-credentials update` - Update credential
  - [ ] `user-credentials delete` - Delete credential
- [ ] Audit logging (High Assurance)
  - [ ] `audit search` - Search audit logs
  - [ ] `audit export` - Export audit logs
- [ ] Diagnostics
  - [ ] `self-diagnostic run` - Run diagnostic tests
  - [ ] `self-diagnostic show` - Show results

**Milestone**: [Phase 4] Advanced Features Complete

#### Issues to Create:
30. #30 - [Phase 4][System Admin] Implement OIDC server management
31. #31 - [Phase 4][Data Manager] Implement sharing policies
32. #32 - [Phase 4][System Admin] Implement auth policies
33. #33 - [Phase 4][Data Manager] Implement user credentials management
34. #34 - [Phase 4][IT Manager] Implement audit log search and export
35. #35 - [Phase 4][System Admin] Implement self-diagnostic utilities

**Success Criteria:**
- All advanced features functional
- Feature parity with Python GCS v5.4.60+
- Documentation complete
- 80%+ test coverage maintained

---

### Phase 5: Polish, Documentation & Release (Weeks 25-28)

**Goal**: Production-ready release with complete documentation

#### Deliverables:
- [ ] Complete documentation
  - [ ] Command reference for every command
  - [ ] User guide with examples
  - [ ] Migration guide from Python CLI
  - [ ] Troubleshooting guide
  - [ ] API documentation (GoDoc)
- [ ] Persona walkthroughs
  - [ ] System Administrator walkthrough
  - [ ] Data Manager walkthrough
  - [ ] Research PI walkthrough
  - [ ] IT Manager walkthrough
- [ ] Release engineering
  - [ ] Multi-platform binaries (Linux x86_64, ARM64, macOS)
  - [ ] .deb and .rpm packages for Linux
  - [ ] Homebrew formula
  - [ ] Docker image
  - [ ] Installation script
- [ ] Performance optimization
  - [ ] Benchmark against Python CLI
  - [ ] Optimize hot paths
  - [ ] Memory profiling
- [ ] Security audit
  - [ ] Token storage security review
  - [ ] Dependency vulnerability scan
  - [ ] Code security review
- [ ] Compatibility testing
  - [ ] Test against all GCS v5.4.x releases
  - [ ] Validate config file compatibility
  - [ ] Test migrations from Python CLI

**Milestone**: [Phase 5] v1.0.0 Release

#### Issues to Create:
36. #36 - [Phase 5][All Personas] Write comprehensive command reference
37. #37 - [Phase 5][All Personas] Create persona walkthroughs
38. #38 - [Phase 5][System Admin] Set up multi-platform release pipeline
39. #39 - [Phase 5][System Admin] Create installation packages (.deb, .rpm, Homebrew)
40. #40 - [Phase 5][IT Manager] Perform security audit
41. #41 - [Phase 5][System Admin] Compatibility testing across GCS versions
42. #42 - [Phase 5][All Personas] Write migration guide from Python CLI
43. #43 - [Phase 5][IT Manager] Performance benchmarking and optimization

**Success Criteria:**
- All commands documented with examples
- Installation on all major platforms working
- Performance equal or better than Python CLI
- Security audit passed
- v1.0.0 release published

---

## GitHub Project Management Setup

### GitHub Project Structure

**Project Name**: Globus Connect Server Go Port
**Project Board**: https://github.com/users/[username]/projects/[N]

### Custom Fields

Following the lens pattern, configure these custom fields programmatically:

1. **Persona (Single Select)**
   - System Administrator (Blue)
   - Data Manager (Green)
   - Research PI (Yellow)
   - IT Manager (Red)

2. **Phase (Single Select)**
   - Phase 0 - Project Setup (Gray)
   - Phase 1 - Proof of Concept (Blue)
   - Phase 2 - Core Management (Green)
   - Phase 3 - Collections & Storage (Yellow)
   - Phase 4 - Advanced Features (Orange)
   - Phase 5 - Polish & Release (Purple)

3. **Estimate (Number)**
   - Story points or days

4. **Priority (Single Select)**
   - Critical (Red)
   - High (Orange)
   - Medium (Yellow)
   - Low (Green)

### Project Board Views

Create these views manually through GitHub UI:

1. **Kanban (Default)**
   - Layout: Board
   - Group by: Status
   - Columns: Triage, Todo, In Progress, In Review, Done

2. **By Phase (Roadmap)**
   - Layout: Board
   - Group by: Phase
   - Sort by: Priority

3. **By Persona**
   - Layout: Board
   - Group by: Persona
   - Sort by: Phase, then Priority

4. **Current Sprint**
   - Layout: Table
   - Filter: Status ≠ Done, Phase = [current phase]
   - Columns: Title, Status, Priority, Persona, Estimate, Assignees

5. **Backlog**
   - Layout: Table
   - Filter: Status = Todo
   - Sort by: Priority, then Phase
   - Columns: Title, Phase, Priority, Persona, Estimate

### Milestones

Create these milestones to track major deliverables:

- **[Phase 0] Project Infrastructure Complete** (Due: Week 2)
- **[Phase 1] Proof of Concept - Authentication Working** (Due: Week 6)
- **[Phase 2] Endpoint & Node Management Complete** (Due: Week 12)
- **[Phase 3] Collection & Storage Management Complete** (Due: Week 18)
- **[Phase 4] Advanced Features Complete** (Due: Week 24)
- **[Phase 5] v1.0.0 Release** (Due: Week 28)

### Label Configuration

Create 147 labels following lens pattern (see .github/labels.yml):

**Type Labels:**
- bug, enhancement, documentation, technical-debt, question

**Priority Labels:**
- priority: critical, priority: high, priority: medium, priority: low

**Area Labels:**
- area: cli, area: config, area: auth, area: endpoint, area: node, area: collection, area: storage-gateway, area: oidc, area: audit, area: diagnostics, area: tests, area: docs, area: build

**Persona Labels:**
- persona: system-admin, persona: data-manager, persona: research-pi, persona: it-manager

**Phase Labels:**
- phase: 0-setup, phase: 1-mvp, phase: 2-core, phase: 3-collections, phase: 4-advanced, phase: 5-release, phase: backlog

**Status Labels:**
- triage, needs-info, blocked, ready, in-progress, in-review, awaiting-merge

**Special Labels:**
- good first issue, help wanted, breaking-change, security, performance, dependencies

### Issue Creation Process

For each issue:

1. Use appropriate issue template (feature_request, bug_report, technical_debt, documentation)
2. Select persona(s) that benefit from this feature
3. Select phase alignment
4. Add all relevant labels (priority, area, persona, phase)
5. Add to project board (will be automatic with GitHub Actions)
6. Set milestone
7. Set estimate in custom field

### GITHUB_ISSUES_SUMMARY.md

Maintain a summary document tracking all issues:

```markdown
# GitHub Issues Summary

| # | Title | Persona | Phase | Priority | Area | Estimate | Status | Milestone |
|---|-------|---------|-------|----------|------|----------|--------|-----------|
| 1 | Initialize Go module | All | Phase 0 | High | build | 1 day | Done | Phase 0 Complete |
| 2 | Set up CI/CD | System Admin | Phase 0 | High | build | 2 days | Done | Phase 0 Complete |
...

## By Phase

### Phase 0 - Project Setup (6 issues)
- #1, #2, #3, #4, #5, #6

### Phase 1 - Proof of Concept (8 issues)
- #7, #8, #9, #10, #11, #12, #13, #14

[...]

## By Persona

### System Administrator (25 issues)
- #1, #2, #7, #8, #9, #11, #15, #16, #17, #18, #19, #20, #22, ...

[...]
```

---

## Success Metrics

### Project-Level Metrics

- **Feature Parity**: 100% of Python CLI commands implemented
- **Config Compatibility**: 100% config file format compatibility
- **Test Coverage**: ≥80% code coverage across all packages
- **Performance**: ≥ Python CLI performance (faster startup, equal API call time)
- **Documentation**: All commands documented with examples
- **Platforms**: Binaries for Linux (x86_64, ARM64), macOS (x86_64, ARM64)

### Persona-Specific Success Metrics

#### System Administrator
- Single binary installation (no Python dependencies)
- Installation time < 5 minutes
- Works on CentOS, Ubuntu, Rocky Linux, macOS
- Config files 100% compatible with Python CLI

#### Data Manager
- Interactive prompts for complex operations
- Clear, actionable error messages
- Tab completion working
- Example commands in --help output

#### Research PI
- Built-in help accessible (--help, examples)
- Audit logging enabled by default
- Role-based access control working
- Compliance-friendly

#### IT Manager
- JSON output for all commands (--format=json)
- Scriptable (no interactive prompts with --non-interactive flag)
- Consistent exit codes (0 = success, non-zero = failure)
- Automation-friendly

---

## Technical Considerations

### Advantages of Go Port

1. **Single Binary Deployment**
   - No Python runtime required
   - No dependency conflicts
   - Smaller deployment footprint

2. **Performance**
   - Faster startup time (no Python interpreter initialization)
   - Lower memory usage
   - Compiled, not interpreted

3. **Cross-Platform**
   - Easy cross-compilation (GOOS, GOARCH)
   - ARM64 support out-of-box
   - Static linking possible

4. **Maintainability**
   - Strong typing catches errors at compile time
   - Better IDE support (gopls)
   - Clear dependency management (go.mod)
   - Built-in testing framework

### Challenges & Mitigations

| Challenge | Mitigation |
|-----------|------------|
| **Config Format Compatibility** | Use same JSON format, validate against Python CLI |
| **API Compatibility** | Extend globus-go-sdk with GCS Manager API client |
| **Feature Discovery** | Reverse-engineer Python CLI, read GCS docs thoroughly |
| **Testing Without Real Endpoint** | Mock API server, staging environment access |
| **Error Message Parity** | Copy error messages from Python CLI where applicable |

### Dependencies

- **globus-go-sdk**: Provides base Globus API clients (auth, transfer)
  - Extend with GCS Manager API client
  - Reuse OAuth2 token management
- **Cobra**: CLI framework (same as globus-go-cli)
- **Viper**: Configuration management
- **golang.org/x/oauth2**: OAuth2 client
- **testify**: Testing assertions and mocks

---

## Development Workflow

### Sprint Planning

- **Sprint Length**: 2 weeks
- **Sprint Planning**: Beginning of each sprint
  - Review current phase backlog
  - Select issues for sprint (based on priority, estimate)
  - Assign issues to developers
  - Update project board

### Daily Standup (Async)

- Update issue status on project board
- Comment on blocked issues
- Review PR status

### Code Review Process

1. Create feature branch from main
2. Implement feature with tests
3. Open PR using template
   - Fill in persona impact section
   - Link related issues (Closes #X)
   - Link related documentation
4. CI runs automatically (tests, lints, builds)
5. Request review from maintainer
6. Address review comments
7. Merge to main after approval
8. Close related issues
9. Update milestone progress

### Release Process

Following semantic versioning (v1.0.0):

1. Update CHANGELOG.md with release notes
2. Create release branch (release/v1.0.0)
3. Update version in code
4. Tag release (git tag v1.0.0)
5. CI builds binaries for all platforms
6. Create GitHub release with binaries
7. Update Homebrew formula
8. Announce release

---

## Documentation Structure

### For Users

- **README.md**: Overview, installation, quick start
- **docs/INSTALLATION.md**: Detailed installation instructions
- **docs/MIGRATION_FROM_PYTHON.md**: Migration guide for Python CLI users
- **docs/COMMAND_REFERENCE.md**: Complete command documentation
- **docs/USER_GUIDE.md**: Task-oriented guide with examples
- **docs/TROUBLESHOOTING.md**: Common issues and solutions
- **docs/USER_SCENARIOS/**: Persona walkthroughs

### For Contributors

- **CONTRIBUTING.md**: Contribution guidelines
- **docs/DEVELOPMENT.md**: Development setup, coding standards
- **docs/ARCHITECTURE.md**: System architecture and design decisions
- **docs/TESTING.md**: Testing strategy and guidelines

### For Project Management

- **PROJECT_PLAN.md**: This document
- **ROADMAP.md**: High-level roadmap and timeline
- **.github/PROJECT_BOARD_SETUP.md**: Project board configuration
- **.github/GITHUB_ISSUES_SUMMARY.md**: Issue tracking summary

---

## Next Steps

### Immediate Actions (Week 1)

1. **Set up GitHub repository**
   - Create repo: globus-go-gcs
   - Add README.md with project overview
   - Add LICENSE (Apache 2.0 to match globus-go-sdk)
   - Add .gitignore for Go projects

2. **Configure GitHub Project Board**
   - Create new project board
   - Add custom fields (Persona, Phase, Estimate, Priority)
   - Create views (Kanban, By Phase, By Persona, Current Sprint, Backlog)

3. **Create labels**
   - Add .github/labels.yml with 147 labels
   - Run label sync workflow

4. **Create issue templates**
   - Feature request with persona and phase fields
   - Bug report with persona field
   - Technical debt template
   - Documentation template

5. **Create PR template**
   - Add persona impact section
   - Add checklist for tests, docs, changelog

6. **Create initial milestones**
   - [Phase 0] Project Infrastructure Complete
   - [Phase 1] Proof of Concept - Authentication Working
   - [Phase 2] Endpoint & Node Management Complete
   - [Phase 3] Collection & Storage Management Complete
   - [Phase 4] Advanced Features Complete
   - [Phase 5] v1.0.0 Release

7. **Create initial issues**
   - Create issues #1-#6 for Phase 0 (project setup)
   - Add all labels, milestones, estimates
   - Add to project board

8. **Initialize Go project**
   - `go mod init github.com/scttfrdmn/globus-go-gcs`
   - Create directory structure
   - Add Makefile
   - Set up CI/CD (GitHub Actions)

### This Week (Week 1)

- Complete all "Immediate Actions" above
- Begin Phase 0 implementation
- Set up development environment
- Research Python GCS CLI implementation (if open source)

### Next Week (Week 2)

- Complete Phase 0
- Begin Phase 1 (Proof of Concept)
- Start OAuth2 authentication implementation

---

## Questions Resolved ✅

1. **GCS Manager API Access**: ✅ **RESOLVED**
   - Python SDK v4 has GCSClient with full API documentation
   - NOT in globus-go-sdk (v3 or v4) - we'll implement it
   - API docs: https://globus-sdk-python.readthedocs.io/en/stable/services/gcs.html
   - GCS Manager API docs: https://docs.globus.org/globus-connect-server/v5/api/

2. **Python CLI Source Code**: ✅ **RESOLVED**
   - Yes, Python GCS CLI is open source
   - globus-go-cli is based on it
   - Can reference implementation for command structure

3. **Testing Environment**: ✅ **RESOLVED**
   - User has Globus subscription for testing
   - Can deploy GCS locally and on AWS
   - Have access to real GCS endpoints for integration testing

4. **Globus SDK Version**: ✅ **RESOLVED**
   - Use globus-go-sdk **v3.65.0-1+** (current stable)
   - Will implement GCS Manager API client in this project
   - Prepare for v4 migration when stable (Q4 2025)

---

## References

- **Python GCS CLI**: https://docs.globus.org/globus-connect-server/v5/reference/
- **globus-go-sdk**: ../globus-go-sdk/
- **globus-go-cli**: ../globus-go-cli/
- **lens project**: ../lens/ (for project management patterns)
- **GCS v5 Installation Guide**: https://docs.globus.org/globus-connect-server/v5/
- **Globus SDK Python**: https://github.com/globus/globus-sdk-python

---

**Document Owner**: Project Lead
**Last Updated**: 2025-10-25
**Status**: Initial Planning
**Next Review**: Weekly during Phase 0, bi-weekly afterwards
