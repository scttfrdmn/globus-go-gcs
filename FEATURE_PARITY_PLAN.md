# Feature Parity Implementation Plan
## Completing the Globus Connect Server Go CLI

**Document Version:** 1.0
**Created:** 2025-10-26
**Status:** READY FOR IMPLEMENTATION
**Current Progress:** 27/83 commands (33%) - All core CRUD complete

---

## Executive Summary

This plan details the implementation of 56 remaining commands to achieve 100% feature parity with the Python globus-connect-server CLI.

**Timeline:** 12-15 weeks for complete implementation
**Approach:** 9 sequential phases, each building on the previous
**Quality Target:** Maintain A+ Go Report Card, 0 lint warnings, 80%+ test coverage

---

## Current Implementation Status

### ✅ Completed Commands (27/83)

**Phase 1-3: Core CRUD** - ✅ COMPLETE
- Authentication: login, logout, whoami
- Collection: create, list, show, update, delete
- Endpoint: show, update
- Node: create, list, show, update, delete
- Role: create, list, show, delete
- Storage Gateway: create, list, show, update, delete

---

## Missing Commands Analysis (56 commands)

### New Command Groups (29 commands)
1. **Audit** (3) - Log management and compliance
2. **Auth Policy** (5) - Authentication policy CRUD
3. **OIDC** (5) - OpenID Connect integration
4. **Session** (3) - CLI session management
5. **Sharing Policy** (4) - Sharing policy CRUD
6. **User Credentials** (9) - User credential management

### Extensions to Existing Groups (27 commands)
1. **Collection** (8) - Batch operations, ownership, domain, validation
2. **Endpoint** (10) - Setup, cleanup, ownership, domain, roles, upgrade
3. **Node** (5) - Setup, cleanup, enable/disable, secrets

---

## Implementation Phases

### Phase 4: Essential Operations (10 commands)
**Timeline:** 2 weeks
**Priority:** CRITICAL
**Dependencies:** None

#### Commands to Implement

**Endpoint (3 commands)**
- `endpoint setup` - Initial endpoint deployment
- `endpoint cleanup` - Permanent endpoint removal
- `endpoint key-convert` - Deployment key management

**Node (5 commands)**
- `node setup` - Node configuration and initialization
- `node cleanup` - Node removal
- `node enable` - Activate node for transfers
- `node disable` - Deactivate node
- `node new-secret` - Generate new node authentication secret

**Collection (2 commands)**
- `collection check` - Validate collection configuration
- `collection batch-delete` - Remove multiple guest collections

#### Implementation Details

**Files to Create:**
```
internal/commands/endpoint/
├── setup.go           # endpoint setup command
├── setup_test.go
├── cleanup.go         # endpoint cleanup command
├── cleanup_test.go
├── key_convert.go     # key-convert command
└── key_convert_test.go

internal/commands/node/
├── setup.go           # node setup command
├── setup_test.go
├── cleanup.go         # node cleanup command
├── cleanup_test.go
├── enable.go          # node enable command
├── enable_test.go
├── disable.go         # node disable command
├── disable_test.go
├── new_secret.go      # new-secret command
└── new_secret_test.go

internal/commands/collection/
├── check.go           # collection check command
├── check_test.go
├── batch_delete.go    # batch-delete command
└── batch_delete_test.go
```

**API Client Extensions:**
```go
// pkg/gcs/endpoint.go
func (c *Client) SetupEndpoint(ctx context.Context, config *EndpointSetupConfig) (*Endpoint, error)
func (c *Client) CleanupEndpoint(ctx context.Context) error
func (c *Client) ConvertDeploymentKey(ctx context.Context, oldKey string) (*DeploymentKeyResult, error)

// pkg/gcs/node.go
func (c *Client) SetupNode(ctx context.Context, config *NodeSetupConfig) (*Node, error)
func (c *Client) CleanupNode(ctx context.Context, nodeID string) error
func (c *Client) EnableNode(ctx context.Context, nodeID string) error
func (c *Client) DisableNode(ctx context.Context, nodeID string) error
func (c *Client) GenerateNodeSecret(ctx context.Context, nodeID string) (*NodeSecret, error)

// pkg/gcs/collection.go
func (c *Client) CheckCollection(ctx context.Context, collectionID string) (*CollectionValidation, error)
func (c *Client) BatchDeleteCollections(ctx context.Context, collectionIDs []string) (*BatchDeleteResult, error)
```

