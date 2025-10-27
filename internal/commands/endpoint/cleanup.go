package endpoint

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewCleanupCmd creates the endpoint cleanup command.
func NewCleanupCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		force        bool
	)

	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Permanently remove endpoint configuration",
		Long: `Permanently remove the GCS endpoint configuration.

WARNING: This action is irreversible and will remove all endpoint
configuration. All collections and nodes must be removed first.

Example:
  globus-connect-server endpoint cleanup \
    --endpoint example.data.globus.org

Use --force to skip confirmation prompt.

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCleanup(cmd.Context(), profile, format, endpointFQDN, force, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runCleanup executes the endpoint cleanup command.
func runCleanup(ctx context.Context, profile, formatStr, endpointFQDN string, force bool, out interface{ Write([]byte) (int, error) }) error {
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
		if err := formatter.Println("WARNING: This will permanently remove the endpoint configuration."); err != nil {
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
		return fmt.Errorf("cleanup cancelled (use --force to proceed)")
	}

	// Create GCS client
	gcsClient, err := gcs.NewClient(
		endpointFQDN,
		gcs.WithAccessToken(token.AccessToken),
	)
	if err != nil {
		return fmt.Errorf("create GCS client: %w", err)
	}

	// Cleanup endpoint
	if err := gcsClient.CleanupEndpoint(ctx); err != nil {
		return fmt.Errorf("cleanup endpoint: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		result := map[string]string{
			"status":  "success",
			"message": "Endpoint cleaned up successfully",
		}
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.PrintText("Endpoint cleaned up successfully.\n"); err != nil {
		return err
	}

	return nil
}
