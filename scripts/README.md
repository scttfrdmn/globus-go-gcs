# Project Management Scripts

This directory contains scripts to set up GitHub Projects and Issues for the Security & Performance v2.0 initiative.

## Prerequisites

### Install GitHub CLI

```bash
# macOS
brew install gh

# Linux
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update
sudo apt install gh

# Windows
winget install --id GitHub.cli
```

### Authenticate with GitHub

```bash
gh auth login
```

Follow the prompts to authenticate with your GitHub account.

### Install jq (JSON processor)

```bash
# macOS
brew install jq

# Linux
sudo apt-get install jq

# Windows
winget install jqlang.jq
```

## Usage

### Step 1: Create GitHub Issues

This script creates all issues defined in the project plan, including:
- 5 Critical security issues (Milestone 1)
- 8 High priority security issues (Milestone 2)
- 5 High priority optimization issues (Milestone 3)
- 4 Testing/documentation issues (Milestone 4)

```bash
./scripts/create-github-issues.sh
```

**What it does**:
1. Creates 4 milestones with due dates
2. Creates ~20 issues with proper labels, descriptions, and assignments
3. Links issues to appropriate milestones
4. Sets up issue dependencies

**Output**:
```
Creating GitHub issues for Security & Performance v2.0 project...
Repository: scttfrdmn/globus-go-gcs

Creating milestones...
Milestones created.

Creating Epic: Security Remediation issues...
Creating issue: [SECURITY] Implement AES-256-GCM encryption for stored OAuth tokens
Creating issue: [SECURITY] Enforce TLS 1.2+ with secure cipher suites
...

✅ GitHub issues created successfully!
```

### Step 2: Set Up Project Board

This script creates a GitHub Projects (v2) board with custom fields and views.

```bash
./scripts/setup-github-project.sh
```

**What it does**:
1. Creates a new GitHub Project named "Security & Performance v2.0"
2. Adds custom fields:
   - **Priority**: P0-Critical, P1-High, P2-Medium, P3-Low
   - **Epic**: Security Remediation, Go Optimizations, Testing
   - **Effort**: Small (1-4h), Medium (4-16h), Large (16-40h), XL (40+h)
   - **Expected Impact**: Text field for performance/security impact
3. Links the repository to the project
4. Adds all issues to the project board

**Output**:
```
Setting up GitHub Project: Security & Performance v2.0

Step 1: Creating project...
✅ Project created: https://github.com/users/scttfrdmn/projects/1
   Project ID: PVT_...
   Project Number: 1

Step 2: Adding custom fields...
✅ Priority field added
✅ Epic field added
✅ Effort field added
✅ Expected Impact field added

Step 3: Creating views...
Default views created:
  - Table view (all issues)
  - Board view (by status)

Step 4: Linking repository...
✅ Repository linked to project

Step 5: Adding issues to project...
  Added issue #1
  Added issue #2
  ...

✅ GitHub Project setup complete!
```

### Step 3: Customize Project Board (Manual)

After running the setup script, visit your project board and create custom views:

1. **By Priority View**:
   - Click "New view" → "Board"
   - Name: "By Priority"
   - Group by: Priority
   - Sort: P0-Critical → P3-Low

2. **By Epic View**:
   - Click "New view" → "Board"
   - Name: "By Epic"
   - Group by: Epic
   - Shows swimlanes for Security, Optimizations, Testing

3. **By Milestone View**:
   - Click "New view" → "Table"
   - Name: "Current Sprint"
   - Filter: Milestone is "Milestone 1: Critical Security Fixes"
   - Group by: Status

4. **Roadmap View**:
   - Click "New view" → "Roadmap"
   - Name: "Timeline"
   - Group by: Milestone
   - Shows Gantt chart of all milestones

## Project Structure

### Milestones