**Testing Requirements:**
- Unit tests for all commands
- Mock API responses
- Error handling validation
- Flag validation tests

---

### Phase 5: Ownership & Identity (10 commands)
**Timeline:** 1.5 weeks
**Priority:** HIGH
**Dependencies:** Phase 4 complete

#### Commands to Implement

**Endpoint Ownership (4 commands)**
- `endpoint set-owner` - Assign endpoint owner role
- `endpoint set-owner-string` - Custom owner display name
- `endpoint reset-owner-string` - Reset to default (ClientID)
- `endpoint set-subscription-id` - Update subscription assignment

**Collection Ownership (5 commands)**
- `collection set-owner` - Designate collection owner
- `collection set-owner-string` - Custom owner display
- `collection reset-owner-string` - Reset owner string
- `collection set-subscription-admin-verified` - Set admin verification
- `collection role` - Collection-specific role management

**Endpoint Roles (1 command)**
- `endpoint role` - Endpoint-specific role management

#### Implementation Details

**Files to Create:**
```
internal/commands/endpoint/
├── set_owner.go
├── set_owner_string.go
├── reset_owner_string.go
├── set_subscription_id.go
└── role.go           # Endpoint role subcommands

internal/commands/collection/
├── set_owner.go
├── set_owner_string.go
├── reset_owner_string.go
├── set_subscription_admin_verified.go
└── role.go           # Collection role subcommands
```

**API Client Extensions:**
```go
// pkg/gcs/endpoint.go
func (c *Client) SetEndpointOwner(ctx context.Context, principalURN string) error
func (c *Client) SetEndpointOwnerString(ctx context.Context, ownerString string) error
func (c *Client) ResetEndpointOwnerString(ctx context.Context) error
func (c *Client) SetSubscriptionID(ctx context.Context, subscriptionID string) error

// pkg/gcs/collection.go
func (c *Client) SetCollectionOwner(ctx context.Context, collectionID, principalURN string) error
func (c *Client) SetCollectionOwnerString(ctx context.Context, collectionID, ownerString string) error
func (c *Client) ResetCollectionOwnerString(ctx context.Context, collectionID string) error
func (c *Client) SetSubscriptionAdminVerified(ctx context.Context, collectionID string, verified bool) error
```

---

### Phase 6: Domain Management (2 commands)
**Timeline:** 1 week
**Priority:** HIGH
**Dependencies:** Phase 5 complete

#### Commands to Implement

- `endpoint domain` - Manage endpoint custom domain
- `collection domain` - Manage collection custom domain

#### Implementation Details

**Files to Create:**
```
internal/commands/endpoint/
└── domain.go         # Domain subcommands (setup, show, delete)

internal/commands/collection/
└── domain.go         # Domain subcommands (setup, show, delete)
```

**API Client Extensions:**
```go
// pkg/gcs/domain.go
type DomainConfig struct {
    Domain       string
    Certificate  string
    PrivateKey   string
    Verified     bool
}

func (c *Client) SetupEndpointDomain(ctx context.Context, config *DomainConfig) error
func (c *Client) GetEndpointDomain(ctx context.Context) (*DomainConfig, error)
func (c *Client) DeleteEndpointDomain(ctx context.Context) error

func (c *Client) SetupCollectionDomain(ctx context.Context, collectionID string, config *DomainConfig) error
func (c *Client) GetCollectionDomain(ctx context.Context, collectionID string) (*DomainConfig, error)
func (c *Client) DeleteCollectionDomain(ctx context.Context, collectionID string) error
```

---

### Phase 7: Authentication & Security (10 commands)
**Timeline:** 2 weeks
**Priority:** HIGH
**Dependencies:** None (can run parallel to Phases 4-6)

#### Commands to Implement

**Auth Policy (5 commands)**
- `auth-policy create` - Create authentication policy
- `auth-policy list` - List all auth policies
- `auth-policy show` - Display policy details
- `auth-policy update` - Modify authentication policy
- `auth-policy delete` - Remove authentication policy

**OIDC (5 commands)**
- `oidc create` - Create OIDC server configuration
- `oidc register` - Register existing OIDC server
- `oidc show` - Display OIDC configuration
- `oidc update` - Modify OIDC settings
- `oidc delete` - Remove OIDC server

#### Implementation Details

