// Package storagegateway provides storage gateway management commands for the GCS CLI.
package storagegateway

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewListCmd creates the storage gateway list command.
func NewListCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		filter       string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List storage gateways on an endpoint",
		Long: `List all storage gateways configured on a Globus Connect Server endpoint.

Storage gateways define the connection between GCS and storage backends
(POSIX filesystems, S3, Azure Blob, Google Cloud Storage, etc.).

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runList(cmd.Context(), profile, format, endpointFQDN, filter, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter storage gateways by name")
	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runList executes the storage gateway list command.
func runList(ctx context.Context, profile, formatStr, endpointFQDN, filter string, out interface{ Write([]byte) (int, error) }) error {
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

	// Build list options
	opts := &gcs.ListStorageGatewaysOptions{}
	if filter != "" {
		opts.Filter = filter
	}

	// Get storage gateways
	list, err := gcsClient.ListStorageGateways(ctx, opts)
	if err != nil {
		return fmt.Errorf("list storage gateways: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(list)
	}

	// Text format
	if len(list.Data) == 0 {
		if err := formatter.Println("No storage gateways found."); err != nil {
			return err
		}
		return nil
	}

	if err := formatter.PrintText("Storage Gateways (%d):\n\n", len(list.Data)); err != nil {
		return err
	}

	for i, gateway := range list.Data {
		if i > 0 {
			if err := formatter.Println(); err != nil {
				return err
			}
		}

		if err := formatter.PrintText("  Display Name:     %s\n", gateway.DisplayName); err != nil {
			return err
		}
		if err := formatter.PrintText("  ID:               %s\n", gateway.ID); err != nil {
			return err
		}
		if err := formatter.PrintText("  Connector:        %s (%s)\n", gateway.ConnectorName, gateway.ConnectorID); err != nil {
			return err
		}
		if gateway.Root != "" {
			if err := formatter.PrintText("  Root:             %s\n", gateway.Root); err != nil {
				return err
			}
		}
		if gateway.HighAssurance {
			if err := formatter.PrintText("  High Assurance:   %t\n", gateway.HighAssurance); err != nil {
				return err
			}
		}
	}

	return nil
}