| Milestone | Duration | Focus | Issues |
|-----------|----------|-------|--------|
| **Milestone 1**: Critical Security Fixes | Week 1-2 | CRITICAL severity security issues | #1-5 |
| **Milestone 2**: High Priority Security | Week 3-4 | HIGH severity security issues | #6-13 |
| **Milestone 3**: High Priority Optimizations | Week 5-6 | Performance improvements | #14-23 |
| **Milestone 4**: Testing & Documentation | Week 7-8 | Validation and docs | #24-28 |

### Labels

#### Priority
- `P0-Critical` - Security vulnerabilities, data loss risk (red)
- `P1-High` - Major security/performance issues (orange)
- `P2-Medium` - Moderate improvements (yellow)
- `P3-Low` - Nice-to-have enhancements (green)

#### Type
- `type: security` - Security remediation
- `type: performance` - Optimization work
- `type: testing` - Test coverage
- `type: documentation` - Docs updates

#### Area
- `area: auth` - Authentication & tokens
- `area: tls` - TLS/HTTPS configuration
- `area: audit` - Audit logging
- `area: api-client` - GCS API client
- `area: cli` - CLI commands
- `area: database` - SQLite operations

#### Effort
- `effort: small` - 1-4 hours
- `effort: medium` - 4-16 hours (1-2 days)
- `effort: large` - 16-40 hours (2-5 days)
- `effort: xl` - 40+ hours (5+ days)

#### Epic
- `epic: security-remediation` - NIST 800-53 compliance
- `epic: go-optimizations` - Performance improvements
- `epic: testing` - Test & validation

## Workflow

### Issue Lifecycle

```
Backlog → Ready → In Progress → In Review → Testing → Done
```

1. **Backlog**: Issue created, awaiting triage
2. **Ready**: Prioritized, ready to start work
3. **In Progress**: Actively being worked (assign yourself)
4. **In Review**: PR submitted, awaiting review
5. **Testing**: Merged to main, needs validation
6. **Done**: Completed and verified

### Working on an Issue

1. **Move to "Ready"**: Triage and prioritize
2. **Assign yourself**: Click "Assignees" → Add yourself
3. **Move to "In Progress"**: Start work
4. **Create feature branch**:
   ```bash
   git checkout -b feature/issue-1-token-encryption
   ```
5. **Make changes and commit**:
   ```bash
   git add .
   git commit -m "feat: implement AES-256-GCM token encryption (#1)"
   ```
6. **Create PR**:
   ```bash
   gh pr create --title "feat: implement AES-256-GCM token encryption" --body "Closes #1"
   ```
7. **Move to "In Review"**: PR automatically updates status
8. **Address review comments**: Make changes, push updates
9. **Merge PR**: Squash and merge
10. **Move to "Done"**: Issue automatically closes

### Branch Naming Convention

- **Security**: `security/issue-N-short-description`
- **Performance**: `perf/issue-N-short-description`
- **Testing**: `test/issue-N-short-description`
- **Documentation**: `docs/issue-N-short-description`

Examples:
```bash
git checkout -b security/issue-1-token-encryption
git checkout -b perf/issue-14-connection-pooling
git checkout -b test/issue-24-security-test-suite
git checkout -b docs/issue-28-hipaa-compliance-guide
```

### Commit Message Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `security`: Security fix
- `perf`: Performance improvement
- `test`: Adding tests
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `chore`: Maintenance tasks

**Examples**:
```
feat(auth): implement AES-256-GCM token encryption (#1)
security(tls): enforce TLS 1.2+ with secure ciphers (#2)
perf(api): add HTTP connection pooling (#14)
test(security): add comprehensive security test suite (#24)
docs(hipaa): add HIPAA compliance guide (#28)
```

## Sprint Planning

### Sprint 1 (Week 1-2): Critical Security - Milestone 1

**Goal**: Address all CRITICAL severity security issues

**Issues**: #1, #2, #3, #4, #5 (44 hours)

**Daily Standup Questions**:
1. What did you complete yesterday?
2. What will you work on today?
3. Any blockers?

**Sprint Planning**:
```bash
# Filter issues for current sprint
gh issue list --milestone "Milestone 1: Critical Security Fixes" --label "P0-Critical"

# Check sprint progress
gh issue list --milestone "Milestone 1: Critical Security Fixes" --json title,state,assignees
```

