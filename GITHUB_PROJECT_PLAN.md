# GitHub Project Plan: Security Remediation & Optimization

## Project Overview

This document outlines the GitHub Projects structure and issues for implementing NIST 800-53 security remediation and Go performance optimizations for globus-go-gcs.

**Project Goals**:
1. Achieve HIPAA/PHI compliance for production deployment
2. Leverage Go's capabilities for 2-5x performance improvements
3. Maintain 100% feature parity with Python version

**Timeline**: 6-8 weeks total
- Security Remediation: 3-4 weeks (15-20 developer days)
- Optimizations: 1.5-2 weeks (8-12 developer days)
- Testing & Validation: 1-2 weeks

---

## GitHub Projects Board Structure

### Project: "Security & Performance v2.0"

**Columns**:
1. **Backlog** - All issues awaiting triage
2. **Ready** - Prioritized and ready to start
3. **In Progress** - Actively being worked
4. **In Review** - PR submitted, awaiting review
5. **Testing** - Merged, needs validation
6. **Done** - Completed and verified

**Views**:
- **By Priority**: Group by priority label (P0-Critical → P3-Low)
- **By Epic**: Group by epic (Security/Optimization/Testing)
- **By Timeline**: Roadmap view with milestone dates

---

## Milestones

### Milestone 1: Critical Security Fixes (Week 1-2)
**Due**: 2 weeks from start
**Goal**: Address all CRITICAL severity findings
- Token encryption at rest
- TLS hardening
- Secrets handling in CLI
- Audit database encryption
- CSRF token strengthening

### Milestone 2: High Priority Security (Week 3-4)
**Due**: 4 weeks from start
**Goal**: Address all HIGH severity findings
- Input validation
- Error information disclosure
- Session timeout enforcement
- Rate limiting
- Path traversal protection

### Milestone 3: High Priority Optimizations (Week 5-6)
**Due**: 6 weeks from start
**Goal**: Implement high-impact performance improvements
- HTTP connection pooling
- JSON buffer pooling
- Batch operation concurrency
- Audit log pipeline
- Database query optimization

### Milestone 4: Security & Performance Testing (Week 7-8)
**Due**: 8 weeks from start
**Goal**: Comprehensive validation
- Security audit
- Performance benchmarking
- Load testing
- Documentation updates

---

## Labels

### Priority
- `P0-Critical` - Security vulnerabilities, data loss risk
- `P1-High` - Major security/performance issues
- `P2-Medium` - Moderate improvements
- `P3-Low` - Nice-to-have enhancements

### Type
- `type: security` - Security remediation
- `type: performance` - Optimization work
- `type: testing` - Test coverage
- `type: documentation` - Docs updates

### Area
- `area: auth` - Authentication & tokens
- `area: tls` - TLS/HTTPS configuration
- `area: audit` - Audit logging
- `area: api-client` - GCS API client
- `area: cli` - CLI commands
- `area: database` - SQLite operations

### Effort
- `effort: small` - 1-4 hours
- `effort: medium` - 4-16 hours (1-2 days)
- `effort: large` - 16-40 hours (2-5 days)
- `effort: xl` - 40+ hours (5+ days)

### Epic
- `epic: security-remediation` - NIST 800-53 compliance
- `epic: go-optimizations` - Performance improvements
- `epic: testing` - Test & validation

### Special
- `breaking-change` - Issues that introduce breaking changes (requires migration guide)

---

## Epic 1: Security Remediation (NIST 800-53 Compliance)

### CRITICAL Severity Issues

#### Issue #1: Token Encryption at Rest
**Title**: Implement AES-256-GCM encryption for stored OAuth tokens
**Priority**: P0-Critical
**Labels**: `type: security`, `area: auth`, `effort: large`, `epic: security-remediation`
**Milestone**: Milestone 1
**Estimated Effort**: 12-16 hours

**Description**:
Currently OAuth tokens are stored in `~/.globus-connect-server/tokens/<profile>.json` with only file permissions (0600) for protection. Tokens must be encrypted at rest using AES-256-GCM.

**NIST 800-53 Controls**: SC-28 (Protection of Information at Rest)

**Implementation**:
- [ ] Derive encryption key from system keyring (use `github.com/zalando/go-keyring`)
- [ ] Implement AES-256-GCM encryption wrapper in `internal/auth/encryption.go`
- [ ] Migrate existing token storage to use encryption in `internal/auth/tokens.go`
- [ ] Add key rotation support
- [ ] Handle keyring unavailable fallback (prompt for passphrase)
- [ ] Add migration tool for existing plaintext tokens

**Files to modify**:
- `internal/auth/tokens.go` (LoadToken, SaveToken functions)
- New: `internal/auth/encryption.go`
- New: `internal/auth/keyring.go`

**Testing**:
- Unit tests for encryption/decryption
- Integration test for token save/load cycle
- Test key rotation
- Test keyring unavailable fallback

**Documentation**:
- Update README with keyring requirements
- Document token encryption behavior
- Add troubleshooting guide for keyring issues

---

#### Issue #2: TLS Configuration Hardening
**Title**: Enforce TLS 1.2+ with secure cipher suites
**Priority**: P0-Critical
**Labels**: `type: security`, `area: tls`, `area: api-client`, `effort: medium`, `epic: security-remediation`
**Milestone**: Milestone 1
**Estimated Effort**: 6-8 hours

**Description**:
HTTP client uses default TLS configuration without enforcing minimum TLS version or approved cipher suites. Must implement strict TLS configuration.

**NIST 800-53 Controls**: SC-8 (Transmission Confidentiality), SC-13 (Cryptographic Protection)

