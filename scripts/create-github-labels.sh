#!/bin/bash

# Script to create GitHub labels for Security & Performance v2.0 project
# Prerequisites: gh CLI tool installed and authenticated
# Usage: ./scripts/create-github-labels.sh

set -e

REPO="scttfrdmn/globus-go-gcs"

echo "Creating GitHub labels for Security & Performance v2.0 project..."
echo "Repository: $REPO"
echo ""

# Function to create a label
create_label() {
    local name="$1"
    local color="$2"
    local description="$3"

    echo "Creating label: $name"
    gh label create "$name" \
        --repo "$REPO" \
        --color "$color" \
        --description "$description" \
        --force 2>/dev/null || echo "  (already exists)"
}

echo "==> Priority Labels"
create_label "P0-Critical" "b60205" "Security vulnerabilities, data loss risk"
create_label "P1-High" "d93f0b" "Major security/performance issues"
create_label "P2-Medium" "fbca04" "Moderate improvements"
create_label "P3-Low" "0e8a16" "Nice-to-have enhancements"

echo ""
echo "==> Type Labels"
create_label "type: security" "d73a4a" "Security remediation"
create_label "type: performance" "1d76db" "Optimization work"
create_label "type: testing" "5319e7" "Test coverage"
create_label "type: documentation" "0075ca" "Documentation updates"

echo ""
echo "==> Area Labels"
create_label "area: auth" "c2e0c6" "Authentication & tokens"
create_label "area: tls" "c2e0c6" "TLS/HTTPS configuration"
create_label "area: audit" "c2e0c6" "Audit logging"
create_label "area: api-client" "c2e0c6" "GCS API client"
create_label "area: cli" "c2e0c6" "CLI commands"
create_label "area: database" "c2e0c6" "SQLite operations"

echo ""
echo "==> Effort Labels"
create_label "effort: small" "d4c5f9" "1-4 hours"
create_label "effort: medium" "c5def5" "4-16 hours (1-2 days)"
create_label "effort: large" "f9d0c4" "16-40 hours (2-5 days)"
create_label "effort: xl" "ffd8b1" "40+ hours (5+ days)"

echo ""
echo "==> Epic Labels"
create_label "epic: security-remediation" "d73a4a" "NIST 800-53 compliance"
create_label "epic: go-optimizations" "1d76db" "Performance improvements"
create_label "epic: testing" "5319e7" "Test & validation"

echo ""
echo "==> Special Labels"
create_label "breaking-change" "b60205" "Introduces breaking changes - requires migration guide"

echo ""
echo "✅ All labels created successfully!"
echo ""
echo "Next steps:"
echo "1. Re-run: ./scripts/create-github-issues.sh (to update issue labels)"
echo "2. Or manually update issues: gh issue edit <number> --add-label \"P0-Critical\""
