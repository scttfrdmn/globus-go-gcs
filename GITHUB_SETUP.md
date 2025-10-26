# GitHub Repository Setup Guide

## ✅ Phase 0 Complete!

All local infrastructure is ready. This guide will help you push to GitHub and set up the project board.

## Step 1: Create GitHub Repository

### Option A: Using GitHub CLI (Recommended)

```bash
# Install gh CLI if needed
brew install gh  # macOS
# or download from https://cli.github.com

# Authenticate
gh auth login

# Create repository
gh repo create globus-go-gcs \
  --public \
  --description "Complete Go port of Globus Connect Server v5 CLI - drop-in replacement for Python globus-connect-server" \
  --homepage "https://github.com/scttfrdmn/globus-go-gcs"

# Push code
git remote add origin https://github.com/scttfrdmn/globus-go-gcs.git
git push -u origin main
```

### Option B: Using GitHub Web UI

1. Go to https://github.com/new
2. Repository name: `globus-go-gcs`
3. Description: "Complete Go port of Globus Connect Server v5 CLI - drop-in replacement for Python globus-connect-server"
4. Public repository
5. Do NOT initialize with README, .gitignore, or license (we already have these)
6. Click "Create repository"
7. Push code:

```bash
git remote add origin https://github.com/scttfrdmn/globus-go-gcs.git
git branch -M main
git push -u origin main
```

## Step 2: Sync Labels

After pushing, sync the 147 labels:

```bash
# This will create all labels from .github/labels.yml
gh label sync --repo scttfrdmn/globus-go-gcs

# Alternatively, the GitHub Action will sync them automatically on next push
```

## Step 3: Create GitHub Project Board

### Create the Project

```bash
# Create project (note the URL that's returned - you'll need the project number)
gh project create \
  --owner scttfrdmn \
  --title "Globus Connect Server Go Port"

# The output will show something like:
# https://github.com/users/scttfrdmn/projects/12
# Note the number (12 in this example)
```

### Add Custom Fields

Replace `PROJECT_NUMBER` with the actual project number from above:

```bash
PROJECT_NUMBER=12  # <-- REPLACE THIS

# Add Persona field
gh project field-create $PROJECT_NUMBER \
  --owner scttfrdmn \
  --name "Persona" \
  --data-type "SINGLE_SELECT" \
  --single-select-options "System Administrator,Data Manager,Research PI,IT Manager"

# Add Phase field
gh project field-create $PROJECT_NUMBER \
  --owner scttfrdmn \
  --name "Phase" \
  --data-type "SINGLE_SELECT" \
  --single-select-options "Phase 0 - Setup,Phase 1 - POC,Phase 2 - Core,Phase 3 - Collections,Phase 4 - Advanced,Phase 5 - Release,Backlog"

# Add Estimate field (number)
gh project field-create $PROJECT_NUMBER \
  --owner scttfrdmn \
  --name "Estimate" \
  --data-type "NUMBER"

# Add Priority field
gh project field-create $PROJECT_NUMBER \
  --owner scttfrdmn \
  --name "Priority" \
  --data-type "SINGLE_SELECT" \
  --single-select-options "Critical,High,Medium,Low"
```

### Create Project Views (Manual - Web UI)

Go to your project board and create these views:

1. **Kanban** (Default)
   - Layout: Board
   - Group by: Status
   - Columns: Triage, Todo, In Progress, In Review, Done

2. **By Phase**
   - Click "+ New view"
   - Layout: Board
   - Group by: Phase
   - Sort by: Priority

3. **By Persona**
   - Click "+ New view"
   - Layout: Board
   - Group by: Persona
   - Sort by: Phase, then Priority

4. **Current Sprint**
   - Click "+ New view"
   - Layout: Table
   - Filter: Status ≠ Done AND Phase = "Phase 1 - POC"
   - Show columns: Title, Status, Priority, Persona, Estimate, Assignees

5. **Backlog**
   - Click "+ New view"
   - Layout: Table
   - Filter: Status = Todo OR Status = Triage
   - Sort by: Priority, then Phase

## Step 4: Create Milestones

```bash
# Phase 0 (Complete)
gh milestone create "[Phase 0] Project Infrastructure Complete" \
  --repo scttfrdmn/globus-go-gcs \
  --due-date 2025-11-01 \
  --description "✅ COMPLETE: Project setup, GitHub infrastructure, planning"

# Phase 1
gh milestone create "[Phase 1] Proof of Concept - Authentication Working" \
  --repo scttfrdmn/globus-go-gcs \
  --due-date 2025-12-15 \
  --description "OAuth2 auth working, token storage, basic endpoint commands, JSON output"

# Phase 2
gh milestone create "[Phase 2] Endpoint & Node Management Complete" \
  --repo scttfrdmn/globus-go-gcs \
  --due-date 2026-02-15 \
  --description "Complete endpoint lifecycle, node management, RBAC"

# Phase 3
gh milestone create "[Phase 3] Collection & Storage Management Complete" \
  --repo scttfrdmn/globus-go-gcs \
  --due-date 2026-04-15 \
  --description "Collections (mapped/guest), storage gateways, all storage connectors"

# Phase 4
gh milestone create "[Phase 4] Advanced Features Complete" \
  --repo scttfrdmn/globus-go-gcs \
  --due-date 2026-06-15 \
  --description "OIDC, policies, user credentials, audit logs, diagnostics"

# Phase 5
gh milestone create "[Phase 5] v1.0.0 Release" \
  --repo scttfrdmn/globus-go-gcs \
  --due-date 2026-07-15 \
  --description "Polish, documentation, multi-platform binaries, v1.0.0 release"
```