**Implementation**:
- [ ] Create custom TLS config with MinVersion: TLS 1.2
- [ ] Define approved cipher suite list (ECDHE-RSA-AES256-GCM, etc.)
- [ ] Add certificate validation options
- [ ] Implement certificate pinning support (optional)
- [ ] Add TLS debugging/logging capability
- [ ] Add configuration options for TLS settings

**Files to modify**:
- `pkg/gcs/client.go` (Client creation, HTTP transport)
- `pkg/gcs/options.go` (Add TLS configuration options)
- New: `pkg/gcs/tls.go`

**Configuration** (add to config.yaml):
```yaml
tls:
  min_version: "1.2"
  cipher_suites:
    - TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
    - TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
  verify_certificate: true
  certificate_pinning: false
```

**Testing**:
- Unit tests for TLS config validation
- Integration test against mock HTTPS server
- Test TLS 1.0/1.1 rejection
- Test weak cipher suite rejection

**Documentation**:
- Document TLS requirements
- Add TLS configuration guide
- Document certificate pinning setup

---

#### Issue #3: Remove Secrets from CLI Arguments
**Title**: 🚨 BREAKING CHANGE: Refactor secret input to use stdin/prompts instead of CLI args
**Priority**: P0-Critical
**Labels**: `type: security`, `area: cli`, `effort: large`, `epic: security-remediation`, `breaking-change`
**Milestone**: Milestone 1
**Estimated Effort**: 12-16 hours
**⚠️ UX Impact**: **YES - BREAKING CHANGE** - Changes user interaction pattern

**📖 Full Migration Guide**: [BREAKING_CHANGES_V2.md](./BREAKING_CHANGES_V2.md)

**Description**:
Multiple commands accept secrets via CLI arguments (e.g., `--secret-access-key`, `--client-secret`), which are visible in process listings and shell history. Must refactor to use secure input methods.

**NIST 800-53 Controls**: IA-5(7) (No Embedded Passwords)

**Affected Commands**:
- `user-credential s3-keys add` (--secret-access-key)
- `user-credential s3-keys update` (--secret-access-key)
- `user-credential activescale-create` (--password)
- `oidc create` (--client-secret)
- `oidc update` (--client-secret)

**Implementation**:
- [ ] Create secure input utility in `internal/input/secure.go`
- [ ] Implement interactive prompt with hidden input (using `golang.org/x/term`)
- [ ] Add `--secret-stdin` flag for pipe input
- [ ] Support environment variable fallback (GLOBUS_SECRET_VALUE)
- [ ] Remove direct secret flags, show deprecation warning
- [ ] Update all affected commands
- [ ] Add validation for secret format/length

**New CLI Patterns**:
```bash
# Interactive prompt (recommended)
$ globus-connect-server user-credential s3-keys add --access-key-id AKIA...
Enter secret access key: ********
Confirm secret access key: ********

# Stdin pipe
$ echo "secret-value" | globus-connect-server oidc create --client-secret-stdin

# Environment variable
$ export GLOBUS_CLIENT_SECRET="secret-value"
$ globus-connect-server oidc create --client-secret-env
```

**Files to modify**:
- New: `internal/input/secure.go`
- `internal/commands/usercredential/s3_keys_add.go`
- `internal/commands/usercredential/s3_keys_update.go`
- `internal/commands/usercredential/activescale_create.go`
- `internal/commands/oidc/create.go`
- `internal/commands/oidc/update.go`

**Testing**:
- Unit tests for secure input functions
- Integration tests for all three input methods
- Test empty/invalid secret handling
- Test confirmation mismatch handling

**Documentation**:
- **Migration guide for users** (breaking change)
- Update all command examples
- Add security best practices guide

---

#### Issue #4: Audit Database Encryption
**Title**: Implement SQLCipher for encrypted audit database
**Priority**: P0-Critical
**Labels**: `type: security`, `area: audit`, `area: database`, `effort: large`, `epic: security-remediation`
**Milestone**: Milestone 1
**Estimated Effort**: 12-16 hours

**Description**:
Audit logs stored in SQLite database at `~/.globus-connect-server/audit/audit.db` are unencrypted. Must use SQLCipher for encryption at rest.

**NIST 800-53 Controls**: AU-9 (Protection of Audit Information), SC-28 (Protection of Information at Rest)

**Implementation**:
- [ ] Replace `modernc.org/sqlite` with SQLCipher-enabled driver
- [ ] Use `github.com/mutecomm/go-sqlcipher/v4` or build SQLCipher support
- [ ] Derive database encryption key from system keyring
- [ ] Add key management functions
- [ ] Migrate existing databases to encrypted format
- [ ] Add database backup encryption
- [ ] Implement secure database deletion (shred)

**Files to modify**:
- `internal/commands/audit/audit.go` (database initialization)
- All audit command files (connection strings)
- New: `internal/commands/audit/migration.go` (plaintext → encrypted)

**Configuration**:
```yaml
audit:
  database_path: ~/.globus-connect-server/audit/audit.db
  encrypted: true
  cipher: aes256
  kdf_iterations: 256000
```

**Testing**:
- Unit tests for encrypted database operations
- Migration test (plaintext → encrypted)
- Performance test (ensure acceptable overhead)
- Test key rotation
- Test database backup/restore

**Documentation**:
- Document SQLCipher requirement
- Add database encryption guide
- Document migration procedure

**Dependencies**:
- Requires Issue #1 (keyring integration)
- May require CGO for SQLCipher (check cross-compilation impact)

---

#### Issue #5: Strengthen CSRF Token Generation
**Title**: Replace timestamp-based CSRF with crypto/rand
**Priority**: P0-Critical
**Labels**: `type: security`, `area: auth`, `effort: small`, `epic: security-remediation`
**Milestone**: Milestone 1
**Estimated Effort**: 2-4 hours

