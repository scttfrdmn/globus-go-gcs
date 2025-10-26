// Package collection provides collection management commands for the GCS CLI.
package collection

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewListCmd creates the collection list command.
func NewListCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		filter       string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List collections on an endpoint",
		Long: `List all collections configured on a Globus Connect Server endpoint.

This command retrieves and displays all collections (both mapped and guest)
on the specified endpoint. Collections can be filtered by name.

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runList(cmd.Context(), profile, format, endpointFQDN, filter, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter collections by name")
	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runList executes the collection list command.
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
	opts := &gcs.ListCollectionsOptions{}
	if filter != "" {
		opts.Filter = filter
	}

	// Get collections
	list, err := gcsClient.ListCollections(ctx, opts)
	if err != nil {
		return fmt.Errorf("list collections: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(list)
	}

	// Text format
	if len(list.Data) == 0 {
		if err := formatter.Println("No collections found."); err != nil {
			return err
		}
		return nil
	}

	if err := formatter.PrintText("Collections (%d):\n\n", len(list.Data)); err != nil {
		return err
	}

	for i, collection := range list.Data {
		if i > 0 {
			if err := formatter.Println(); err != nil {
				return err
			}
		}

		if err := formatter.PrintText("  Display Name:     %s\n", collection.DisplayName); err != nil {
			return err
		}
		if err := formatter.PrintText("  ID:               %s\n", collection.ID); err != nil {
			return err
		}
		if err := formatter.PrintText("  Type:             %s\n", collection.CollectionType); err != nil {
			return err
		}
		if collection.StorageGatewayID != "" {
			if err := formatter.PrintText("  Storage Gateway:  %s\n", collection.StorageGatewayID); err != nil {
				return err
			}
		}
		if err := formatter.PrintText("  Public:           %t\n", collection.Public); err != nil {
			return err
		}
	}

	return nil
}
