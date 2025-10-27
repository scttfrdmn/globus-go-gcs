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

// NewResetOwnerStringCmd creates the endpoint reset-owner-string command.
func NewResetOwnerStringCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "reset-owner-string",
		Short: "Reset endpoint owner string to default (ClientID)",
		Long: `Reset the endpoint owner display name to the default value.

This command removes any custom owner string that was set using 'set-owner-string'
and reverts to displaying the ClientID as the owner name.

Example:
  globus-connect-server endpoint reset-owner-string \
    --endpoint example.data.globus.org

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runResetOwnerString(cmd.Context(), profile, format, endpointFQDN, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runResetOwnerString executes the endpoint reset-owner-string command.
func runResetOwnerString(ctx context.Context, profile, formatStr, endpointFQDN string, out interface{ Write([]byte) (int, error) }) error {
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

	// Reset owner string
	if err := gcsClient.ResetEndpointOwnerString(ctx); err != nil {
		return fmt.Errorf("reset endpoint owner string: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		result := map[string]string{
			"status":  "success",
			"message": "Endpoint owner string reset to default (ClientID)",
		}
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.PrintText("Endpoint owner string reset successfully.\n"); err != nil {
		return err
	}
	if err := formatter.PrintText("The owner display name is now set to the default (ClientID).\n"); err != nil {
		return err
	}

	return nil
}