**Description**:
CSRF state token in OAuth flow uses timestamp-based generation (`state-%d`), which is predictable. Must use cryptographically secure random generation.

**NIST 800-53 Controls**: IA-9 (Service Identification and Authentication), SC-13 (Cryptographic Protection)

**Current Code** (`internal/commands/auth/login.go:86-88`):
```go
// Generate state parameter for CSRF protection
state := fmt.Sprintf("state-%d", time.Now().Unix())
```

**Implementation**:
- [ ] Replace with crypto/rand.Read (32 bytes)
- [ ] Base64 URL-safe encode the random bytes
- [ ] Store state with timestamp for expiration (5 minutes)
- [ ] Add state validation in callback handler
- [ ] Clean up expired states

**New Code**:
```go
// generateSecureState generates a cryptographically secure random state token
func generateSecureState() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", fmt.Errorf("generate random state: %w", err)
    }
    return base64.URLEncoding.EncodeToString(b), nil
}
```

**Files to modify**:
- `internal/commands/auth/login.go`

**Testing**:
- Unit test for randomness (no duplicates in 10000 generations)
- Test state validation
- Test state expiration

**Documentation**:
- Update OAuth flow documentation

---

### HIGH Severity Issues

#### Issue #6: Comprehensive Input Validation
**Title**: Add validation for all user inputs (paths, IDs, names)
**Priority**: P1-High
**Labels**: `type: security`, `area: cli`, `effort: large`, `epic: security-remediation`
**Milestone**: Milestone 2
**Estimated Effort**: 16-20 hours

**Description**:
Many commands lack input validation for user-provided values, potentially allowing injection attacks or malformed data.

**NIST 800-53 Controls**: SI-10 (Information Input Validation)

**Validation Needed**:
- Collection IDs (UUID format)
- Storage Gateway IDs (UUID format)
- Endpoint FQDNs (valid domain name)
- File paths (no directory traversal)
- Identity URNs (valid URN format)
- Email addresses (RFC 5322)
- Display names (length, character restrictions)
- JSON payloads (schema validation)

**Implementation**:
- [ ] Create validation package `internal/validation/`
- [ ] Implement validators for each data type
- [ ] Add `github.com/go-playground/validator/v10` for struct validation
- [ ] Update all command RunE functions to validate inputs
- [ ] Add clear error messages for validation failures
- [ ] Add regex patterns for each format

**New Package Structure**:
```
internal/validation/
├── validators.go      # Core validation functions
├── formats.go        # Format constants and regex patterns
├── errors.go         # Validation error types
└── validators_test.go
```

**Files to modify**:
- All 83 command files (add validation calls)
- `pkg/gcs/types.go` (add validation struct tags)

**Example**:
```go
// Before
func runCreate(ctx context.Context, name, gatewayID string, ...) error {
    collection := &gcs.Collection{
        DisplayName: name,
        StorageGatewayID: gatewayID,
    }
    // ...
}

// After
func runCreate(ctx context.Context, name, gatewayID string, ...) error {
    if err := validation.ValidateDisplayName(name); err != nil {
        return fmt.Errorf("invalid collection name: %w", err)
    }
    if err := validation.ValidateUUID(gatewayID); err != nil {
        return fmt.Errorf("invalid storage gateway ID: %w", err)
    }
    collection := &gcs.Collection{
        DisplayName: name,
        StorageGatewayID: gatewayID,
    }
    // ...
}
```

**Testing**:
- Unit tests for each validator (valid/invalid cases)
- Integration tests for command validation
- Fuzzing tests for each validator

**Documentation**:
- Document valid formats for each input type
- Add input validation guide

---

#### Issue #7: Reduce Error Information Disclosure
**Title**: Sanitize error messages to prevent information leakage
**Priority**: P1-High
**Labels**: `type: security`, `area: api-client`, `area: cli`, `effort: medium`, `epic: security-remediation`
**Milestone**: Milestone 2
**Estimated Effort**: 8-12 hours

**Description**:
Error messages may expose internal paths, tokens, or system information. Need structured error handling with user-safe messages.

**NIST 800-53 Controls**: SI-11 (Error Handling)

**Current Issues**:
- Token load errors expose full file paths
- API errors may include sensitive headers
- Database errors expose SQL queries
- Stack traces in some error paths

**Implementation**:
- [ ] Create error sanitization package `internal/errors/`
- [ ] Define public vs. internal error messages
- [ ] Add debug mode for detailed errors (--debug flag)
- [ ] Wrap all errors with sanitized messages
- [ ] Log detailed errors to debug log file
- [ ] Audit all error returns

**Error Handling Pattern**:
```go
// Before
if err != nil {
    return fmt.Errorf("failed to load token from %s: %w", tokenPath, err)
}

// After
if err != nil {
    log.Debugf("Failed to load token from %s: %v", tokenPath, err)
    return errors.NewUserError("authentication failed: token not found or invalid")
}
```

**Files to modify**:
- New: `internal/errors/errors.go`
- All command files (error handling)
- `pkg/gcs/client.go` (API error handling)

**Testing**:
- Test error messages don't contain sensitive info
- Test debug mode exposes full details
- Test error logging

**Documentation**:
- Document error handling guidelines
- Add troubleshooting guide with debug mode

---

#### Issue #8: Implement Session Timeout Enforcement
**Title**: Add configurable session timeouts with automatic refresh
**Priority**: P1-High
**Labels**: `type: security`, `area: auth`, `effort: medium`, `epic: security-remediation`
**Milestone**: Milestone 2
**Estimated Effort**: 8-12 hours

**Description**:
No session timeout enforcement - tokens remain valid until expiration. Should implement configurable timeouts with automatic refresh.

