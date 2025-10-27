# Implementation Issues for Feature Parity
## Tracking Document for 56 Missing Commands

**Status:** READY TO CREATE
**Total Issues:** 68 (56 command implementations + 12 testing/integration)

---

## Issue Template

Each issue should include:
- **Title:** [PhaseN] Brief description
- **Labels:** phase-N, priority-X, area-Y
- **Assignee:** Developer
- **Milestone:** Phase N Complete
- **Estimate:** Story points/days
- **Description:** Detailed requirements
- **Acceptance Criteria:** Definition of done
- **Dependencies:** Blocking issues

---

## Phase 4: Essential Operations (12 issues)

### Issue #44: [Phase 4] Implement endpoint setup command
**Priority:** Critical
**Area:** endpoint
**Estimate:** 2 days
**Dependencies:** None

**Description:**
Implement `globus-connect-server endpoint setup` command for initial endpoint deployment.

**Tasks:**
- [ ] Add `SetupEndpoint()` method to `pkg/gcs/endpoint.go`
- [ ] Create `internal/commands/endpoint/setup.go`
- [ ] Add comprehensive flag support (name, organization, contact-email, etc.)
- [ ] Implement setup wizard for interactive mode
- [ ] Add validation for required fields
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Update help documentation

**Acceptance Criteria:**
- Command successfully creates new endpoint
- All flags properly validated
- Interactive wizard working
- Config files generated correctly
- Tests passing with >80% coverage
- Zero lint warnings

---

### Issue #45: [Phase 4] Implement endpoint cleanup command
**Priority:** Critical
**Area:** endpoint
**Estimate:** 1 day
**Dependencies:** #44

**Description:**
Implement `globus-connect-server endpoint cleanup` command for permanent endpoint removal.

**Tasks:**
- [ ] Add `CleanupEndpoint()` method to `pkg/gcs/endpoint.go`
- [ ] Create `internal/commands/endpoint/cleanup.go`
- [ ] Add confirmation prompt for destructive operation
- [ ] Add `--force` flag to skip confirmation
- [ ] Handle cleanup of related resources
- [ ] Write unit tests
- [ ] Update help documentation

**Acceptance Criteria:**
- Command removes endpoint successfully
- Confirmation prompt working
- Related resources cleaned up
- Error handling for active transfers
- Tests passing
- Zero lint warnings

---

### Issue #46: [Phase 4] Implement endpoint key-convert command
**Priority:** High
**Area:** endpoint
**Estimate:** 1 day
**Dependencies:** None

**Description:**
Implement `globus-connect-server endpoint key-convert` for deployment key management.

**Tasks:**
- [ ] Add `ConvertDeploymentKey()` method to `pkg/gcs/endpoint.go`
- [ ] Create `internal/commands/endpoint/key_convert.go`
- [ ] Handle key conversion logic
- [ ] Add key validation
- [ ] Write unit tests
- [ ] Update help documentation

**Acceptance Criteria:**
- Key conversion working
- Old keys invalidated
- New keys stored securely
- Tests passing
- Zero lint warnings

---

### Issue #47: [Phase 4] Implement node setup command
**Priority:** Critical
**Area:** node
**Estimate:** 2 days
**Dependencies:** None

**Description:**
Implement `globus-connect-server node setup` for node configuration and initialization.

**Tasks:**
- [ ] Add `SetupNode()` method to `pkg/gcs/node.go`
- [ ] Create `internal/commands/node/setup.go`
- [ ] Add flags for node configuration
- [ ] Implement setup wizard
- [ ] Add node secret generation
- [ ] Write unit tests
- [ ] Update help documentation

**Acceptance Criteria:**
- Node setup working
- Configuration files created
- Secrets generated securely
- Tests passing with >80% coverage
- Zero lint warnings

---

### Issue #48: [Phase 4] Implement node cleanup command
**Priority:** High
**Area:** node
**Estimate:** 1 day
**Dependencies:** #47

**Description:**
Implement `globus-connect-server node cleanup` for node removal.

**Tasks:**
- [ ] Add `CleanupNode()` method to `pkg/gcs/node.go`
- [ ] Create `internal/commands/node/cleanup.go`
- [ ] Add confirmation prompt
- [ ] Handle cleanup of node state
- [ ] Write unit tests
- [ ] Update help documentation

**Acceptance Criteria:**
- Node cleanup working
- Confirmation prompt implemented
- State cleaned up properly
- Tests passing
- Zero lint warnings

---

### Issue #49: [Phase 4] Implement node enable/disable commands
**Priority:** Medium
**Area:** node
**Estimate:** 1 day
**Dependencies:** None

**Description:**
Implement `globus-connect-server node enable` and `node disable` for node activation control.

