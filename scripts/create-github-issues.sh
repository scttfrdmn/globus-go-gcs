#!/bin/bash

# Script to create GitHub issues for Security Remediation & Optimization project
# Prerequisites: gh CLI tool installed and authenticated
# Usage: ./scripts/create-github-issues.sh

set -e

REPO="scttfrdmn/globus-go-gcs"

echo "Creating GitHub issues for Security & Performance v2.0 project..."
echo "Repository: $REPO"
echo ""

# Check if gh is installed
if ! command -v gh &> /dev/null; then
    echo "Error: gh CLI tool is not installed"
    echo "Install from: https://cli.github.com/"
    exit 1
fi

# Check if authenticated
if ! gh auth status &> /dev/null; then
    echo "Error: Not authenticated with GitHub"
    echo "Run: gh auth login"
    exit 1
fi

echo "Creating milestones..."

# Create milestones
gh api repos/$REPO/milestones -X POST -f title="Milestone 1: Critical Security Fixes" -f due_on="$(date -u -v+2w +%Y-%m-%dT%H:%M:%SZ)" -f description="Address all CRITICAL severity security findings" || echo "Milestone 1 may already exist"

gh api repos/$REPO/milestones -X POST -f title="Milestone 2: High Priority Security" -f due_on="$(date -u -v+4w +%Y-%m-%dT%H:%M:%SZ)" -f description="Address all HIGH severity security findings" || echo "Milestone 2 may already exist"

gh api repos/$REPO/milestones -X POST -f title="Milestone 3: High Priority Optimizations" -f due_on="$(date -u -v+6w +%Y-%m-%dT%H:%M:%SZ)" -f description="Implement high-impact performance improvements" || echo "Milestone 3 may already exist"

gh api repos/$REPO/milestones -X POST -f title="Milestone 4: Testing & Documentation" -f due_on="$(date -u -v+8w +%Y-%m-%dT%H:%M:%SZ)" -f description="Comprehensive validation and documentation" || echo "Milestone 4 may already exist"

echo "Milestones created."
echo ""

# Function to create an issue
create_issue() {
    local title="$1"
    local body="$2"
    local labels="$3"
    local milestone="$4"

    echo "Creating issue: $title"
    gh issue create \
        --repo "$REPO" \
        --title "$title" \
        --body "$body" \
        --label "$labels" \
        --milestone "$milestone" || echo "Failed to create issue: $title"
}

echo "Creating Epic: Security Remediation issues..."

# Issue #1: Token Encryption
create_issue \
    "[SECURITY] Implement AES-256-GCM encryption for stored OAuth tokens" \
    "## Security Finding

**Severity**: P0-Critical
**NIST 800-53 Control**: SC-28 (Protection of Information at Rest)
**HIPAA Applicability**: Yes

## Description
Currently OAuth tokens are stored in ~/.globus-connect-server/tokens/<profile>.json with only file permissions (0600) for protection. Tokens must be encrypted at rest using AES-256-GCM.

## Current Behavior
- Tokens stored in plaintext JSON files
- Protected only by file permissions
- Vulnerable if system is compromised

## Security Impact
- CRITICAL: Plaintext tokens accessible if attacker gains file system access
- Fails HIPAA Security Rule § 164.312(a)(2)(iv) - encryption requirements
- Non-compliant with NIST 800-53 SC-28

## Expected Behavior
- Tokens encrypted at rest with AES-256-GCM
- Encryption key derived from system keyring
- Key rotation support
- Migration for existing plaintext tokens

## Implementation Details
- [ ] Integrate github.com/zalando/go-keyring for key management
- [ ] Implement AES-256-GCM encryption wrapper in internal/auth/encryption.go
- [ ] Update LoadToken/SaveToken in internal/auth/tokens.go
- [ ] Add key rotation support
- [ ] Handle keyring unavailable fallback (prompt for passphrase)
- [ ] Create migration tool for existing tokens

## Files Affected
- internal/auth/tokens.go
- New: internal/auth/encryption.go
- New: internal/auth/keyring.go

## Estimated Effort
Large - 12-16 hours" \
    "P0-Critical,type: security,area: auth,effort: large,epic: security-remediation" \
    "Milestone 1: Critical Security Fixes"

# Issue #2: TLS Configuration
create_issue \
    "[SECURITY] Enforce TLS 1.2+ with secure cipher suites" \
    "## Security Finding

**Severity**: P0-Critical
**NIST 800-53 Control**: SC-8 (Transmission Confidentiality), SC-13 (Cryptographic Protection)
**HIPAA Applicability**: Yes

