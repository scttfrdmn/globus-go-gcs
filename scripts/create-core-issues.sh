#!/bin/bash

# Script to create core v2.0 issues with simplified bodies
# Prerequisites: gh CLI tool installed and authenticated
# Usage: ./scripts/create-core-issues.sh

set -e

REPO="scttfrdmn/globus-go-gcs"

echo "Creating core v2.0 issues..."
echo "Repository: $REPO"
echo ""

# Milestone 1: Critical Security Fixes
echo "==> Creating Milestone 1 issues (Critical Security)"

gh issue create --repo "$REPO" \
  --title "[SECURITY] Implement AES-256-GCM encryption for stored OAuth tokens" \
  --milestone "Milestone 1: Critical Security Fixes" \
  --label "P0-Critical,type: security,area: auth,effort: large,epic: security-remediation" \
  --body "**Priority**: P0-Critical
**NIST 800-53**: SC-28 (Protection of Information at Rest)
**HIPAA**: Yes
**Effort**: 12-16 hours

## Description
Currently OAuth tokens are stored in plaintext JSON files with only file permissions (0600). Must implement AES-256-GCM encryption at rest.

## Implementation
- [ ] Integrate github.com/zalando/go-keyring
- [ ] Implement AES-256-GCM encryption in internal/auth/encryption.go
- [ ] Update LoadToken/SaveToken in internal/auth/tokens.go
- [ ] Add key rotation support
- [ ] Create migration tool

## Files
- internal/auth/tokens.go
- New: internal/auth/encryption.go
- New: internal/auth/keyring.go

See GITHUB_PROJECT_PLAN.md for full details."

gh issue create --repo "$REPO" \
  --title "[SECURITY] Enforce TLS 1.2+ with secure cipher suites" \
  --milestone "Milestone 1: Critical Security Fixes" \
  --label "P0-Critical,type: security,area: tls,area: api-client,effort: medium,epic: security-remediation" \
  --body "**Priority**: P0-Critical
**NIST 800-53**: SC-8, SC-13
**HIPAA**: Yes
**Effort**: 6-8 hours

## Description
HTTP client uses default TLS configuration. Must enforce TLS 1.2+ with approved cipher suites.

## Implementation
- [ ] Create custom TLS config with MinVersion: TLS 1.2
- [ ] Define approved cipher suite list
- [ ] Add certificate validation options
- [ ] Add TLS configuration options

## Files
- pkg/gcs/client.go
- pkg/gcs/options.go
- New: pkg/gcs/tls.go

See GITHUB_PROJECT_PLAN.md for full details."

gh issue create --repo "$REPO" \
  --title "[SECURITY] 🚨 BREAKING CHANGE: Remove secrets from CLI arguments" \
  --milestone "Milestone 1: Critical Security Fixes" \
  --label "P0-Critical,type: security,area: cli,effort: large,epic: security-remediation,breaking-change" \
  --body "**⚠️ BREAKING CHANGE - User-Visible Impact**

**Priority**: P0-Critical
**NIST 800-53**: IA-5(7)
**HIPAA**: Yes
**Effort**: 12-16 hours

## Description
Secrets via CLI arguments are visible in process listings and shell history. Critical security vulnerability.

**📖 Full Migration Guide**: BREAKING_CHANGES_V2.md

## What's Changing
Removing flags: --secret-access-key, --client-secret, --password

## New Secure Methods
1. Interactive prompt (recommended)
2. --secret-stdin flag
3. --secret-env flag

## Affected Commands (5)
- user-credential s3-keys add/update
- user-credential activescale-create
- oidc create/update

See BREAKING_CHANGES_V2.md for complete migration guide."

gh issue create --repo "$REPO" \
  --title "[SECURITY] Implement SQLCipher for encrypted audit database" \
  --milestone "Milestone 1: Critical Security Fixes" \
  --label "P0-Critical,type: security,area: audit,area: database,effort: large,epic: security-remediation" \
  --body "**Priority**: P0-Critical
**NIST 800-53**: AU-9, SC-28
**HIPAA**: Yes
**Effort**: 12-16 hours

## Description
Audit logs in SQLite database are unencrypted. Must use SQLCipher for encryption at rest.

## Implementation
- [ ] Replace modernc.org/sqlite with SQLCipher
- [ ] Derive encryption key from system keyring
- [ ] Migrate existing databases
- [ ] Add secure database deletion

## Files
- internal/commands/audit/audit.go
- All audit command files
- New: internal/commands/audit/migration.go

## Dependencies
Depends on #8 (keyring integration)

See GITHUB_PROJECT_PLAN.md for full details."