## Step 5: Create Initial Issues

### Phase 0 Issues (Documentation)

```bash
# Issue #1 - Already done but document it
gh issue create \
  --repo scttfrdmn/globus-go-gcs \
  --title "[Phase 0] Initialize Go module and project structure" \
  --label "phase: 0-setup,priority: high,area: build" \
  --milestone "[Phase 0] Project Infrastructure Complete" \
  --body "**Status**: ✅ COMPLETE

Created:
- Go module (github.com/scttfrdmn/globus-go-gcs)
- Directory structure (cmd/, pkg/, internal/)
- Basic CLI entry point
- Makefile for common tasks

**Persona**: System Administrator"

# Mark it as closed immediately
gh issue close 1 --repo scttfrdmn/globus-go-gcs --reason completed

# Issue #2 - CI/CD (complete)
gh issue create \
  --repo scttfrdmn/globus-go-gcs \
  --title "[Phase 0] Set up CI/CD pipeline" \
  --label "phase: 0-setup,priority: high,area: build" \
  --milestone "[Phase 0] Project Infrastructure Complete" \
  --body "**Status**: ✅ COMPLETE

Created:
- .github/workflows/ci.yml (test, lint, build)
- .github/workflows/labels.yml (label sync)
- Multi-platform testing (Ubuntu, macOS)
- Go versions: 1.21, 1.22, 1.23

**Persona**: System Administrator"

gh issue close 2 --repo scttfrdmn/globus-go-gcs --reason completed

# Issue #3 - GitHub infrastructure (complete)
gh issue create \
  --repo scttfrdmn/globus-go-gcs \
  --title "[Phase 0] Configure GitHub infrastructure" \
  --label "phase: 0-setup,priority: high,area: docs" \
  --milestone "[Phase 0] Project Infrastructure Complete" \
  --body "**Status**: ✅ COMPLETE

Created:
- 147 labels (.github/labels.yml)
- Issue templates (feature request, bug report)
- PR template with persona impact
- Project board with custom fields

**Persona**: All"

gh issue close 3 --repo scttfrdmn/globus-go-gcs --reason completed
```

### Phase 1 Issues (Next Up)

