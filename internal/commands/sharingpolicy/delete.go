package sharingpolicy

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewDeleteCmd creates the sharing-policy delete command.
func NewDeleteCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		force        bool
	)

	cmd := &cobra.Command{
		Use:   "delete POLICY_ID",
		Short: "Delete a sharing policy",
		Long: `Delete a sharing policy.

This removes the sharing policy configuration. Use --force to skip the
confirmation prompt.

Example:
  globus-connect-server sharing-policy delete abc123 \
    --endpoint example.data.globus.org

  # Skip confirmation
  globus-connect-server sharing-policy delete abc123 \
    --endpoint example.data.globus.org \
    --force

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(cmd.Context(), profile, format, endpointFQDN, args[0], force, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runDelete executes the sharing-policy delete command.
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

	// Confirmation prompt unless --force
	if !force {
		fmt.Fprintf(os.Stderr, "Are you sure you want to delete sharing policy %s? (yes/no): ", policyID)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read confirmation: %w", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "yes" && response != "y" {
			return fmt.Errorf("deletion cancelled")
		}
	}

	// Create GCS client
	gcsClient, err := gcs.NewClient(
		endpointFQDN,
		gcs.WithAccessToken(token.AccessToken),
	)
	if err != nil {
		return fmt.Errorf("create GCS client: %w", err)
	}

	// Delete sharing policy
	if err := gcsClient.DeleteSharingPolicy(ctx, policyID); err != nil {
		return fmt.Errorf("delete sharing policy: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(map[string]string{"status": "deleted", "policy_id": policyID})
	}

	// Text format
	return formatter.PrintText("Sharing policy %s deleted successfully\n", policyID)
}