## Description
HTTP client uses default TLS configuration without enforcing minimum TLS version or approved cipher suites.

## Current Behavior
- Uses Go's default TLS configuration
- Accepts TLS 1.0/1.1 connections
- No cipher suite restrictions

## Security Impact
- Vulnerable to downgrade attacks
- May use weak cipher suites
- Non-compliant with NIST 800-53 SC-8, SC-13

## Implementation Details
- [ ] Create custom TLS config with MinVersion: TLS 1.2
- [ ] Define approved cipher suite list
- [ ] Add certificate validation options
- [ ] Implement optional certificate pinning
- [ ] Add TLS configuration options

## Files Affected
- pkg/gcs/client.go
- pkg/gcs/options.go
- New: pkg/gcs/tls.go

## Estimated Effort
Medium - 6-8 hours" \
    "P0-Critical,type: security,area: tls,area: api-client,effort: medium,epic: security-remediation" \
    "Milestone 1: Critical Security Fixes"

# Issue #3: Secrets in CLI
create_issue \
    "[SECURITY] 🚨 BREAKING CHANGE: Remove secrets from CLI arguments" \
    "## 🚨 BREAKING CHANGE - User-Visible Impact

This issue implements a **security-critical breaking change** that affects 5 commands.

**📖 Full Documentation**: See [BREAKING_CHANGES_V2.md](../BREAKING_CHANGES_V2.md) for complete migration guide.

---

## Security Finding

**Severity**: P0-Critical
**NIST 800-53 Control**: IA-5(7) (No Embedded Passwords)
**HIPAA Applicability**: Yes
**⚠️ UX Impact**: **YES - Changes user interaction pattern**

## Description
Multiple commands accept secrets via CLI arguments, which are visible in process listings and shell history. This is a **critical security vulnerability** that makes v1.x unsuitable for HIPAA/PHI environments.

