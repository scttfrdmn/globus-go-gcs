package usercredential

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewListCmd creates the list command.
func NewListCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all user credentials",
		Long: `List all user storage credentials configured on the endpoint.

This command displays credentials for ActiveScale, OAuth2, and S3 storage
backends.

Example:
  globus-connect-server user-credential list \
    --endpoint example.data.globus.org

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runList(cmd.Context(), profile, format, endpointFQDN, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runList executes the list command.
func runList(ctx context.Context, profile, formatStr, endpointFQDN string, out interface{ Write([]byte) (int, error) }) error {
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

	// List credentials
	credentials, err := gcsClient.ListUserCredentials(ctx)
	if err != nil {
		return fmt.Errorf("list user credentials: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(credentials)
	}

	// Text format
	if len(credentials.Data) == 0 {
		return formatter.Println("No user credentials found")
	}

	if err := formatter.Println("User Credentials"); err != nil {
		return err
	}
	if err := formatter.Println("================"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	for _, credential := range credentials.Data {
		if credential.ID != "" {
			if err := formatter.PrintText("ID:               %s\n", credential.ID); err != nil {
				return err
			}
		}
		if credential.Type != "" {
			if err := formatter.PrintText("Type:             %s\n", credential.Type); err != nil {
				return err
			}
		}
		if credential.IdentityID != "" {
			if err := formatter.PrintText("Identity:         %s\n", credential.IdentityID); err != nil {
				return err
			}
		}
		if credential.StorageGatewayID != "" {
			if err := formatter.PrintText("Storage Gateway:  %s\n", credential.StorageGatewayID); err != nil {
				return err
			}
		}
		if credential.Username != "" {
			if err := formatter.PrintText("Username:         %s\n", credential.Username); err != nil {
				return err
			}
		}
		if err := formatter.Println(); err != nil {
			return err
		}
	}

	return nil
}