**NIST 800-53 Controls**: AC-12 (Session Termination)

**Implementation**:
- [ ] Add last_used timestamp to token storage
- [ ] Check last_used on every command execution
- [ ] Implement configurable inactivity timeout (default: 30 minutes)
- [ ] Add automatic token refresh when near expiration
- [ ] Clear token on timeout
- [ ] Add session timeout warnings

**Configuration**:
```yaml
auth:
  session_timeout_mins: 60          # Absolute timeout
  inactivity_timeout_mins: 30       # Inactivity timeout
  auto_refresh: true                # Auto-refresh tokens
  refresh_threshold_mins: 5         # Refresh when <5 mins remaining
```

**Files to modify**:
- `internal/auth/tokens.go` (add timeout checking)
- `pkg/config/config.go` (add timeout config)
- All commands (check timeout before execution)

**Testing**:
- Test inactivity timeout enforcement
- Test automatic refresh
- Test absolute timeout
- Test session timeout warnings

**Documentation**:
- Document session timeout behavior
- Add configuration guide

---

#### Issue #9: API Rate Limiting and Retry Logic
**Title**: Implement exponential backoff and rate limit handling
**Priority**: P1-High
**Labels**: `type: security`, `type: performance`, `area: api-client`, `effort: medium`, `epic: security-remediation`
**Milestone**: Milestone 2
**Estimated Effort**: 6-8 hours

**Description**:
API client doesn't handle rate limits or implement retry logic, potentially causing service disruption or unintentional DoS.

**NIST 800-53 Controls**: SC-5 (Denial of Service Protection)

**Implementation**:
- [ ] Detect 429 (Too Many Requests) responses
- [ ] Implement exponential backoff with jitter
- [ ] Add Retry-After header parsing
- [ ] Add configurable max retries (default: 3)
- [ ] Add rate limit headroom monitoring
- [ ] Add request throttling for batch operations

**Configuration**:
```yaml
api:
  max_retries: 3
  initial_retry_delay: 1s
  max_retry_delay: 30s
  retry_jitter: 0.1
  requests_per_second: 10          # Client-side throttle
```

**Files to modify**:
- `pkg/gcs/client.go` (add retry logic to doRequest)
- New: `pkg/gcs/retry.go`
- `pkg/gcs/options.go` (add rate limit config)

**Implementation Example**:
```go
func (c *Client) doRequestWithRetry(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
    var resp *http.Response
    var err error

    backoff := c.initialRetryDelay
    for attempt := 0; attempt <= c.maxRetries; attempt++ {
        resp, err = c.doRequest(ctx, method, path, body)

        if err == nil && resp.StatusCode != http.StatusTooManyRequests {
            return resp, nil
        }

        if attempt < c.maxRetries {
            // Parse Retry-After header
            retryAfter := parseRetryAfter(resp)
            if retryAfter > 0 {
                time.Sleep(retryAfter)
            } else {
                // Exponential backoff with jitter
                jitter := time.Duration(rand.Float64() * float64(backoff) * c.retryJitter)
                time.Sleep(backoff + jitter)
                backoff *= 2
                if backoff > c.maxRetryDelay {
                    backoff = c.maxRetryDelay
                }
            }
        }
    }

    return nil, fmt.Errorf("max retries exceeded: %w", err)
}
```

**Testing**:
- Test 429 handling and backoff
- Test Retry-After header parsing
- Test max retry limit
- Test jitter distribution

**Documentation**:
- Document retry behavior
- Document rate limit configuration

---

#### Issue #10: Path Traversal Protection
**Title**: Validate and sanitize all file path operations
**Priority**: P1-High
**Labels**: `type: security`, `area: cli`, `effort: small`, `epic: security-remediation`
**Milestone**: Milestone 2
**Estimated Effort**: 4-6 hours

**Description**:
File operations in audit commands don't validate paths, potentially allowing directory traversal attacks.

**NIST 800-53 Controls**: SI-10 (Information Input Validation)

**Affected Commands**:
- `audit dump --output` (arbitrary file write)
- Config file loading (arbitrary file read)

**Implementation**:
- [ ] Create path validation utility
- [ ] Reject paths with `..` components
- [ ] Validate paths are within expected directories
- [ ] Add allowlist for valid directories
- [ ] Require absolute paths or resolve to absolute
- [ ] Add path permission checks

**Files to modify**:
- New: `internal/validation/paths.go`
- `internal/commands/audit/dump.go`
- `pkg/config/config.go`

**Example**:
```go
func ValidateOutputPath(path string) error {
    // Resolve to absolute path
    absPath, err := filepath.Abs(path)
    if err != nil {
        return fmt.Errorf("invalid path: %w", err)
    }

    // Check for directory traversal
    if strings.Contains(path, "..") {
        return fmt.Errorf("path traversal not allowed")
    }

    // Ensure writable directory
    dir := filepath.Dir(absPath)
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        return fmt.Errorf("directory does not exist: %s", dir)
    }

    return nil
}
```

**Testing**:
- Test path traversal rejection (../../etc/passwd)
- Test valid path acceptance
- Test permission handling

**Documentation**:
- Document path validation rules

---

### Additional HIGH Severity Issues (Issues #11-#13)

*(Abbreviated for length - follow same format)*

- **Issue #11**: Implement Secure Configuration File Permissions
- **Issue #12**: Add Token Rotation Support
- **Issue #13**: Implement Audit Log Integrity Verification

---

## Epic 2: Go Performance Optimizations

### HIGH Priority Optimizations

