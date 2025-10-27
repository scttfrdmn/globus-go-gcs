#!/bin/bash

# Script to set up GitHub Projects board for Security & Performance v2.0
# Prerequisites: gh CLI tool with projects extension
# Usage: ./scripts/setup-github-project.sh

set -e

REPO="scttfrdmn/globus-go-gcs"
OWNER="scttfrdmn"
PROJECT_NAME="Security & Performance v2.0"

echo "Setting up GitHub Project: $PROJECT_NAME"
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

echo "Step 1: Creating project..."
# Note: GitHub Projects v2 uses GraphQL API
# We'll use gh api for this

# Create project
PROJECT_JSON=$(gh api graphql -f query='
mutation {
  createProjectV2(input: {
    ownerId: "'$(gh api user -q .node_id)'"
    title: "'"$PROJECT_NAME"'"
  }) {
    projectV2 {
      id
      number
      url
    }
  }
}' 2>/dev/null || echo '{"data":{"createProjectV2":null}}')

PROJECT_ID=$(echo "$PROJECT_JSON" | jq -r '.data.createProjectV2.projectV2.id // empty')
PROJECT_NUMBER=$(echo "$PROJECT_JSON" | jq -r '.data.createProjectV2.projectV2.number // empty')
PROJECT_URL=$(echo "$PROJECT_JSON" | jq -r '.data.createProjectV2.projectV2.url // empty')

if [ -z "$PROJECT_ID" ]; then
    echo "Error: Failed to create project. It may already exist."
    echo "Please check: https://github.com/users/$OWNER/projects"
    exit 1
fi

echo "✅ Project created: $PROJECT_URL"
echo "   Project ID: $PROJECT_ID"
echo "   Project Number: $PROJECT_NUMBER"
echo ""

echo "Step 2: Adding custom fields..."

# Add Priority field (single select)
gh api graphql -f query='
mutation($projectId: ID!) {
  createProjectV2Field(input: {
    projectId: $projectId
    dataType: SINGLE_SELECT
    name: "Priority"
    singleSelectOptions: [
      {name: "P0-Critical", color: RED, description: "Critical security/data loss"}
      {name: "P1-High", color: ORANGE, description: "Major issues"}
      {name: "P2-Medium", color: YELLOW, description: "Moderate improvements"}
      {name: "P3-Low", color: GREEN, description: "Nice-to-have"}
    ]
  }) {
    projectV2Field {
      ... on ProjectV2SingleSelectField {
        id
        name
      }
    }
  }
}' -f projectId="$PROJECT_ID"

echo "✅ Priority field added"

# Add Epic field (single select)
gh api graphql -f query='
mutation($projectId: ID!) {
  createProjectV2Field(input: {
    projectId: $projectId
    dataType: SINGLE_SELECT
    name: "Epic"
    singleSelectOptions: [
      {name: "Security Remediation", color: RED}
      {name: "Go Optimizations", color: BLUE}
      {name: "Testing", color: PURPLE}
    ]
  }) {
    projectV2Field {
      ... on ProjectV2SingleSelectField {
        id
        name
      }
    }
  }
}' -f projectId="$PROJECT_ID"

echo "✅ Epic field added"

# Add Effort field (single select)
gh api graphql -f query='
mutation($projectId: ID!) {
  createProjectV2Field(input: {
    projectId: $projectId
    dataType: SINGLE_SELECT
    name: "Effort"
    singleSelectOptions: [
      {name: "Small (1-4h)", color: GREEN}
      {name: "Medium (4-16h)", color: YELLOW}
      {name: "Large (16-40h)", color: ORANGE}
      {name: "XL (40+h)", color: RED}
    ]
  }) {
    projectV2Field {
      ... on ProjectV2SingleSelectField {
        id
        name
      }
    }
  }
}' -f projectId="$PROJECT_ID"

echo "✅ Effort field added"

# Add Expected Impact field (text)
gh api graphql -f query='
mutation($projectId: ID!) {
  createProjectV2Field(input: {
    projectId: $projectId
    dataType: TEXT
    name: "Expected Impact"
  }) {
    projectV2Field {
      ... on ProjectV2Field {
        id
        name
      }
    }
  }
}' -f projectId="$PROJECT_ID"

echo "✅ Expected Impact field added"
echo ""

echo "Step 3: Creating views..."

# Note: Views are automatically created with the project
# Users can customize these in the GitHub UI

echo "Default views created:"
echo "  - Table view (all issues)"
echo "  - Board view (by status)"
echo ""
echo "Recommended custom views to create in the UI:"
echo "  1. 'By Priority' - Group by Priority field"
echo "  2. 'By Epic' - Group by Epic field"
echo "  3. 'By Milestone' - Group by Milestone"
echo "  4. 'Current Sprint' - Filter by current milestone"
echo ""

echo "Step 4: Linking repository..."

# Link the repository to the project
REPO_ID=$(gh api "repos/$REPO" -q .node_id)

gh api graphql -f query='
mutation($projectId: ID!, $repositoryId: ID!) {
  linkProjectV2ToRepository(input: {
    projectId: $projectId
    repositoryId: $repositoryId
  }) {
    repository {
      id
    }
  }
}' -f projectId="$PROJECT_ID" -f repositoryId="$REPO_ID"

echo "✅ Repository linked to project"
echo ""

echo "Step 5: Adding issues to project..."

# Get all issues with the epic labels
ISSUE_NUMBERS=$(gh issue list --repo "$REPO" --label "epic: security-remediation" --json number --jq '.[].number')

for issue_num in $ISSUE_NUMBERS; do
    ISSUE_ID=$(gh api "repos/$REPO/issues/$issue_num" -q .node_id)

    gh api graphql -f query='
    mutation($projectId: ID!, $contentId: ID!) {
      addProjectV2ItemById(input: {
        projectId: $projectId
        contentId: $contentId
      }) {
        item {
          id
        }
      }
    }' -f projectId="$PROJECT_ID" -f contentId="$ISSUE_ID" 2>/dev/null || true

    echo "  Added issue #$issue_num"
done

# Add optimization issues
ISSUE_NUMBERS=$(gh issue list --repo "$REPO" --label "epic: go-optimizations" --json number --jq '.[].number')

for issue_num in $ISSUE_NUMBERS; do
    ISSUE_ID=$(gh api "repos/$REPO/issues/$issue_num" -q .node_id)

    gh api graphql -f query='
    mutation($projectId: ID!, $contentId: ID!) {
      addProjectV2ItemById(input: {
        projectId: $projectId
        contentId: $contentId
      }) {
        item {
          id
        }
      }
    }' -f projectId="$PROJECT_ID" -f contentId="$ISSUE_ID" 2>/dev/null || true

    echo "  Added issue #$issue_num"
done

# Add testing issues
ISSUE_NUMBERS=$(gh issue list --repo "$REPO" --label "epic: testing" --json number --jq '.[].number')

for issue_num in $ISSUE_NUMBERS; do
    ISSUE_ID=$(gh api "repos/$REPO/issues/$issue_num" -q .node_id)

    gh api graphql -f query='
    mutation($projectId: ID!, $contentId: ID!) {
      addProjectV2ItemById(input: {
        projectId: $projectId
        contentId: $contentId
      }) {
        item {
          id
        }
      }
    }' -f projectId="$PROJECT_ID" -f contentId="$ISSUE_ID" 2>/dev/null || true

    echo "  Added issue #$issue_num"
done

echo "✅ All issues added to project"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ GitHub Project setup complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Project URL: $PROJECT_URL"
echo ""
echo "Next steps:"
echo "1. Visit the project board: $PROJECT_URL"
echo "2. Customize views in the UI:"
echo "   - Create 'By Priority' view (group by Priority)"
echo "   - Create 'By Epic' view (group by Epic)"
echo "   - Create 'By Milestone' view (group by Milestone)"
echo "3. Start working on Milestone 1 issues"
echo "4. Move issues through the workflow: Backlog → Ready → In Progress → In Review → Done"
echo ""
echo "Workflow tips:"
echo "- Use draft PRs for work in progress"
echo "- Link PRs to issues with 'Closes #N' in PR description"
echo "- Update issue status as work progresses"
echo "- Use project views to track sprint progress"
echo ""