**Tasks:**
- [ ] Add `EnableNode()` method to `pkg/gcs/node.go`
- [ ] Add `DisableNode()` method to `pkg/gcs/node.go`
- [ ] Create `internal/commands/node/enable.go`
- [ ] Create `internal/commands/node/disable.go`
- [ ] Add status verification
- [ ] Write unit tests for both commands
- [ ] Update help documentation

**Acceptance Criteria:**
- Enable/disable working
- Node status changes reflected
- Tests passing
- Zero lint warnings

---

### Issue #50: [Phase 4] Implement node new-secret command
**Priority:** Medium
**Area:** node
**Estimate:** 1 day
**Dependencies:** None

**Description:**
Implement `globus-connect-server node new-secret` for generating new node authentication secrets.

**Tasks:**
- [ ] Add `GenerateNodeSecret()` method to `pkg/gcs/node.go`
- [ ] Create `internal/commands/node/new_secret.go`
- [ ] Implement secure secret generation
- [ ] Add secret rotation logic
- [ ] Write unit tests
- [ ] Update help documentation

**Acceptance Criteria:**
- Secret generation working
- Old secrets invalidated
- New secrets stored securely
- Tests passing
- Zero lint warnings

---

### Issue #51: [Phase 4] Implement collection check command
**Priority:** High
**Area:** collection
**Estimate:** 2 days
**Dependencies:** None

**Description:**
Implement `globus-connect-server collection check` for collection configuration validation.

**Tasks:**
- [ ] Add `CheckCollection()` method to `pkg/gcs/collection.go`
- [ ] Create `internal/commands/collection/check.go`
- [ ] Implement validation rules
- [ ] Add detailed validation reporting
- [ ] Handle pre-deployment checks
- [ ] Write unit tests
- [ ] Update help documentation

**Acceptance Criteria:**
- Validation working correctly
- Clear error messages for issues
- All validation rules implemented
- Tests passing with >80% coverage
- Zero lint warnings

---

### Issue #52: [Phase 4] Implement collection batch-delete command
**Priority:** Medium
**Area:** collection
**Estimate:** 1 day
**Dependencies:** None

**Description:**
Implement `globus-connect-server collection batch-delete` for removing multiple guest collections.

**Tasks:**
- [ ] Add `BatchDeleteCollections()` method to `pkg/gcs/collection.go`
- [ ] Create `internal/commands/collection/batch_delete.go`
- [ ] Add collection ID list handling
- [ ] Implement transaction/rollback logic
- [ ] Add progress reporting
- [ ] Write unit tests
- [ ] Update help documentation

**Acceptance Criteria:**
- Batch delete working
- Error handling for partial failures
- Progress reporting implemented
- Tests passing
- Zero lint warnings

---

### Issue #53: [Phase 4] Integration testing for essential operations
**Priority:** High
**Area:** testing
**Estimate:** 2 days
**Dependencies:** #44-#52

**Description:**
Create integration tests for all Phase 4 commands against staging GCS endpoint.

**Tasks:**
- [ ] Set up test GCS endpoint
- [ ] Write end-to-end test suite
- [ ] Test full endpoint lifecycle (setup → use → cleanup)
- [ ] Test full node lifecycle (setup → enable → disable → cleanup)
- [ ] Test collection validation
- [ ] Test batch operations
- [ ] Document test environment setup

**Acceptance Criteria:**
- All integration tests passing
- Test coverage >80%
- Documentation complete
- CI pipeline updated

---

### Issue #54: [Phase 4] Update main.go and wire Phase 4 commands
**Priority:** High
**Area:** cli
**Estimate:** 0.5 days
**Dependencies:** #44-#52

**Description:**
Wire all Phase 4 commands into the CLI hierarchy and update main.go.

**Tasks:**
- [ ] Update `internal/commands/endpoint/endpoint.go` to add new subcommands
- [ ] Update `internal/commands/node/node.go` to add new subcommands
- [ ] Update `internal/commands/collection/collection.go` to add new subcommands
- [ ] Verify command hierarchy
- [ ] Update help text
- [ ] Test all commands accessible

**Acceptance Criteria:**
- All commands accessible via CLI
- Help text complete and accurate
- Command hierarchy logical
- Tests passing

---

### Issue #55: [Phase 4] Documentation for essential operations
**Priority:** Medium
**Area:** documentation
**Estimate:** 1 day
**Dependencies:** #44-#54

**Description:**
Document all Phase 4 commands with examples and use cases.

**Tasks:**
- [ ] Add command reference entries
- [ ] Write usage examples
- [ ] Document setup workflows
- [ ] Add troubleshooting guide
- [ ] Update README if needed

**Acceptance Criteria:**
- All commands documented
- Examples tested and working
- Troubleshooting guide complete
- Documentation reviewed

---

## Phase 5: Ownership & Identity (12 issues)

### Issue #56-65: [Phase 5] Individual ownership commands
*Similar structure to Phase 4 issues*

