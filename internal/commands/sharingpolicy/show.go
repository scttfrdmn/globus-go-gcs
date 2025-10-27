package sharingpolicy

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

// NewShowCmd creates the sharing-policy show command.
func NewShowCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "show POLICY_ID",
		Short: "Display details of a sharing policy",
		Long: `Display detailed information about a specific sharing policy.

Shows all configuration including restrictions, allowed/denied users and groups.

Example:
  globus-connect-server sharing-policy show abc123 \
    --endpoint example.data.globus.org

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShow(cmd.Context(), profile, format, endpointFQDN, args[0], cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runShow executes the sharing-policy show command.
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

	// Get sharing policy
	policy, err := gcsClient.GetSharingPolicy(ctx, policyID)
	if err != nil {
		return fmt.Errorf("get sharing policy: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(policy)
	}

	// Text format
	if err := formatter.Println("Sharing Policy Details"); err != nil {
		return err
	}
	if err := formatter.Println("======================"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if policy.ID != "" {
		if err := formatter.PrintText("%-20s%s\n", "ID:", policy.ID); err != nil {
			return err
		}
	}
	if policy.Name != "" {
		if err := formatter.PrintText("%-20s%s\n", "Name:", policy.Name); err != nil {
			return err
		}
	}
	if policy.CollectionID != "" {
		if err := formatter.PrintText("%-20s%s\n", "Collection ID:", policy.CollectionID); err != nil {
			return err
		}
	}
	if policy.Description != "" {
		if err := formatter.PrintText("%-20s%s\n", "Description:", policy.Description); err != nil {
			return err
		}
	}
	if policy.SharingRestrict != "" {
		if err := formatter.PrintText("%-20s%s\n", "Sharing Restrict:", policy.SharingRestrict); err != nil {
			return err
		}
	}
	if len(policy.SharingUsersAllow) > 0 {
		if err := formatter.PrintText("%-20s%s\n", "Users Allow:", strings.Join(policy.SharingUsersAllow, ", ")); err != nil {
			return err
		}
	}
	if len(policy.SharingUsersDeny) > 0 {
		if err := formatter.PrintText("%-20s%s\n", "Users Deny:", strings.Join(policy.SharingUsersDeny, ", ")); err != nil {
			return err
		}
	}
	if len(policy.SharingGroupsAllow) > 0 {
		if err := formatter.PrintText("%-20s%s\n", "Groups Allow:", strings.Join(policy.SharingGroupsAllow, ", ")); err != nil {
			return err
		}
	}
	if len(policy.SharingGroupsDeny) > 0 {
		if err := formatter.PrintText("%-20s%s\n", "Groups Deny:", strings.Join(policy.SharingGroupsDeny, ", ")); err != nil {
			return err
		}
	}

	return nil
}