**Sprint Review**:
- Demo token encryption
- Demo secure secret input
- Review TLS configuration
- Discuss audit database encryption

**Sprint Retrospective**:
- What went well?
- What could be improved?
- Action items for next sprint

### Sprint 2 (Week 3-4): High Priority Security - Milestone 2

**Goal**: Complete HIGH severity security issues

**Issues**: #4, #6, #7, #8, #9, #10 (74 hours)

### Sprint 3 (Week 5-6): High Priority Optimizations - Milestone 3

**Goal**: Implement high-impact performance improvements

**Issues**: #14, #15, #16, #17, #18 (36 hours)

### Sprint 4 (Week 7-8): Testing & Documentation - Milestone 4

**Goal**: Comprehensive testing and documentation

**Issues**: #24, #25, #26, #28 (88 hours)

## Metrics & Reporting

### View Sprint Progress

```bash
# Issues by milestone
gh issue list --milestone "Milestone 1: Critical Security Fixes" --json number,title,state

# Issues by assignee
gh issue list --assignee @me --state open

# Issues by label
gh issue list --label "P0-Critical" --state open

# Recently closed issues
gh issue list --state closed --limit 10
```

### Generate Sprint Report

```bash
# Issues completed in last week
gh issue list --state closed --search "closed:>=$(date -v-7d +%Y-%m-%d)" --json number,title,closedAt

# Issues in progress
gh issue list --label "in-progress" --json number,title,assignees

# Blocked issues
gh issue list --label "blocked" --json number,title,body
```

### Track Burndown

Use GitHub's built-in insights:
1. Go to project board
2. Click "Insights" tab
3. View burndown chart by milestone
4. Track velocity over sprints

## CI/CD Integration

### Automated Issue Updates

GitHub Actions can automatically update issue status based on PR events:

```yaml
# .github/workflows/issue-automation.yml
name: Issue Automation
on:
  pull_request:
    types: [opened, closed, reopened]

jobs:
  update-issue:
    runs-on: ubuntu-latest
    steps:
      - name: Move to In Review
        if: github.event.action == 'opened'
        run: |
          gh issue edit ${{ github.event.pull_request.number }} \
            --add-label "in-review"

      - name: Move to Done
        if: github.event.action == 'closed' && github.event.pull_request.merged
        run: |
          gh issue edit ${{ github.event.pull_request.number }} \
            --add-label "done" \
            --remove-label "in-review"
```

## Troubleshooting

### Issues not appearing in project

```bash
# Manually add an issue to project
gh project item-add <PROJECT_NUMBER> --owner scttfrdmn --url https://github.com/scttfrdmn/globus-go-gcs/issues/1
```

### Permission errors

Ensure you have:
- Write access to the repository
- Admin access to create projects
- Valid authentication token

```bash
# Check authentication
gh auth status

# Refresh authentication
gh auth refresh
```

### Script fails with GraphQL errors

GitHub Projects v2 uses GraphQL API. Errors may occur if:
- Project already exists (check: https://github.com/users/scttfrdmn/projects)
- Field names conflict with existing fields
- Rate limits exceeded (wait and retry)

## Additional Resources

- [GitHub Projects Documentation](https://docs.github.com/en/issues/planning-and-tracking-with-projects)
- [GitHub CLI Documentation](https://cli.github.com/manual/)
- [GitHub GraphQL API](https://docs.github.com/en/graphql)
- [Conventional Commits](https://www.conventionalcommits.org/)

## Support

For issues with these scripts:
1. Check script output for error messages
2. Verify prerequisites are installed
3. Check GitHub authentication status
4. Review GitHub API rate limits
5. Open an issue in the repository

---

**Project Status**: Ready to start
**Total Issues**: 28 (5 critical, 8 high security + 5 high perf + 5 testing)
**Estimated Timeline**: 6-8 weeks
**Team Size**: 2 developers recommended
