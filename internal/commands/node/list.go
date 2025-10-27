// Package node provides node management commands for the GCS CLI.
package node

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewListCmd creates the node list command.
func NewListCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		filter       string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List nodes on an endpoint",
		Long: `List all nodes configured on a Globus Connect Server endpoint.

Nodes represent the compute resources that handle data transfer operations
for the endpoint. Each node can be configured for incoming and/or outgoing
transfers.

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runList(cmd.Context(), profile, format, endpointFQDN, filter, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter nodes by name")
	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runList executes the node list command.
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
	opts := &gcs.ListNodesOptions{}
	if filter != "" {
		opts.Filter = filter
	}

	// Get nodes
	list, err := gcsClient.ListNodes(ctx, opts)
	if err != nil {
		return fmt.Errorf("list nodes: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(list)
	}

	// Text format
	if len(list.Data) == 0 {
		if err := formatter.Println("No nodes found."); err != nil {
			return err
		}
		return nil
	}

	if err := formatter.PrintText("Nodes (%d):\n\n", len(list.Data)); err != nil {
		return err
	}

	for i, node := range list.Data {
		if i > 0 {
			if err := formatter.Println(); err != nil {
				return err
			}
		}

		if err := formatter.PrintText("  Name:         %s\n", node.Name); err != nil {
			return err
		}
		if err := formatter.PrintText("  ID:           %s\n", node.ID); err != nil {
			return err
		}
		if err := formatter.PrintText("  Incoming:     %t\n", node.Incoming); err != nil {
			return err
		}
		if err := formatter.PrintText("  Outgoing:     %t\n", node.Outgoing); err != nil {
			return err
		}
	}

	return nil
}