#### Issue #14: HTTP Connection Pooling
**Title**: Implement HTTP client connection pooling with tuned parameters
**Priority**: P1-High
**Labels**: `type: performance`, `area: api-client`, `effort: small`, `epic: go-optimizations`
**Milestone**: Milestone 3
**Estimated Effort**: 3-4 hours
**Expected Impact**: 30-50% faster API calls

**Description**:
Currently each GCS client creates a new HTTP client with default settings. Implement shared connection pooling with optimized parameters.

**Implementation**:
- [ ] Create shared HTTP transport with connection pooling
- [ ] Configure MaxIdleConns: 100
- [ ] Configure MaxIdleConnsPerHost: 10
- [ ] Configure IdleConnTimeout: 90s
- [ ] Enable HTTP/2
- [ ] Add connection metrics

**Files to modify**:
- `pkg/gcs/client.go`
- `pkg/gcs/options.go`
- New: `pkg/gcs/transport.go`

**Code Example**:
```go
var defaultTransport = &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    IdleConnTimeout:     90 * time.Second,
    TLSHandshakeTimeout: 10 * time.Second,
    ForceAttemptHTTP2:   true,
}

func NewClient(endpoint string, opts ...ClientOption) (*Client, error) {
    c := &Client{
        endpoint:   endpoint,
        httpClient: &http.Client{Transport: defaultTransport},
    }
    // ...
}
```

**Testing**:
- Benchmark connection reuse
- Test concurrent requests
- Measure connection pool metrics

**Performance Target**: 30-50% reduction in API call latency

---

#### Issue #15: JSON Buffer Pooling with sync.Pool
**Title**: Implement sync.Pool for JSON marshaling buffers
**Priority**: P1-High
**Labels**: `type: performance`, `area: api-client`, `effort: small`, `epic: go-optimizations`
**Milestone**: Milestone 3
**Estimated Effort**: 3-4 hours
**Expected Impact**: 40-60% fewer allocations

**Description**:
70+ calls to `json.Marshal` across codebase. Implement buffer pooling to reduce allocations.

**Implementation**:
- [ ] Create buffer pool with sync.Pool
- [ ] Create helper function for pooled JSON encoding
- [ ] Replace json.Marshal calls in hot paths
- [ ] Add buffer pool metrics

**Files to modify**:
- `pkg/gcs/client.go` (doRequest method)
- New: `pkg/gcs/pool.go`

**Code Example**:
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func marshalJSON(v interface{}) ([]byte, error) {
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufferPool.Put(buf)

    if err := json.NewEncoder(buf).Encode(v); err != nil {
        return nil, err
    }

    // Copy to new slice since buf will be reused
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    return result, nil
}
```

**Testing**:
- Benchmark memory allocations (before/after)
- Test pool exhaustion under load
- Verify correctness

**Performance Target**: 40-60% reduction in allocations

---

#### Issue #16: Concurrent Batch Operations
**Title**: Implement goroutine worker pool for batch-delete operations
**Priority**: P1-High
**Labels**: `type: performance`, `area: cli`, `effort: medium`, `epic: go-optimizations`
**Milestone**: Milestone 3
**Estimated Effort**: 6-8 hours
**Expected Impact**: 3-5x faster batch operations

**Description**:
Batch delete operations are sequential. Implement concurrent deletion with goroutine worker pool.

**Implementation**:
- [ ] Create worker pool pattern with configurable workers
- [ ] Implement job queue with channels
- [ ] Add progress reporting with atomic counters
- [ ] Add error aggregation
- [ ] Add rate limiting to respect API limits
- [ ] Add graceful cancellation with context

**Files to modify**:
- `internal/commands/collection/batch_delete.go`
- `internal/commands/role/batch_delete.go`
- New: `internal/concurrent/worker_pool.go`

**Code Example**:
```go
type WorkerPool struct {
    workers   int
    jobs      chan Job
    results   chan Result
    errors    chan error
    wg        sync.WaitGroup
}

func (p *WorkerPool) Process(ctx context.Context, items []string) error {
    // Start workers
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(ctx)
    }

    // Send jobs
    go func() {
        for _, item := range items {
            select {
            case p.jobs <- Job{ID: item}:
            case <-ctx.Done():
                return
            }
        }
        close(p.jobs)
    }()

    // Wait for completion
    p.wg.Wait()
    close(p.results)
    close(p.errors)

    return nil
}
```

**Configuration**:
```yaml
performance:
  batch_workers: 10              # Concurrent goroutines
  batch_buffer_size: 100         # Channel buffer size
```

**Testing**:
- Benchmark sequential vs concurrent (100 items)
- Test error handling
- Test cancellation
- Test rate limiting integration

**Performance Target**: 3-5x faster for 100+ item batches

---

#### Issue #17: Audit Log Pipeline Architecture
**Title**: Implement channel-based pipeline for audit log processing
**Priority**: P1-High
**Labels**: `type: performance`, `area: audit`, `area: database`, `effort: large`, `epic: go-optimizations`
**Milestone**: Milestone 3
**Estimated Effort**: 12-16 hours
**Expected Impact**: 4-6x faster audit load

**Description**:
Audit log loading is sequential (fetch → parse → insert). Implement pipeline pattern with concurrent stages.

**Implementation**:
- [ ] Stage 1: Fetch from API (producer)
- [ ] Stage 2: Parse and transform (workers)
- [ ] Stage 3: Batch insert to DB (consumer)
- [ ] Use buffered channels between stages
- [ ] Implement batch insert (100 records per transaction)
- [ ] Add backpressure handling
- [ ] Add progress reporting

**Files to modify**:
- `internal/commands/audit/load.go`
- New: `internal/commands/audit/pipeline.go`

**Architecture**:
```
API Fetch ──[chan]──> Parse Workers ──[chan]──> Batch Inserter
  (1)                    (N=10)                      (1)

