# Getting Started with globus-go-gcs Development

## Project Setup Complete! ✅

The following project infrastructure has been set up:

### Core Documents
- ✅ **PROJECT_PLAN.md** - Complete 28-week implementation roadmap
- ✅ **README.md** - Project overview and quick start
- ✅ **go.mod** - Go module initialized
- ✅ **LICENSE** - Apache 2.0
- ✅ **.gitignore** - Standard Go project ignores

### GitHub Infrastructure
- ✅ **.github/labels.yml** - 147 labels (type, priority, area, persona, phase, status)
- ✅ **.github/ISSUE_TEMPLATE/feature_request.yml** - Feature request with persona fields
- ✅ **.github/ISSUE_TEMPLATE/bug_report.yml** - Bug report template
- ✅ **.github/pull_request_template.md** - PR template with persona impact

### Key Decisions Made

1. **SDK Version**: Use globus-go-sdk **v3.65.0-1** (stable)
   - Implement GCS Manager API client in this project
   - Migrate to v4 later using module versioning

2. **GCS Manager API Client**: We'll implement it ourselves
   - Python SDK v4 has `GCSClient` but Go SDK doesn't
   - Based on `https://gcs.example.edu/api` per-endpoint pattern
   - Will contribute back to globus-go-sdk after proven

3. **Project Management**: Lens-style methodology
   - 4 personas: System Admin, Data Manager, Research PI, IT Manager
   - 6 phases: Setup → POC → Core → Collections → Advanced → Release
   - GitHub Project Board with custom fields

4. **Testing**: Have access to Globus subscription
   - Can deploy GCS locally and on AWS
   - Integration tests with real endpoints

---

## Next Steps

### 1. Initialize Git Repository (Immediate)

```bash
cd /Users/scttfrdmn/src/globus-go-gcs

# Initialize git
git init
git add .
git commit -m "Initial project setup

- Add PROJECT_PLAN.md with complete 28-week roadmap
- Add README.md with project overview
- Set up GitHub infrastructure (labels, issue templates, PR template)
- Initialize Go module
- Add .gitignore and LICENSE (Apache 2.0)

Phase 0 - Project Setup"

# Create GitHub repository and push
# (Follow GitHub's instructions for creating a new repo)
git remote add origin https://github.com/scttfrdmn/globus-go-gcs.git
git branch -M main
git push -u origin main
```

### 2. Set Up GitHub Project Board (Week 1)

Follow [PROJECT_PLAN.md](PROJECT_PLAN.md) Section "GitHub Project Management Setup":

1. **Create Project Board**
   - Go to https://github.com/users/scttfrdmn/projects
   - Create new project: "Globus Connect Server Go Port"

2. **Add Custom Fields** (programmatically using `gh` CLI):
   ```bash
   # Install gh CLI if needed: brew install gh

   # Authenticate
   gh auth login

   # Create project (note the project number)
   gh project create --owner scttfrdmn --title "Globus Connect Server Go Port"

   # Add custom fields (replace PROJECT_NUMBER with actual number)
   gh project field-create PROJECT_NUMBER --owner scttfrdmn \
     --name "Persona" --data-type "SINGLE_SELECT" \
     --single-select-options "System Administrator,Data Manager,Research PI,IT Manager"

   gh project field-create PROJECT_NUMBER --owner scttfrdmn \
     --name "Phase" --data-type "SINGLE_SELECT" \
     --single-select-options "Phase 0,Phase 1,Phase 2,Phase 3,Phase 4,Phase 5,Backlog"

   gh project field-create PROJECT_NUMBER --owner scttfrdmn \
     --name "Estimate" --data-type "NUMBER"

   gh project field-create PROJECT_NUMBER --owner scttfrdmn \
     --name "Priority" --data-type "SINGLE_SELECT" \
     --single-select-options "Critical,High,Medium,Low"
   ```

3. **Create Project Views** (manually in GitHub UI):
   - Kanban (default)
   - By Phase (roadmap view)
   - By Persona (user-centric)
   - Current Sprint (filtered)
   - Backlog (future work)

4. **Sync Labels**:
   ```bash
   # This will create all 147 labels from .github/labels.yml
   gh label sync --repo scttfrdmn/globus-go-gcs
   ```