```bash
# Issue #4 - OAuth2 authentication
gh issue create \
  --repo scttfrdmn/globus-go-gcs \
  --title "[Phase 1] Implement OAuth2 authentication flow" \
  --label "phase: 1-poc,priority: critical,area: auth,persona: system-admin,good first issue" \
  --milestone "[Phase 1] Proof of Concept - Authentication Working" \
  --body "**Goal**: Implement OAuth2 authentication (login/logout/whoami)

**Tasks**:
- [ ] Create auth package (pkg/auth/)
- [ ] Implement \`login\` command (browser-based OAuth2)
- [ ] Implement \`logout\` command (clear tokens)
- [ ] Implement \`whoami\` command (show user info)
- [ ] Use globus-go-sdk v3 Auth client
- [ ] Add tests

**Success Criteria**:
- Can authenticate via browser OAuth2 flow
- Tokens stored in \`~/.globus-connect-server/tokens.json\`
- Can query user information

**Estimate**: 1 week
**Persona**: System Administrator"

# Issue #5 - Token storage
gh issue create \
  --repo scttfrdmn/globus-go-gcs \
  --title "[Phase 1] Implement token storage compatible with Python CLI" \
  --label "phase: 1-poc,priority: critical,area: config,persona: system-admin" \
  --milestone "[Phase 1] Proof of Concept - Authentication Working" \
  --body "**Goal**: Token storage compatible with Python GCS CLI

**Tasks**:
- [ ] Create config package (pkg/config/)
- [ ] Read/write \`~/.globus-connect-server/tokens.json\`
- [ ] Parse Python CLI token format
- [ ] Implement token refresh logic
- [ ] Add tests

**Success Criteria**:
- Can read tokens from Python CLI installation
- Can write tokens that Python CLI can read
- Token refresh works automatically

**Estimate**: 3 days
**Persona**: System Administrator"

# Issue #6 - Config management
gh issue create \
  --repo scttfrdmn/globus-go-gcs \
  --title "[Phase 1] Implement configuration file management" \
  --label "phase: 1-poc,priority: high,area: config,persona: system-admin" \
  --milestone "[Phase 1] Proof of Concept - Authentication Working" \
  --body "**Goal**: Configuration management compatible with Python CLI

**Tasks**:
- [ ] Read/write \`~/.globus-connect-server/config.json\`
- [ ] Parse Python CLI config format
- [ ] Support all config options
- [ ] Add tests

**Success Criteria**:
- Can read Python CLI config files
- Can write compatible config files
- Configuration values accessible to commands

**Estimate**: 2 days
**Persona**: System Administrator"

# Issue #7 - GCS Manager API client
gh issue create \
  --repo scttfrdmn/globus-go-gcs \
  --title "[Phase 1] Create GCS Manager API base client" \
  --label "phase: 1-poc,priority: critical,area: gcs-client" \
  --milestone "[Phase 1] Proof of Concept - Authentication Working" \
  --body "**Goal**: Implement GCS Manager API client

**Tasks**:
- [ ] Create pkg/client/gcs/ package
- [ ] Implement base client (connects by FQDN)
- [ ] Add GetGCSInfo() method
- [ ] Add GetEndpoint() method
- [ ] Error handling with globus-go-sdk patterns
- [ ] Add tests with mock server

**Reference**: Python SDK v4 GCSClient
https://globus-sdk-python.readthedocs.io/en/stable/services/gcs.html

**Success Criteria**:
- Can connect to GCS endpoint by FQDN
- Can query endpoint information
- Error handling works correctly

**Estimate**: 1 week
**Persona**: System Administrator"

# Issue #8 - Endpoint show command
gh issue create \
  --repo scttfrdmn/globus-go-gcs \
  --title "[Phase 1] Implement \`endpoint show\` command" \
  --label "phase: 1-poc,priority: high,area: endpoint,persona: system-admin,good first issue" \
  --milestone "[Phase 1] Proof of Concept - Authentication Working" \
  --body "**Goal**: Implement basic endpoint command (proof of concept)

**Tasks**:
- [ ] Create internal/commands/endpoint/ package
- [ ] Implement \`endpoint show\` command
- [ ] Use GCS Manager API client
- [ ] Format output as table
- [ ] Add tests

**Success Criteria**:
- Can run \`globus-connect-server endpoint show\`
- Displays endpoint information
- Works with authenticated session

**Estimate**: 2 days
**Persona**: System Administrator"

# Issue #9 - JSON output
gh issue create \
  --repo scttfrdmn/globus-go-gcs \
  --title "[Phase 1] Add JSON output formatting" \
  --label "phase: 1-poc,priority: high,area: output,persona: it-manager" \
  --milestone "[Phase 1] Proof of Concept - Authentication Working" \
  --body "**Goal**: Support JSON output for all commands

**Tasks**:
- [ ] Create pkg/output/ package
- [ ] Implement JSON formatter
- [ ] Implement text table formatter
- [ ] Add \`--format=json\` flag
- [ ] Add tests

**Success Criteria**:
- All commands support \`--format=json\`
- JSON output is valid and parseable
- Table output is human-readable

**Estimate**: 2 days
**Persona**: IT Manager (for automation)"
```

## Step 6: Link Repository to Project

After creating issues, add them to the project board:

```bash
# Get project URL
PROJECT_URL="https://github.com/users/scttfrdmn/projects/$PROJECT_NUMBER"

# Add issues to project (issues 4-9)
for i in {4..9}; do
  gh project item-add $PROJECT_NUMBER \
    --owner scttfrdmn \
    --url "https://github.com/scttfrdmn/globus-go-gcs/issues/$i"
done
```

## Step 7: Set Up Repository Settings

In GitHub web UI (https://github.com/scttfrdmn/globus-go-gcs/settings):

1. **General**:
   - Description: "Complete Go port of Globus Connect Server v5 CLI - drop-in replacement for Python globus-connect-server"
   - Topics: `globus`, `go`, `cli`, `data-transfer`, `globus-connect-server`, `gcs`
   - Disable: Wikis, Projects (we use Projects v2 at user level)

2. **Actions**:
   - Enable GitHub Actions
   - Allow all actions

3. **Branches** (later, when collaborating):
   - Add branch protection rule for `main`
   - Require PR reviews
   - Require status checks to pass

## Verification

After setup, verify:

```bash
# Check repository
gh repo view scttfrdmn/globus-go-gcs

# Check labels
gh label list --repo scttfrdmn/globus-go-gcs

# Check milestones
gh milestone list --repo scttfrdmn/globus-go-gcs

# Check issues
gh issue list --repo scttfrdmn/globus-go-gcs

# Check CI status
gh run list --repo scttfrdmn/globus-go-gcs
```

## Next Steps

Once GitHub is set up:

1. ✅ Phase 0 is complete!
2. 🚀 Start Phase 1 development
3. Pick up Issue #4 (OAuth2 authentication)
4. Follow the development workflow in GETTING_STARTED.md

---

**Ready to push to GitHub and start Phase 1!** 🎉