Commands:
- #56: endpoint set-owner
- #57: endpoint set-owner-string
- #58: endpoint reset-owner-string
- #59: endpoint set-subscription-id
- #60: collection set-owner
- #61: collection set-owner-string
- #62: collection reset-owner-string
- #63: collection set-subscription-admin-verified
- #64: collection role
- #65: endpoint role

### Issue #66: [Phase 5] Integration testing for ownership
### Issue #67: [Phase 5] Wire Phase 5 commands

---

## Phase 6: Domain Management (4 issues)

### Issue #68-69: [Phase 6] Domain commands
- #68: endpoint domain (with subcommands)
- #69: collection domain (with subcommands)

### Issue #70: [Phase 6] Domain integration testing
### Issue #71: [Phase 6] Wire Phase 6 commands

---

## Phase 7: Authentication & Security (12 issues)

### Issue #72-76: [Phase 7] Auth Policy commands
- #72: auth-policy create
- #73: auth-policy list
- #74: auth-policy show
- #75: auth-policy update
- #76: auth-policy delete

### Issue #77-81: [Phase 7] OIDC commands
- #77: oidc create
- #78: oidc register
- #79: oidc show
- #80: oidc update
- #81: oidc delete

### Issue #82: [Phase 7] Integration testing for auth & security
### Issue #83: [Phase 7] Wire Phase 7 commands

---

## Phase 8: Session Management (5 issues)

### Issue #84-86: [Phase 8] Session commands
- #84: session show
- #85: session update
- #86: session consent

### Issue #87: [Phase 8] Integration testing for sessions
### Issue #88: [Phase 8] Wire Phase 8 commands

---

## Phase 9: Sharing Policies (6 issues)

### Issue #89-92: [Phase 9] Sharing Policy commands
- #89: sharing-policy create
- #90: sharing-policy list
- #91: sharing-policy show
- #92: sharing-policy delete

### Issue #93: [Phase 9] Integration testing for sharing policies
### Issue #94: [Phase 9] Wire Phase 9 commands

---

## Phase 10: User Credentials (12 issues)

### Issue #95-103: [Phase 10] User Credential commands
- #95: user-credential activescale-create
- #96: user-credential oauth-create
- #97: user-credential s3-create
- #98: user-credential s3-keys-add
- #99: user-credential s3-keys-update
- #100: user-credential s3-keys-delete
- #101: user-credential list
- #102: user-credential show
- #103: user-credential delete

### Issue #104: [Phase 10] Integration testing for user credentials
### Issue #105: [Phase 10] Wire Phase 10 commands
### Issue #106: [Phase 10] Security audit for credential handling

---

## Phase 11: Audit & Compliance (6 issues)

### Issue #107-109: [Phase 11] Audit commands
- #107: audit load
- #108: audit query
- #109: audit dump

### Issue #110: [Phase 11] Audit database implementation (SQLite)
### Issue #111: [Phase 11] Integration testing for audit
### Issue #112: [Phase 11] Wire Phase 11 commands

---

## Phase 12: Advanced Operations (3 issues)

### Issue #113: [Phase 12] Implement endpoint upgrade command
### Issue #114: [Phase 12] Integration testing for upgrade
### Issue #115: [Phase 12] Wire Phase 12 commands

---

## Summary by Phase

| Phase | Command Issues | Testing Issues | Integration Issues | Total |
|-------|---------------|----------------|-------------------|-------|
| Phase 4 | 9 | 1 | 2 | 12 |
| Phase 5 | 10 | 1 | 1 | 12 |
| Phase 6 | 2 | 1 | 1 | 4 |
| Phase 7 | 10 | 1 | 1 | 12 |
| Phase 8 | 3 | 1 | 1 | 5 |
| Phase 9 | 4 | 1 | 1 | 6 |
| Phase 10 | 9 | 2 | 1 | 12 |
| Phase 11 | 3 | 1 | 2 | 6 |
| Phase 12 | 1 | 1 | 1 | 3 |
| **TOTAL** | **51** | **10** | **11** | **72** |

---

## Issue Creation Checklist

For each issue, ensure:
- [ ] Clear, descriptive title
- [ ] Detailed description with context
- [ ] Specific tasks listed
- [ ] Clear acceptance criteria
- [ ] Dependencies identified
- [ ] Appropriate labels applied
- [ ] Milestone assigned
- [ ] Estimate provided
- [ ] Priority set

---

## Next Actions

1. **Create GitHub Issues** - Use this document to create issues #44-#115
2. **Link Issues** - Set up dependency relationships
3. **Create Project Board Views** - Organize by phase
4. **Assign Issues** - Distribute to team members
5. **Start Phase 4** - Begin implementation

---

**Document Status:** READY FOR ISSUE CREATION
**Total Work Items:** 72 issues
**Estimated Timeline:** 12-15 weeks
**Start Date:** Upon issue creation