## Current Behavior (INSECURE)
\`\`\`bash
# v1.x - Secrets visible in process list
globus-connect-server user-credential s3-keys add \\
  --secret-access-key wJalrXUtnFEMI/K7MDENG...

# Any user can see:
$ ps aux | grep secret
user  12345  ... --secret-access-key wJalrXUtnFEMI/K7MDENG...
\`\`\`

## Security Impact
- **CRITICAL**: Secrets visible in ps output to all users
- **CRITICAL**: Secrets stored in shell history files
- **CRITICAL**: Secrets may be logged by system monitoring
- **Non-compliant**: NIST 800-53 IA-5(7), HIPAA § 164.312(a)(2)(iv), PCI DSS 8.2.1

## New Behavior (SECURE)

### Option 1: Interactive Prompt (Recommended)
\`\`\`bash
# v2.0 - Prompt for secret with hidden input
globus-connect-server user-credential s3-keys add \\
  --access-key-id AKIA...
# Prompts: Enter secret access key: ********
\`\`\`

### Option 2: Stdin (For Automation)
\`\`\`bash
# v2.0 - Pipe secret from secure source
echo \"\$SECRET\" | globus-connect-server user-credential s3-keys add \\
  --access-key-id AKIA... \\
  --secret-stdin
\`\`\`

### Option 3: Environment Variable
\`\`\`bash
# v2.0 - Read from environment
export GLOBUS_SECRET_VALUE=\"...\"
globus-connect-server user-credential s3-keys add \\
  --access-key-id AKIA... \\
  --secret-env
\`\`\`

## Implementation Details
- [ ] Create secure input utility in internal/input/secure.go
- [ ] Implement interactive prompt with hidden input (golang.org/x/term)
- [ ] Add --secret-stdin flag for pipe input
- [ ] Add --secret-env flag for environment variable
- [ ] Remove old --secret-* flags completely
- [ ] Update all affected commands
- [ ] Add validation for secret format/length
- [ ] Add confirmation prompt for interactive mode

## Affected Commands (5)
1. \`user-credential s3-keys add\` - removes --secret-access-key
2. \`user-credential s3-keys update\` - removes --secret-access-key
3. \`user-credential activescale-create\` - removes --password
4. \`oidc create\` - removes --client-secret
5. \`oidc update\` - removes --client-secret

## Documentation Requirements
- [ ] Update BREAKING_CHANGES_V2.md with examples
- [ ] Update command help text
- [ ] Add migration examples for shell scripts
- [ ] Add migration examples for CI/CD (GitHub Actions, GitLab CI)
- [ ] Add migration examples for IaC (Ansible, Terraform)
- [ ] Update README.md with prominent warning
- [ ] Create detection script for finding old usage

## Testing Requirements
- [ ] Unit tests for secure input functions
- [ ] Integration tests for all three input methods
- [ ] Test empty/invalid secret handling
- [ ] Test confirmation mismatch handling
- [ ] Test in CI/CD environment
- [ ] Manual testing on Linux/macOS/Windows

## Migration Support
- [ ] Provide migration guide with clear examples
- [ ] Provide detection script to find old usage patterns
- [ ] Show deprecation warning in v2.0-beta
- [ ] Document in release notes prominently

## Estimated Effort
Large - 12-16 hours

## ⚠️ Release Blockers
- [ ] Migration guide complete and reviewed
- [ ] All examples tested on multiple platforms
- [ ] Breaking change documented in release notes
- [ ] Beta testing period completed (2 weeks minimum)" \
    "P0-Critical,type: security,area: cli,effort: large,epic: security-remediation,breaking-change" \
    "Milestone 1: Critical Security Fixes"

# Issue #4: Audit Database Encryption
create_issue \
    "[SECURITY] Implement SQLCipher for encrypted audit database" \
    "## Security Finding

**Severity**: P0-Critical
**NIST 800-53 Control**: AU-9 (Protection of Audit Information), SC-28
**HIPAA Applicability**: Yes

## Description
Audit logs stored in SQLite database are unencrypted.

## Security Impact
- Audit logs contain sensitive information (PHI access events)
- Unencrypted database violates AU-9 and SC-28
- HIPAA audit trail must be protected

## Implementation Details
- [ ] Replace modernc.org/sqlite with SQLCipher
- [ ] Derive encryption key from system keyring
- [ ] Migrate existing databases to encrypted format
- [ ] Add secure database deletion

## Files Affected
- internal/commands/audit/audit.go
- All audit command files
- New: internal/commands/audit/migration.go

## Dependencies
- Depends on #1 (keyring integration)

## Estimated Effort
Large - 12-16 hours" \
    "P0-Critical,type: security,area: audit,area: database,effort: large,epic: security-remediation" \
    "Milestone 1: Critical Security Fixes"

# Issue #5: CSRF Token
create_issue \
    "[SECURITY] Strengthen CSRF token generation" \
    "## Security Finding

**Severity**: P0-Critical
**NIST 800-53 Control**: IA-9, SC-13
**HIPAA Applicability**: Yes

## Description
CSRF state token uses timestamp-based generation which is predictable.

## Current Code
\`\`\`go
state := fmt.Sprintf(\"state-%d\", time.Now().Unix())
\`\`\`

## Implementation Details
- [ ] Use crypto/rand for 32 bytes of randomness
- [ ] Base64 URL-safe encode
- [ ] Add state expiration (5 minutes)
- [ ] Implement state validation

## Files Affected
- internal/commands/auth/login.go

## Estimated Effort
Small - 2-4 hours" \
    "P0-Critical,type: security,area: auth,effort: small,epic: security-remediation" \
    "Milestone 1: Critical Security Fixes"

echo ""
echo "Creating Epic: Performance Optimization issues..."

# Issue #14: HTTP Connection Pooling
create_issue \
    "[PERF] Implement HTTP client connection pooling" \
    "## Optimization Opportunity

**Priority**: P1-High
**Expected Impact**: 30-50% faster API calls
**Estimated Effort**: Small - 3-4 hours

## Current Performance
Each GCS client creates a new HTTP client with default settings.

## Proposed Optimization
Implement shared HTTP transport with connection pooling:
- MaxIdleConns: 100
- MaxIdleConnsPerHost: 10
- IdleConnTimeout: 90s
- Enable HTTP/2

## Implementation Details
- [ ] Create shared transport in pkg/gcs/transport.go
- [ ] Configure connection pool parameters
- [ ] Add connection metrics

## Files Affected
- pkg/gcs/client.go
- pkg/gcs/options.go
- New: pkg/gcs/transport.go

## Performance Target
30-50% reduction in API call latency

## User Impact
**UX Changes**: None
**Functionality Changes**: None" \
    "P1-High,type: performance,area: api-client,effort: small,epic: go-optimizations" \
    "Milestone 3: High Priority Optimizations"

# Issue #15: JSON Buffer Pooling
create_issue \
    "[PERF] Implement sync.Pool for JSON marshaling buffers" \
    "## Optimization Opportunity

**Priority**: P1-High
**Expected Impact**: 40-60% fewer allocations
**Estimated Effort**: Small - 3-4 hours

## Current Performance
70+ calls to json.Marshal allocate new buffers each time.

## Proposed Optimization
Use sync.Pool to reuse JSON encoding buffers.

## Implementation Details
- [ ] Create buffer pool with sync.Pool
- [ ] Create helper function for pooled JSON encoding
- [ ] Replace json.Marshal in hot paths
- [ ] Add buffer pool metrics

## Files Affected
- pkg/gcs/client.go
- New: pkg/gcs/pool.go

## Performance Target
40-60% reduction in allocations

## User Impact
**UX Changes**: None
**Functionality Changes**: None" \
    "P1-High,type: performance,area: api-client,effort: small,epic: go-optimizations" \
    "Milestone 3: High Priority Optimizations"

# Issue #16: Batch Concurrency
create_issue \
    "[PERF] Implement goroutine worker pool for batch operations" \
    "## Optimization Opportunity

**Priority**: P1-High
**Expected Impact**: 3-5x faster batch operations
**Estimated Effort**: Medium - 6-8 hours

## Current Performance
Batch delete operations are sequential.

## Proposed Optimization
Implement concurrent deletion with goroutine worker pool:
- Configurable worker count (default: 10)
- Channel-based job queue
- Progress reporting with atomic counters
- Error aggregation
- Rate limiting integration

## Implementation Details
- [ ] Create worker pool pattern
- [ ] Implement job queue with channels
- [ ] Add progress reporting
- [ ] Add error aggregation
- [ ] Add graceful cancellation

## Files Affected
- internal/commands/collection/batch_delete.go
- internal/commands/role/batch_delete.go
- New: internal/concurrent/worker_pool.go

## Performance Target
3-5x faster for 100+ item batches

## User Impact
**UX Changes**: None
**Functionality Changes**: None" \
    "P1-High,type: performance,area: cli,effort: medium,epic: go-optimizations" \
    "Milestone 3: High Priority Optimizations"

# Issue #17: Audit Pipeline
create_issue \
    "[PERF] Implement channel-based pipeline for audit log processing" \
    "## Optimization Opportunity

**Priority**: P1-High
**Expected Impact**: 4-6x faster audit load
**Estimated Effort**: Large - 12-16 hours

## Current Performance
Audit log loading is sequential (fetch → parse → insert).

## Proposed Optimization
Implement pipeline pattern with concurrent stages:
- Stage 1: Fetch from API (producer)
- Stage 2: Parse and transform (10 workers)
- Stage 3: Batch insert to DB (consumer, 100 records/transaction)

## Implementation Details
- [ ] Create pipeline architecture
- [ ] Implement fetch stage
- [ ] Implement parse worker pool
- [ ] Implement batch insert stage
- [ ] Add backpressure handling
- [ ] Add progress reporting

## Files Affected
- internal/commands/audit/load.go
- New: internal/commands/audit/pipeline.go

## Performance Target
4-6x faster for large datasets (10k+ records)

## User Impact
**UX Changes**: None
**Functionality Changes**: None" \
    "P1-High,type: performance,area: audit,area: database,effort: large,epic: go-optimizations" \
    "Milestone 3: High Priority Optimizations"

# Issue #18: Database Optimization
create_issue \
    "[PERF] Add indexes and optimize SQLite queries" \
    "## Optimization Opportunity

**Priority**: P1-High
**Expected Impact**: 25-40% faster queries
**Estimated Effort**: Small - 3-4 hours

## Proposed Optimization
- Add composite indexes for common query patterns
- Use prepared statements
- Enable SQLite performance pragmas (WAL, mmap)
- Consider FTS5 for full-text search

## Implementation Details
- [ ] Add composite indexes
- [ ] Implement prepared statement caching
- [ ] Configure SQLite pragmas for performance
- [ ] Add query result caching

## Files Affected
- internal/commands/audit/audit.go
- internal/commands/audit/query.go

## Performance Target
25-40% faster query execution

## User Impact
**UX Changes**: None
**Functionality Changes**: None" \
    "P1-High,type: performance,area: database,effort: small,epic: go-optimizations" \
    "Milestone 3: High Priority Optimizations"

echo ""
echo "Creating Epic: Testing issues..."

# Issue #24: Security Test Suite
create_issue \
    "[TEST] Comprehensive security testing for all remediation" \
    "## Testing Task

**Type**: Security Tests
**Priority**: P1-High
**Estimated Effort**: 16-20 hours

## Test Scope
Comprehensive security testing for all NIST 800-53 remediation work.

## Test Coverage Goals
- [ ] 90%+ test coverage on security-critical code
- [ ] All security features tested
- [ ] Fuzzing for input validation
- [ ] Race condition detection

## Test Scenarios
- Token encryption/decryption
- TLS configuration validation
- Secret input security
- Database encryption
- CSRF token randomness
- Input validation fuzzing
- Error sanitization
- Session timeout enforcement
- Path traversal protection
- Rate limit handling

## Tools
- go test -race
- go test -fuzz
- gosec
- go-cve-check

## Success Criteria
- [ ] 90%+ test coverage on security code
- [ ] 0 gosec issues
- [ ] All fuzz tests pass (100k iterations)
- [ ] No known CVEs in dependencies

## Dependencies
- Depends on #1-#13 (security implementations)" \
    "P1-High,type: testing,epic: testing,effort: large" \
    "Milestone 4: Testing & Documentation"

# Issue #25: Performance Benchmarks
create_issue \
    "[TEST] Comprehensive benchmarking for all optimizations" \
    "## Testing Task

**Type**: Performance Benchmarks
**Priority**: P1-High
**Estimated Effort**: 16-20 hours

## Test Scope
Benchmark all performance optimizations against baseline.

## Benchmarks Needed
- API client throughput (requests/sec)
- Connection pool efficiency
- JSON marshaling allocations
- Batch operation concurrency
- Audit log pipeline throughput
- Database query performance
- Memory usage profiling
- CPU profiling

## Tools
- go test -bench -benchmem
- go tool pprof
- benchstat

## Success Criteria
- [ ] 30%+ improvement in API throughput
- [ ] 40%+ reduction in allocations
- [ ] 3x+ improvement in batch operations
- [ ] 4x+ improvement in audit loading

## Dependencies
- Depends on #14-#23 (optimization implementations)" \
    "P1-High,type: testing,type: performance,epic: testing,effort: large" \
    "Milestone 4: Testing & Documentation"

# Issue #26: Integration Testing
create_issue \
    "[TEST] End-to-end integration tests for all commands" \
    "## Testing Task

**Type**: Integration Tests
**Priority**: P1-High
**Estimated Effort**: 24-32 hours

## Test Scope
End-to-end testing of all 83 commands.

## Test Scenarios
- Full authentication flow
- Endpoint setup and configuration
- Collection CRUD operations
- Storage gateway management
- User credential workflows
- Audit log end-to-end
- Batch operations with errors
- Token refresh and expiration
- TLS certificate validation
- Rate limit handling

## Test Environment
- Mock GCS API server (httptest)
- Mock OAuth server
- Test SQLite database
- Test keyring implementation

## Success Criteria
- [ ] All 83 commands have integration tests
- [ ] Tests run in CI/CD pipeline
- [ ] Tests complete in <5 minutes

## Dependencies
- Depends on #24, #25" \
    "P1-High,type: testing,epic: testing,effort: xl" \
    "Milestone 4: Testing & Documentation"

# Issue #28: Documentation
create_issue \
    "[DOC] Update all documentation for security and performance changes" \
    "## Documentation Task

**Priority**: P1-High
**Estimated Effort**: 12-16 hours

## Documentation Needed
- [ ] Security architecture document
- [ ] HIPAA compliance guide
- [ ] Configuration reference (security settings)
- [ ] Migration guide (secret handling UX change)
- [ ] Performance tuning guide
- [ ] Troubleshooting guide
- [ ] API reference updates
- [ ] README updates
- [ ] CHANGELOG for v2.0

## Files to Create/Update
- docs/SECURITY.md
- docs/HIPAA_COMPLIANCE.md
- docs/MIGRATION_V2.md
- docs/PERFORMANCE.md
- docs/CONFIGURATION.md
- README.md
- CHANGELOG.md

## Dependencies
- Depends on all implementation issues" \
    "P1-High,type: documentation,epic: testing,effort: large" \
    "Milestone 4: Testing & Documentation"

echo ""
echo "✅ GitHub issues created successfully!"
echo ""
echo "Next steps:"
echo "1. Review issues at: https://github.com/$REPO/issues"
echo "2. Create project board: gh project create --title 'Security & Performance v2.0'"
echo "3. Add issues to project board"
echo "4. Start with Milestone 1 (Critical Security Fixes)"
echo ""
echo "To view milestones: gh api repos/$REPO/milestones"
echo "To view issues: gh issue list --repo $REPO"