**Files to Create:**
```
internal/commands/
├── authpolicy/
│   ├── authpolicy.go
│   ├── create.go
│   ├── create_test.go
│   ├── list.go
│   ├── list_test.go
│   ├── show.go
│   ├── show_test.go
│   ├── update.go
│   ├── update_test.go
│   ├── delete.go
│   └── delete_test.go
└── oidc/
    ├── oidc.go
    ├── create.go
    ├── create_test.go
    ├── register.go
    ├── register_test.go
    ├── show.go
    ├── show_test.go
    ├── update.go
    ├── update_test.go
    ├── delete.go
    └── delete_test.go

pkg/gcs/
├── authpolicy.go      # Auth policy API client
├── authpolicy_test.go
├── oidc.go            # OIDC API client
└── oidc_test.go
```

**API Client Types:**
```go
// pkg/gcs/authpolicy.go
type AuthPolicy struct {
    ID                    string
    Name                  string
    Description           string
    RequireMFA            bool
    RequireHighAssurance  bool
    AllowedDomains        []string
    BlockedDomains        []string
}

// pkg/gcs/oidc.go
type OIDCServer struct {
    ID           string
    Issuer       string
    ClientID     string
    ClientSecret string
    Audience     string
    Scopes       []string
}
```

---

### Phase 8: Session Management (3 commands)
**Timeline:** 1 week
**Priority:** MEDIUM
**Dependencies:** Phase 7 complete

#### Commands to Implement

- `session show` - Display current authentication session
- `session update` - Modify session settings
- `session consent` - Update session consents

#### Implementation Details

**Files to Create:**
```
internal/commands/session/
├── session.go
├── show.go
├── show_test.go
├── update.go
├── update_test.go
├── consent.go
└── consent_test.go
```

---

### Phase 9: Sharing Policies (4 commands)
**Timeline:** 1 week
**Priority:** MEDIUM
**Dependencies:** None

#### Commands to Implement

- `sharing-policy create` - Create sharing policy
- `sharing-policy list` - List all sharing policies
- `sharing-policy show` - Display policy details
- `sharing-policy delete` - Remove sharing policy

#### Implementation Details

**Files to Create:**
```
internal/commands/sharingpolicy/
├── sharingpolicy.go
├── create.go
├── create_test.go
├── list.go
├── list_test.go
├── show.go
├── show_test.go
├── delete.go
└── delete_test.go

pkg/gcs/
├── sharingpolicy.go
└── sharingpolicy_test.go
```

---

### Phase 10: User Credentials (9 commands)
**Timeline:** 2 weeks
**Priority:** MEDIUM
**Dependencies:** Phase 7 complete

#### Commands to Implement

- `user-credential activescale-create` - ActiveScale credentials
- `user-credential oauth-create` - OAuth2 credentials
- `user-credential s3-create` - S3 credentials
- `user-credential s3-keys-add` - Add S3 IAM keys
- `user-credential s3-keys-update` - Update S3 keys
- `user-credential s3-keys-delete` - Remove S3 keys
- `user-credential list` - List all credentials
- `user-credential show` - Display credential details
- `user-credential delete` - Remove credential

#### Implementation Details

**Files to Create:**
```
internal/commands/usercredential/
├── usercredential.go
├── activescale_create.go
├── oauth_create.go
├── s3_create.go
├── s3_keys_add.go
├── s3_keys_update.go
├── s3_keys_delete.go
├── list.go
├── show.go
├── delete.go
└── *_test.go files

pkg/gcs/
├── usercredential.go
└── usercredential_test.go
```

**API Client Types:**
```go
type UserCredential struct {
    ID           string
    IdentityID   string
    StorageGatewayID string
    Type         string  // "activescale", "oauth", "s3"
    // Type-specific fields
    S3Keys       []S3Key
    OAuthToken   string
}

type S3Key struct {
    AccessKeyID     string
    SecretAccessKey string
    CreatedAt       time.Time
}
```

---

### Phase 11: Audit & Compliance (3 commands)
**Timeline:** 2 weeks
**Priority:** LOW
**Dependencies:** None

#### Commands to Implement

- `audit load` - Load audit logs into search database
- `audit query` - Search audit logs with filters
- `audit dump` - Export audit logs to file

#### Implementation Details

**Files to Create:**
```
internal/commands/audit/
├── audit.go
├── load.go
├── load_test.go
├── query.go
├── query_test.go
├── dump.go
└── dump_test.go

pkg/gcs/
├── audit.go
└── audit_test.go
```

**Features:**
- SQLite backend for local audit database
- Query language support (filters, date ranges)
- Export formats (JSON, CSV)

---

### Phase 12: Advanced Operations (1 command)
**Timeline:** 1 week
**Priority:** LOW
**Dependencies:** All phases complete