Channel Buffer: 1000 records
Batch Size: 100 records per transaction
```

**Code Example**:
```go
func runPipeline(ctx context.Context, client *gcs.Client, db *sql.DB) error {
    // Create channels
    fetchChan := make(chan *gcs.AuditLog, 1000)
    insertChan := make(chan *gcs.AuditLog, 1000)

    // Stage 1: Fetch
    go fetchLogs(ctx, client, fetchChan)

    // Stage 2: Parse (10 workers)
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            parseLogs(ctx, fetchChan, insertChan)
        }()
    }

    // Close insertChan when all parsers done
    go func() {
        wg.Wait()
        close(insertChan)
    }()

    // Stage 3: Batch insert
    return batchInsertLogs(ctx, db, insertChan)
}
```

**Testing**:
- Benchmark pipeline vs sequential (10k records)
- Test stage backpressure
- Test error propagation
- Test graceful shutdown

**Performance Target**: 4-6x faster for large datasets (10k+ records)

---

#### Issue #18: Database Query Optimization
**Title**: Add indexes and optimize SQLite queries for audit operations
**Priority**: P1-High
**Labels**: `type: performance`, `area: database`, `effort: small`, `epic: go-optimizations`
**Milestone**: Milestone 3
**Estimated Effort**: 3-4 hours
**Expected Impact**: 25-40% faster queries

**Description**:
Audit queries could be optimized with better indexes and prepared statements.

**Implementation**:
- [ ] Add composite indexes for common query patterns
- [ ] Use prepared statements for repeated queries
- [ ] Enable SQLite performance pragmas
- [ ] Add query result caching for frequent queries
- [ ] Optimize LIKE queries with FTS5

**Files to modify**:
- `internal/commands/audit/audit.go` (schema)
- `internal/commands/audit/query.go`

**New Schema Indexes**:
```sql
-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_event_time ON audit_logs(event_type, timestamp);
CREATE INDEX IF NOT EXISTS idx_identity_time ON audit_logs(identity_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_result_time ON audit_logs(result, timestamp);

-- Full-text search for message field
CREATE VIRTUAL TABLE IF NOT EXISTS audit_logs_fts USING fts5(message, content='audit_logs');
```

**SQLite Pragmas**:
```go
db.Exec("PRAGMA journal_mode = WAL")           // Write-ahead logging
db.Exec("PRAGMA synchronous = NORMAL")         // Faster writes
db.Exec("PRAGMA cache_size = -64000")          // 64MB cache
db.Exec("PRAGMA temp_store = MEMORY")          // In-memory temp tables
db.Exec("PRAGMA mmap_size = 268435456")        // 256MB memory-mapped I/O
```

**Prepared Statements**:
```go
type AuditDB struct {
    db               *sql.DB
    insertStmt       *sql.Stmt
    queryByEventStmt *sql.Stmt
    queryByTimeStmt  *sql.Stmt
}

func (db *AuditDB) PrepareStatements() error {
    var err error
    db.insertStmt, err = db.db.Prepare(`
        INSERT INTO audit_logs (...) VALUES (?, ?, ?, ...)
    `)
    // ...
}
```

**Testing**:
- Benchmark query performance with indexes
- Test query plan with EXPLAIN QUERY PLAN
- Benchmark prepared statements vs dynamic queries

**Performance Target**: 25-40% faster query execution

---

### MEDIUM Priority Optimizations (Issues #19-#23)

*(Abbreviated for length)*

- **Issue #19**: Implement Context Timeout Tuning (4-6 hours)
- **Issue #20**: Add Response Streaming for Large Datasets (8-12 hours)
- **Issue #21**: Optimize String Operations with strings.Builder (2-3 hours)
- **Issue #22**: Implement Struct Embedding for Type Composition (4-6 hours)
- **Issue #23**: Add Profiling and Metrics Collection (6-8 hours)

---

## Epic 3: Testing & Validation

### Issue #24: Security Test Suite
**Title**: Comprehensive security testing for all remediation
**Priority**: P1-High
**Labels**: `type: testing`, `epic: testing`, `effort: large`
**Milestone**: Milestone 4
**Estimated Effort**: 16-20 hours

**Test Coverage**:
- [ ] Token encryption/decryption tests
- [ ] TLS configuration validation tests
- [ ] Secret input security tests
- [ ] Database encryption tests
- [ ] CSRF token randomness tests
- [ ] Input validation fuzzing
- [ ] Error sanitization tests
- [ ] Session timeout tests
- [ ] Path traversal attack tests
- [ ] Rate limit handling tests

**Tools**:
- `go test -race` (race detector)
- `go test -fuzz` (fuzzing)
- `gosec` (security scanner)
- `go-cve-check` (vulnerability scanner)

**Success Criteria**:
- 90%+ test coverage on security-critical code
- 0 gosec issues
- All fuzz tests pass (100k iterations)
- No known CVEs in dependencies

---

### Issue #25: Performance Benchmark Suite
**Title**: Comprehensive benchmarking for all optimizations
**Priority**: P1-High
**Labels**: `type: testing`, `type: performance`, `epic: testing`, `effort: large`
**Milestone**: Milestone 4
**Estimated Effort**: 16-20 hours

**Benchmarks**:
- [ ] API client throughput (requests/sec)
- [ ] Connection pool efficiency
- [ ] JSON marshaling allocations
- [ ] Batch operation concurrency
- [ ] Audit log pipeline throughput
- [ ] Database query performance
- [ ] Memory usage profiling
- [ ] CPU profiling

**Tools**:
- `go test -bench` (benchmarks)
- `go test -benchmem` (memory allocations)
- `go tool pprof` (profiling)
- `benchstat` (statistical comparison)

**Success Criteria**:
- 30%+ improvement in API throughput
- 40%+ reduction in allocations
- 3x+ improvement in batch operations
- 4x+ improvement in audit loading
- Memory usage within acceptable bounds

**Benchmark Format**:
```go
func BenchmarkAPIClient(b *testing.B) {
    // Baseline benchmark
}

func BenchmarkAPIClientOptimized(b *testing.B) {
    // Optimized version
}

// Run: benchstat baseline.txt optimized.txt
```

---

### Issue #26: Integration Testing
**Title**: End-to-end integration tests for all commands
**Priority**: P1-High
**Labels**: `type: testing`, `epic: testing`, `effort: xl`
**Milestone**: Milestone 4
**Estimated Effort**: 24-32 hours

**Test Scenarios**:
- [ ] Full authentication flow (login → command → logout)
- [ ] Endpoint setup and configuration
- [ ] Collection CRUD operations
- [ ] Storage gateway management
- [ ] User credential workflows
- [ ] Audit log end-to-end
- [ ] Batch operations with errors
- [ ] Token refresh and expiration
- [ ] TLS certificate validation
- [ ] Rate limit handling

**Test Environment**:
- Mock GCS API server (httptest)
- Mock OAuth server
- Test SQLite database
- Test keyring implementation

**Success Criteria**:
- All 83 commands have integration tests
- Tests run in CI/CD pipeline
- Tests complete in <5 minutes

---

### Issue #27: Load Testing
**Title**: Load testing to validate performance under stress
**Priority**: P2-Medium
**Labels**: `type: testing`, `type: performance`, `epic: testing`, `effort: medium`
**Milestone**: Milestone 4
**Estimated Effort**: 8-12 hours

**Test Scenarios**:
- [ ] 1000 concurrent API requests
- [ ] 10,000 batch delete operations
- [ ] 100,000 audit log records
- [ ] 1000 database queries/sec
- [ ] Connection pool exhaustion
- [ ] Memory usage under load

**Tools**:
- `vegeta` (HTTP load testing)
- Custom Go load generators
- Resource monitoring (pprof)

**Success Criteria**:
- No crashes under load
- Performance degrades gracefully
- Memory usage remains bounded
- Connection pool handles burst traffic

---

### Issue #28: Documentation Updates
**Title**: Update all documentation for security and performance changes
**Priority**: P1-High
**Labels**: `type: documentation`, `epic: testing`, `effort: large`
**Milestone**: Milestone 4
**Estimated Effort**: 12-16 hours

**Documentation Needed**:
- [ ] Security architecture document
- [ ] HIPAA compliance guide
- [ ] Configuration reference (security settings)
- [ ] Migration guide (secret handling UX change)
- [ ] Performance tuning guide
- [ ] Troubleshooting guide (keyring, TLS, etc.)
- [ ] API reference updates
- [ ] README updates
- [ ] CHANGELOG for v2.0

**Files to create/update**:
- `docs/SECURITY.md`
- `docs/HIPAA_COMPLIANCE.md`
- `docs/MIGRATION_V2.md`
- `docs/PERFORMANCE.md`
- `docs/CONFIGURATION.md`
- `README.md`
- `CHANGELOG.md`

---

## Issue Dependencies

```
Security Critical Path:
#1 (Token Encryption) ──┬──> #4 (Audit DB Encryption)
                        └──> #8 (Session Timeouts)
                        └──> #12 (Token Rotation)

#2 (TLS Hardening) ────────> #9 (Rate Limiting)

#3 (Secret Input) ─────────> #28 (Documentation)

#6 (Input Validation) ─────> #10 (Path Traversal)

Performance Critical Path:
#14 (Connection Pool) ─────> #9 (Rate Limiting)

#15 (Buffer Pool) ─────────> #16 (Batch Concurrency)

#17 (Audit Pipeline) ──────> #18 (DB Optimization)

Testing Dependencies:
#1-#13 (Security) ─────────> #24 (Security Tests)

#14-#23 (Performance) ─────> #25 (Benchmarks)

#24 + #25 ────────────────> #26 (Integration Tests)

#26 ──────────────────────> #27 (Load Tests)

All Issues ────────────────> #28 (Documentation)
```

---

## CI/CD Pipeline Updates

### New GitHub Actions Workflows

#### `.github/workflows/security.yml`
```yaml
name: Security Tests
on: [push, pull_request]
jobs:
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run gosec
        run: gosec ./...
      - name: Run go-cve-check
        run: go install github.com/google/go-cve-check@latest && go-cve-check ./...
      - name: Run tests with race detector
        run: go test -race ./...
      - name: Run fuzzing (short)
        run: go test -fuzz=. -fuzztime=30s ./...
```

#### `.github/workflows/benchmarks.yml`
```yaml
name: Performance Benchmarks
on: [push, pull_request]
jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run benchmarks
        run: |
          go test -bench=. -benchmem -run=^$ ./... | tee benchmark.txt
      - name: Store benchmark results
        uses: benchmark-action/github-action-benchmark@v1
        with:
          tool: 'go'
          output-file-path: benchmark.txt
          github-token: ${{ secrets.GITHUB_TOKEN }}
          auto-push: true
```

---

## Issue Templates

### Security Issue Template (`.github/ISSUE_TEMPLATE/security.md`)
```markdown
---
name: Security Issue
about: Report a security vulnerability or compliance gap
title: '[SECURITY] '
labels: 'type: security'
---

## Security Finding

**Severity**: [Critical/High/Medium/Low]
**NIST 800-53 Control**: [e.g., SC-28]
**HIPAA Applicability**: [Yes/No]

## Description
[Clear description of the security issue]

## Current Behavior
[What is currently happening]

## Expected Behavior
[What should happen for compliance]

## Proposed Solution
[How to fix the issue]

## Files Affected
- [ ] file1.go
- [ ] file2.go

## Testing
- [ ] Unit tests
- [ ] Security tests
- [ ] Manual verification

## Documentation
- [ ] Security docs updated
- [ ] Configuration guide updated
```

### Performance Issue Template (`.github/ISSUE_TEMPLATE/performance.md`)
```markdown
---
name: Performance Optimization
about: Propose a performance improvement
title: '[PERF] '
labels: 'type: performance'
---

## Optimization Opportunity

**Priority**: [High/Medium/Low]
**Expected Impact**: [e.g., 30% faster, 50% fewer allocations]
**Estimated Effort**: [hours]

## Current Performance
[Benchmark results or profiling data]

## Proposed Optimization
[Description of the optimization]

## Implementation Details
[Technical approach]

## Benchmarks
```go
// Include benchmark code
```

## Files Affected
- [ ] file1.go
- [ ] file2.go

## Testing
- [ ] Benchmarks added
- [ ] Correctness verified
- [ ] Memory profiling
```

---

## Sprint Planning (Example)

### Sprint 1 (Week 1-2): Critical Security - Milestone 1
**Goal**: Address all CRITICAL severity security issues

**Issues**:
- #1: Token Encryption (16h) - @developer1
- #2: TLS Hardening (8h) - @developer2
- #3: Secret Input (16h) - @developer1
- #5: CSRF Token (4h) - @developer2

**Total**: 44 hours (2 developers × 40 hours/week × 2 weeks = 160 hours available)

**Daily Standup Focus**:
- Token encryption integration
- Keyring fallback handling
- Secret input UX validation
- TLS testing

**Sprint Review**:
- Demo token encryption
- Demo secure secret input
- TLS configuration review

---

### Sprint 2 (Week 3-4): High Priority Security + Audit DB - Milestone 2
**Goal**: Complete HIGH severity security issues

**Issues**:
- #4: Audit DB Encryption (16h) - @developer1
- #6: Input Validation (20h) - @developer1 + @developer2
- #7: Error Sanitization (12h) - @developer2
- #8: Session Timeouts (12h) - @developer1
- #9: Rate Limiting (8h) - @developer2
- #10: Path Traversal (6h) - @developer2

**Total**: 74 hours

---

### Sprint 3 (Week 5-6): High Priority Optimizations - Milestone 3
**Goal**: Implement high-impact performance improvements

**Issues**:
- #14: HTTP Connection Pooling (4h) - @developer2
- #15: JSON Buffer Pooling (4h) - @developer2
- #16: Batch Concurrency (8h) - @developer1
- #17: Audit Pipeline (16h) - @developer1
- #18: DB Optimization (4h) - @developer2

**Total**: 36 hours

---

### Sprint 4 (Week 7-8): Testing & Documentation - Milestone 4
**Goal**: Comprehensive testing and documentation

**Issues**:
- #24: Security Test Suite (20h) - @developer1 + @developer2
- #25: Performance Benchmarks (20h) - @developer1 + @developer2
- #26: Integration Tests (32h) - @developer1 + @developer2
- #28: Documentation (16h) - @developer1 + @developer2

**Total**: 88 hours

**Release**: v2.0.0 - Security & Performance Release

---

## Release Checklist

### v2.0.0 Release
- [ ] All critical security issues resolved
- [ ] All high priority security issues resolved
- [ ] All high priority optimizations implemented
- [ ] Security test suite passes
- [ ] Performance benchmarks meet targets
- [ ] Integration tests pass
- [ ] Documentation complete
- [ ] Migration guide published
- [ ] CHANGELOG updated
- [ ] Release notes drafted
- [ ] Security audit completed (external review recommended)
- [ ] Load testing completed
- [ ] Beta testing period (2 weeks recommended)

**Breaking Changes**:
- Secret input via CLI arguments removed (use stdin/prompt/env)

**Upgrade Path**:
1. Backup token and config files
2. Install v2.0.0
3. Run token migration: `globus-connect-server auth migrate-tokens`
4. Update scripts to use new secret input methods
5. Review security configuration

---

## Metrics & KPIs

### Security Metrics
- **Target**: 0 critical vulnerabilities
- **Target**: 0 high vulnerabilities
- **Target**: 90%+ test coverage on security code
- **Target**: Pass external security audit

### Performance Metrics
- **API Throughput**: 30-50% improvement
- **Memory Allocations**: 40-60% reduction
- **Batch Operations**: 3-5x faster
- **Audit Loading**: 4-6x faster
- **Database Queries**: 25-40% faster

### Quality Metrics
- **Test Coverage**: >80% overall, >90% security-critical
- **Linter Issues**: 0
- **Go Report Card**: A+
- **Documentation**: 100% of new features documented

---

## Summary Statistics

**Total Issues**: 28
- Security: 13 issues (5 critical, 8 high)
- Performance: 10 issues (5 high, 5 medium)
- Testing: 5 issues

**Total Estimated Effort**: 242-318 hours
- Security: 120-160 hours (15-20 days)
- Performance: 64-96 hours (8-12 days)
- Testing: 58-82 hours (7-10 days)

**Timeline**: 6-8 weeks
- Sprint 1-2: Security (4 weeks)
- Sprint 3: Optimizations (2 weeks)
- Sprint 4: Testing & Documentation (2 weeks)

**Team**: 2 developers
**Working Hours**: 160 hours/developer/month
**Total Capacity**: 320 hours/month

**Expected Outcomes**:
- HIPAA/PHI ready deployment
- 2-5x performance improvement
- Production-ready v2.0.0 release
