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

// NewSetOwnerStringCmd creates the endpoint set-owner-string command.
func NewSetOwnerStringCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		ownerString  string
	)

	cmd := &cobra.Command{
		Use:   "set-owner-string",
		Short: "Set custom display name for endpoint owner",
		Long: `Set a custom display name for the endpoint owner.

This command allows you to override the default owner display name (which is
the ClientID) with a custom string. This is useful for providing a more
user-friendly display name that will appear in the Globus web interface.

To reset to the default, use the 'reset-owner-string' command.

Example:
  globus-connect-server endpoint set-owner-string \
    --endpoint example.data.globus.org \
    --owner-string "Research Data Team"

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSetOwnerString(cmd.Context(), profile, format, endpointFQDN, ownerString, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&ownerString, "owner-string", "", "Custom owner display name")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("owner-string")

	return cmd
}

// runSetOwnerString executes the endpoint set-owner-string command.
func runSetOwnerString(ctx context.Context, profile, formatStr, endpointFQDN, ownerString string, out interface{ Write([]byte) (int, error) }) error {
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

	// Set owner string
	if err := gcsClient.SetEndpointOwnerString(ctx, ownerString); err != nil {
		return fmt.Errorf("set endpoint owner string: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		result := map[string]string{
			"status":       "success",
			"owner_string": ownerString,
			"message":      "Endpoint owner string set successfully",
		}
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.PrintText("Endpoint owner string set successfully.\n"); err != nil {
		return err
	}
	if err := formatter.PrintText("Owner String: %s\n", ownerString); err != nil {
		return err
	}

	return nil
}