#### Commands to Implement

- `endpoint upgrade` - Upgrade endpoint to latest version

#### Implementation Details

**Files to Create:**
```
internal/commands/endpoint/
├── upgrade.go
└── upgrade_test.go
```

**Features:**
- Version compatibility checks
- Pre-upgrade validation
- Post-upgrade verification
- Rollback support

---

## Dependency Graph

```
START
  │
  ├─→ Phase 4 (Essential Operations)
  │     ├─→ Phase 5 (Ownership & Identity)
  │     │     └─→ Phase 6 (Domain Management)
  │     │
  │     └─→ Phase 7 (Auth & Security)
  │           └─→ Phase 8 (Session Management)
  │                 └─→ Phase 10 (User Credentials)
  │
  ├─→ Phase 9 (Sharing Policy) - Independent
  │
  ├─→ Phase 11 (Audit) - Independent
  │
  └─→ Phase 12 (Upgrade) - Requires all phases
```

---

## Implementation Standards

### For Each Command:

1. **API Client Method**
   - Add to `pkg/gcs/*.go`
   - Include full type definitions
   - Comprehensive error handling
   - Unit tests with httptest

2. **CLI Command**
   - Add to `internal/commands/*/`
   - Follow existing patterns
   - Cobra command structure
   - Flag definitions and validation

3. **Tests**
   - Command structure tests
   - Flag validation tests
   - Mock execution tests
   - Error case coverage

4. **Integration**
   - Wire into command hierarchy
   - Update help text
   - Add to main.go if new group

5. **Quality Checks**
   - All tests passing
   - Zero lint warnings
   - Cyclomatic complexity < 30
   - Test coverage > 80%

---

## Timeline & Estimates

| Phase | Commands | Weeks | Start | End | Cumulative |
|-------|----------|-------|-------|-----|------------|
| **Phase 4** | 10 | 2 | Week 1 | Week 2 | 2 weeks |
| **Phase 5** | 10 | 1.5 | Week 3 | Week 4 | 4 weeks |
| **Phase 6** | 2 | 1 | Week 5 | Week 5 | 5 weeks |
| **Phase 7** | 10 | 2 | Week 1* | Week 2* | 7 weeks |
| **Phase 8** | 3 | 1 | Week 6 | Week 6 | 8 weeks |
| **Phase 9** | 4 | 1 | Week 3* | Week 3* | 9 weeks |
| **Phase 10** | 9 | 2 | Week 7 | Week 8 | 10 weeks |
| **Phase 11** | 3 | 2 | Week 5* | Week 6* | 12 weeks |
| **Phase 12** | 1 | 1 | Week 9 | Week 9 | 13 weeks |
| **Buffer** | - | 1.5 | Week 10 | Week 11 | **15 weeks** |

*Phases 7, 9, 11 can run in parallel with sequential phases

---

## Success Metrics

### Technical Metrics:
- ✅ 83/83 commands implemented (100%)
- ✅ All tests passing (>80% coverage)
- ✅ Zero lint warnings
- ✅ A+ Go Report Card grade
- ✅ Cyclomatic complexity < 30 for all functions

### Functional Metrics:
- ✅ 100% feature parity with Python CLI
- ✅ All documented use cases supported
- ✅ Config file compatibility maintained
- ✅ Performance ≥ Python CLI

---

## Risk Assessment

### High Risk Areas:
1. **User Credentials** - Multiple connector types, security-sensitive
2. **Audit System** - Database integration, query language
3. **OIDC** - External auth integration complexity
4. **Endpoint Setup** - Critical deployment path

### Mitigation Strategies:
- Start with high-risk items early
- Extensive testing for security-sensitive code
- Reference Python implementation closely
- Staging environment testing before production

---

## Resource Requirements

### Development:
- 1 senior Go developer (full-time)
- GCS test endpoint access
- Python CLI for reference

### Testing:
- GCS v5.4+ endpoint
- Multiple storage connector types
- OIDC test provider
- Audit log samples

---

## Next Steps

1. **Review & Approve** this plan
2. **Phase 4 Kickoff** - Begin essential operations
3. **Daily Progress** - Update todos, track issues
4. **Weekly Review** - Phase completion, blockers
5. **Phase Completion** - Commit, test, document

---

**Document Status:** READY FOR IMPLEMENTATION
**Approval Required:** Yes
**Start Date:** Upon approval
**Target Completion:** 12-15 weeks from start
