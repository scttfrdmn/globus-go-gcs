package authpolicy

import (
	"context"
	"fmt"
	"strings"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewUpdateCmd creates the auth-policy update command.
func NewUpdateCmd() *cobra.Command {
	var (
		profile              string
		format               string
		endpointFQDN         string
		name                 string
		description          string
		requireMFA           string
		requireHighAssurance string
		allowedDomains       string
		blockedDomains       string
	)

	cmd := &cobra.Command{
		Use:   "update POLICY_ID",
		Short: "Update an authentication policy",
		Long: `Update an existing authentication policy.

Only specified fields will be updated. Omitted fields remain unchanged.

Example:
  globus-connect-server auth-policy update abc123 \
    --endpoint example.data.globus.org \
    --name "Updated Policy Name" \
    --require-mfa=true

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			policyID := args[0]
			return runUpdate(cmd.Context(), profile, format, endpointFQDN, policyID, name, description,
				requireMFA, requireHighAssurance, allowedDomains, blockedDomains, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&name, "name", "", "Policy name")
	cmd.Flags().StringVar(&description, "description", "", "Policy description")
	cmd.Flags().StringVar(&requireMFA, "require-mfa", "", "Require multi-factor authentication (true/false)")
	cmd.Flags().StringVar(&requireHighAssurance, "require-high-assurance", "", "Require high assurance authentication (true/false)")
	cmd.Flags().StringVar(&allowedDomains, "allowed-domains", "", "Comma-separated list of allowed domains")
	cmd.Flags().StringVar(&blockedDomains, "blocked-domains", "", "Comma-separated list of blocked domains")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runUpdate executes the auth-policy update command.
func runUpdate(ctx context.Context, profile, formatStr, endpointFQDN, policyID, name, description string,
	requireMFA, requireHighAssurance, allowedDomains, blockedDomains string,
	out interface{ Write([]byte) (int, error) }) error {

	// Load token
	token, err := auth.LoadToken(profile)
	if err != nil {
		return fmt.Errorf("not logged in: %w (use 'login' command first)", err)
	}

	// Check if token is valid
	if !token.IsValid() {
		return fmt.Errorf("token expired, please login again")
	}

	// Create output formatter
	formatter := output.NewFormatter(output.Format(formatStr), out)

	// Create GCS client
	gcsClient, err := gcs.NewClient(
		endpointFQDN,
		gcs.WithAccessToken(token.AccessToken),
	)
	if err != nil {
		return fmt.Errorf("create GCS client: %w", err)
	}

	// Build update
	policy := &gcs.AuthPolicy{}

	if name != "" {
		policy.Name = name
	}
	if description != "" {
		policy.Description = description
	}

	// Parse boolean flags
	if requireMFA != "" {
		mfa := requireMFA == "true"
		policy.RequireMFA = mfa
	}
	if requireHighAssurance != "" {
		ha := requireHighAssurance == "true"
		policy.RequireHighAssurance = ha
	}

	// Parse domains
	if allowedDomains != "" {
		domains := strings.Split(allowedDomains, ",")
		for i, d := range domains {
			domains[i] = strings.TrimSpace(d)
		}
		policy.AllowedDomains = domains
	}

	if blockedDomains != "" {
		domains := strings.Split(blockedDomains, ",")
		for i, d := range domains {
			domains[i] = strings.TrimSpace(d)
		}
		policy.BlockedDomains = domains
	}

	// Update policy
	updated, err := gcsClient.UpdateAuthPolicy(ctx, policyID, policy)
	if err != nil {
		return fmt.Errorf("update auth policy: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(updated)
	}

	// Text format
	if err := formatter.Println("Authentication policy updated successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}
	if updated.ID != "" {
		if err := formatter.PrintText("%-20s%s\n", "ID:", updated.ID); err != nil {
			return err
		}
	}
	if updated.Name != "" {
		if err := formatter.PrintText("%-20s%s\n", "Name:", updated.Name); err != nil {
			return err
		}
	}

	return nil
}
