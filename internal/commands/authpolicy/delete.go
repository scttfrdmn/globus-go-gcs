package authpolicy

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewDeleteCmd creates the auth-policy delete command.
func NewDeleteCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		force        bool
	)

	cmd := &cobra.Command{
		Use:   "delete POLICY_ID",
		Short: "Delete an authentication policy",
		Long: `Delete an authentication policy from the endpoint.

WARNING: This action cannot be undone. Deleting a policy may affect
access controls for users relying on this policy.

Example:
  globus-connect-server auth-policy delete abc123 \
    --endpoint example.data.globus.org \
    --force

Use --force to skip confirmation prompt.

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			policyID := args[0]
			return runDelete(cmd.Context(), profile, format, endpointFQDN, policyID, force, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runDelete executes the auth-policy delete command.
func runDelete(ctx context.Context, profile, formatStr, endpointFQDN, policyID string, force bool, out interface{ Write([]byte) (int, error) }) error {
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

	// Confirmation prompt (unless --force)
	if !force {
		if err := formatter.PrintText("WARNING: This will permanently delete authentication policy %s.\n", policyID); err != nil {
			return err
		}
		if err := formatter.Println("This action cannot be undone."); err != nil {
			return err
		}
		if err := formatter.Println(); err != nil {
			return err
		}
		if err := formatter.PrintText("To proceed, use --force flag.\n"); err != nil {
			return err
		}
		return fmt.Errorf("delete cancelled (use --force to proceed)")
	}

	// Create GCS client
	gcsClient, err := gcs.NewClient(
		endpointFQDN,
		gcs.WithAccessToken(token.AccessToken),
	)
	if err != nil {
		return fmt.Errorf("create GCS client: %w", err)
	}

	// Delete policy
	if err := gcsClient.DeleteAuthPolicy(ctx, policyID); err != nil {
		return fmt.Errorf("delete auth policy: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		result := map[string]string{
			"status":    "success",
			"policy_id": policyID,
			"message":   "Authentication policy deleted successfully",
		}
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.PrintText("Authentication policy %s deleted successfully.\n", policyID); err != nil {
		return err
	}

	return nil
}
