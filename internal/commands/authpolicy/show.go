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

// NewShowCmd creates the auth-policy show command.
func NewShowCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "show POLICY_ID",
		Short: "Display authentication policy details",
		Long: `Display detailed information about a specific authentication policy.

This command shows all configuration details for an authentication policy,
including security requirements, domain restrictions, and other settings.

Example:
  globus-connect-server auth-policy show abc123 \
    --endpoint example.data.globus.org

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			policyID := args[0]
			return runShow(cmd.Context(), profile, format, endpointFQDN, policyID, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runShow executes the auth-policy show command.
func runShow(ctx context.Context, profile, formatStr, endpointFQDN, policyID string, out interface{ Write([]byte) (int, error) }) error {
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

	// Get auth policy
	policy, err := gcsClient.GetAuthPolicy(ctx, policyID)
	if err != nil {
		return fmt.Errorf("get auth policy: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(policy)
	}

	// Text format
	if err := formatter.Println("Authentication Policy"); err != nil {
		return err
	}
	if err := formatter.Println("====================="); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if policy.ID != "" {
		if err := formatter.PrintText("%-25s%s\n", "ID:", policy.ID); err != nil {
			return err
		}
	}
	if policy.Name != "" {
		if err := formatter.PrintText("%-25s%s\n", "Name:", policy.Name); err != nil {
			return err
		}
	}
	if policy.Description != "" {
		if err := formatter.PrintText("%-25s%s\n", "Description:", policy.Description); err != nil {
			return err
		}
	}
	if err := formatter.PrintText("%-25s%v\n", "Require MFA:", policy.RequireMFA); err != nil {
		return err
	}
	if err := formatter.PrintText("%-25s%v\n", "Require High Assurance:", policy.RequireHighAssurance); err != nil {
		return err
	}

	if len(policy.AllowedDomains) > 0 {
		if err := formatter.PrintText("%-25s%s\n", "Allowed Domains:", strings.Join(policy.AllowedDomains, ", ")); err != nil {
			return err
		}
	}
	if len(policy.BlockedDomains) > 0 {
		if err := formatter.PrintText("%-25s%s\n", "Blocked Domains:", strings.Join(policy.BlockedDomains, ", ")); err != nil {
			return err
		}
	}

	return nil
}