gh issue create --repo "$REPO" \
  --title "[SECURITY] Strengthen CSRF token generation" \
  --milestone "Milestone 1: Critical Security Fixes" \
  --label "P0-Critical,type: security,area: auth,effort: small,epic: security-remediation" \
  --body "**Priority**: P0-Critical
**NIST 800-53**: IA-9, SC-13
**HIPAA**: Yes
**Effort**: 2-4 hours

## Description
CSRF state token uses timestamp-based generation (predictable). Must use crypto/rand.

## Current Code
\`\`\`go
state := fmt.Sprintf(\"state-%d\", time.Now().Unix())
\`\`\`

## Implementation
- [ ] Use crypto/rand for 32 bytes
- [ ] Base64 URL-safe encode
- [ ] Add state expiration (5 minutes)
- [ ] Implement state validation

## Files
- internal/commands/auth/login.go

See GITHUB_PROJECT_PLAN.md for full details."

echo ""
echo "==> Creating Milestone 3 issues (High Priority Optimizations)"

gh issue create --repo "$REPO" \
  --title "[PERF] Implement HTTP client connection pooling" \
  --milestone "Milestone 3: High Priority Optimizations" \
  --label "P1-High,type: performance,area: api-client,effort: small,epic: go-optimizations" \
  --body "**Priority**: P1-High
**Expected Impact**: 30-50% faster API calls
**Effort**: 3-4 hours

## Description
Each GCS client creates new HTTP client. Implement shared connection pooling.

## Implementation
- [ ] Create shared transport with MaxIdleConns: 100
- [ ] Configure IdleConnTimeout: 90s
- [ ] Enable HTTP/2
- [ ] Add connection metrics

## Files
- pkg/gcs/client.go
- pkg/gcs/options.go
- New: pkg/gcs/transport.go

See GITHUB_PROJECT_PLAN.md for full details."

gh issue create --repo "$REPO" \
  --title "[PERF] Implement sync.Pool for JSON marshaling buffers" \
  --milestone "Milestone 3: High Priority Optimizations" \
  --label "P1-High,type: performance,area: api-client,effort: small,epic: go-optimizations" \
  --body "**Priority**: P1-High
**Expected Impact**: 40-60% fewer allocations
**Effort**: 3-4 hours

## Description
70+ calls to json.Marshal. Implement buffer pooling with sync.Pool.

## Implementation
- [ ] Create buffer pool with sync.Pool
- [ ] Create helper for pooled JSON encoding
- [ ] Replace json.Marshal in hot paths
- [ ] Add buffer pool metrics

## Files
- pkg/gcs/client.go
- New: pkg/gcs/pool.go

See GITHUB_PROJECT_PLAN.md for full details."

gh issue create --repo "$REPO" \
  --title "[PERF] Implement goroutine worker pool for batch operations" \
  --milestone "Milestone 3: High Priority Optimizations" \
  --label "P1-High,type: performance,area: cli,effort: medium,epic: go-optimizations" \
  --body "**Priority**: P1-High
**Expected Impact**: 3-5x faster batch operations
**Effort**: 6-8 hours

## Description
Batch delete operations are sequential. Implement concurrent deletion with goroutine worker pool.

## Implementation
- [ ] Create worker pool pattern
- [ ] Implement job queue with channels
- [ ] Add progress reporting
- [ ] Add error aggregation
- [ ] Add graceful cancellation

## Files
- internal/commands/collection/batch_delete.go
- internal/commands/role/batch_delete.go
- New: internal/concurrent/worker_pool.go

See GITHUB_PROJECT_PLAN.md for full details."

gh issue create --repo "$REPO" \
  --title "[PERF] Implement channel-based pipeline for audit log processing" \
  --milestone "Milestone 3: High Priority Optimizations" \
  --label "P1-High,type: performance,area: audit,area: database,effort: large,epic: go-optimizations" \
  --body "**Priority**: P1-High
**Expected Impact**: 4-6x faster audit load
**Effort**: 12-16 hours

## Description
Audit log loading is sequential. Implement pipeline with concurrent stages.

## Architecture
- Stage 1: Fetch from API (producer)
- Stage 2: Parse and transform (10 workers)
- Stage 3: Batch insert to DB (consumer)

## Implementation
- [ ] Create pipeline architecture
- [ ] Implement fetch stage
- [ ] Implement parse worker pool
- [ ] Implement batch insert (100 records/txn)
- [ ] Add backpressure handling

## Files
- internal/commands/audit/load.go
- New: internal/commands/audit/pipeline.go

See GITHUB_PROJECT_PLAN.md for full details."

gh issue create --repo "$REPO" \
  --title "[PERF] Add indexes and optimize SQLite queries" \
  --milestone "Milestone 3: High Priority Optimizations" \
  --label "P1-High,type: performance,area: database,effort: small,epic: go-optimizations" \
  --body "**Priority**: P1-High
**Expected Impact**: 25-40% faster queries
**Effort**: 3-4 hours

## Description
Optimize audit queries with indexes and prepared statements.

## Implementation
- [ ] Add composite indexes for common patterns
- [ ] Implement prepared statement caching
- [ ] Configure SQLite pragmas (WAL, mmap)
- [ ] Consider FTS5 for full-text search

## Files
- internal/commands/audit/audit.go
- internal/commands/audit/query.go

See GITHUB_PROJECT_PLAN.md for full details."

echo ""
echo "==> Creating Milestone 4 issues (Testing & Documentation)"

gh issue create --repo "$REPO" \
  --title "[TEST] Comprehensive security testing for all remediation" \
  --milestone "Milestone 4: Testing & Documentation" \
  --label "P1-High,type: testing,epic: testing,effort: large" \
  --body "**Priority**: P1-High
**Effort**: 16-20 hours

## Description
Comprehensive security testing for NIST 800-53 remediation.

## Test Coverage
- Token encryption/decryption
- TLS configuration validation
- Secret input security
- Database encryption
- CSRF token randomness
- Input validation fuzzing
- Error sanitization
- Session timeout
- Path traversal protection
- Rate limit handling

## Success Criteria
- [ ] 90%+ test coverage on security code
- [ ] 0 gosec issues
- [ ] All fuzz tests pass (100k iterations)
- [ ] No known CVEs

See GITHUB_PROJECT_PLAN.md for full details."

gh issue create --repo "$REPO" \
  --title "[TEST] Comprehensive benchmarking for all optimizations" \
  --milestone "Milestone 4: Testing & Documentation" \
  --label "P1-High,type: testing,type: performance,epic: testing,effort: large" \
  --body "**Priority**: P1-High
**Effort**: 16-20 hours

## Description
Benchmark all performance optimizations against baseline.

## Benchmarks
- API client throughput
- Connection pool efficiency
- JSON marshaling allocations
- Batch operation concurrency
- Audit log pipeline
- Database query performance
- Memory/CPU profiling

## Success Criteria
- [ ] 30%+ API throughput improvement
- [ ] 40%+ reduction in allocations
- [ ] 3x+ batch operation improvement
- [ ] 4x+ audit loading improvement

See GITHUB_PROJECT_PLAN.md for full details."

gh issue create --repo "$REPO" \
  --title "[TEST] End-to-end integration tests for all commands" \
  --milestone "Milestone 4: Testing & Documentation" \
  --label "P1-High,type: testing,epic: testing,effort: xl" \
  --body "**Priority**: P1-High
**Effort**: 24-32 hours

## Description
End-to-end integration tests for all 83 commands.

## Test Scenarios
- Full authentication flow
- Endpoint setup/config
- Collection CRUD
- Storage gateway management
- User credential workflows
- Audit log end-to-end
- Batch operations with errors
- Token refresh/expiration
- TLS validation
- Rate limit handling

## Success Criteria
- [ ] All 83 commands tested
- [ ] Tests run in CI/CD
- [ ] Tests complete in <5 minutes

See GITHUB_PROJECT_PLAN.md for full details."

gh issue create --repo "$REPO" \
  --title "[DOC] Update all documentation for security and performance changes" \
  --milestone "Milestone 4: Testing & Documentation" \
  --label "P1-High,type: documentation,epic: testing,effort: large" \
  --body "**Priority**: P1-High
**Effort**: 12-16 hours

## Documentation Needed
- [ ] Security architecture doc
- [ ] HIPAA compliance guide
- [ ] Configuration reference
- [ ] Migration guide (v1.x → v2.0)
- [ ] Performance tuning guide
- [ ] Troubleshooting guide
- [ ] API reference updates
- [ ] README updates
- [ ] CHANGELOG for v2.0

## Files
- docs/SECURITY.md
- docs/HIPAA_COMPLIANCE.md
- docs/MIGRATION_V2.md (already exists)
- docs/PERFORMANCE.md
- docs/CONFIGURATION.md
- README.md (already updated)
- CHANGELOG.md (already exists)

See GITHUB_PROJECT_PLAN.md for full details."

echo ""
echo "✅ Core v2.0 issues created successfully!"
echo ""
echo "Created issues:"
echo "  Milestone 1 (Critical Security): 5 issues"
echo "  Milestone 3 (Performance): 5 issues"
echo "  Milestone 4 (Testing/Docs): 4 issues"
echo ""
echo "Next steps:"
echo "1. View issues: gh issue list --repo $REPO"
echo "2. View milestones: gh api repos/$REPO/milestones"
echo "3. Set up project board: ./scripts/setup-github-project.sh"
echo "4. Start with issue #8 (Token encryption)"