5. **Create Milestones**:
   ```bash
   gh milestone create "[Phase 0] Project Infrastructure Complete" \
     --repo scttfrdmn/globus-go-gcs \
     --due-date 2025-11-15 \
     --description "Project setup, GitHub infra, planning complete"

   gh milestone create "[Phase 1] Proof of Concept - Authentication Working" \
     --repo scttfrdmn/globus-go-gcs \
     --due-date 2025-12-15 \
     --description "OAuth2 auth working, basic endpoint commands"

   # Create remaining milestones...
   ```

6. **Create Initial Issues** (see PROJECT_PLAN.md for issue list):
   ```bash
   # Phase 0 issues
   gh issue create --repo scttfrdmn/globus-go-gcs \
     --title "[Phase 0] Initialize Go module and project structure" \
     --label "phase: 0-setup,priority: high,area: build" \
     --milestone "[Phase 0] Project Infrastructure Complete" \
     --body "Set up basic Go project structure..."

   # Repeat for issues #2-6...
   ```

### 3. Begin Phase 1 Development (Week 2+)

Once GitHub infrastructure is complete:

1. **Add globus-go-sdk dependency**:
   ```bash
   go get github.com/scttfrdmn/globus-go-sdk/v3@latest
   go get github.com/spf13/cobra@latest
   go get github.com/spf13/viper@latest
   ```

2. **Create directory structure**:
   ```bash
   mkdir -p cmd/globus-connect-server
   mkdir -p pkg/client/gcs
   mkdir -p pkg/config
   mkdir -p pkg/models
   mkdir -p pkg/output
   mkdir -p internal/commands
   ```

3. **Implement Phase 1 POC** (see PROJECT_PLAN.md):
   - OAuth2 authentication flow
   - Token storage (compatible with Python CLI)
   - Config management
   - Basic `endpoint show` command
   - JSON output support

---

## Development Workflow

### Daily Workflow

1. **Pick an issue** from project board (status = "Ready")
2. **Move to "In Progress"** on project board
3. **Create feature branch**: `git checkout -b feature/issue-7-oauth2-login`
4. **Implement** with tests
5. **Open PR** using template (includes persona impact)
6. **CI runs** automatically (lint, test, build)
7. **Code review** and address feedback
8. **Merge** and close issue
9. **Update milestone** progress

### Creating Issues

Always use the issue templates:
- Feature Request: `.github/ISSUE_TEMPLATE/feature_request.yml`
- Bug Report: `.github/ISSUE_TEMPLATE/bug_report.yml`

Include:
- Persona(s) affected
- Phase alignment
- Area(s) impacted
- Priority level
- Clear success criteria

### Pull Requests

Always use the PR template:
- Fill in **Persona Impact Assessment** (required!)
- Link related issues (`Closes #7`)
- Add tests
- Update CHANGELOG.md
- Verify Python CLI compatibility (if applicable)

---

## Quick Reference

### Key Files
- **PROJECT_PLAN.md** - Complete roadmap (1200+ lines)
- **README.md** - Project overview
- **.github/labels.yml** - All 147 labels
- **.github/ISSUE_TEMPLATE/** - Issue templates
- **.github/pull_request_template.md** - PR template

### Key Commands
```bash
# Build
go build ./cmd/globus-connect-server

# Test
go test ./...

# Lint
golangci-lint run

# Create issue
gh issue create --web

# View project board
gh project view PROJECT_NUMBER
```

### Important Links
- **Python GCS CLI Docs**: https://docs.globus.org/globus-connect-server/v5/reference/
- **GCS Manager API Docs**: https://docs.globus.org/globus-connect-server/v5/api/
- **Python SDK GCSClient**: https://globus-sdk-python.readthedocs.io/en/stable/services/gcs.html
- **globus-go-sdk**: ../globus-go-sdk/
- **globus-go-cli**: ../globus-go-cli/

---

## Questions?

- Review [PROJECT_PLAN.md](PROJECT_PLAN.md) for detailed planning
- Check existing issues on GitHub
- Reference adjacent projects: globus-go-sdk, globus-go-cli, lens

---

**Ready to start coding!** 🚀

**First task**: Set up GitHub repository and project board, then begin Phase 1 (OAuth2 authentication).
