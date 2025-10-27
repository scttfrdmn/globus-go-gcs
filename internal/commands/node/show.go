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

// NewShowCmd creates the node show command.
func NewShowCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "show NODE_ID",
		Short: "Display node details",
		Long: `Display detailed information about a specific node.

This command retrieves and displays the configuration of a node including
its name, ID, and transfer capabilities (incoming and outgoing).

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeID := args[0]
			return runShow(cmd.Context(), profile, format, endpointFQDN, nodeID, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runShow executes the node show command.
func runShow(ctx context.Context, profile, formatStr, endpointFQDN, nodeID string, out interface{ Write([]byte) (int, error) }) error {
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

	// Get node
	node, err := gcsClient.GetNode(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("get node: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(node)
	}

	// Text format
	if err := formatter.Println("Node Details:"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if node.Name != "" {
		if err := formatter.PrintText("%-15s%s\n", "Name:", node.Name); err != nil {
			return err
		}
	}
	if node.ID != "" {
		if err := formatter.PrintText("%-15s%s\n", "ID:", node.ID); err != nil {
			return err
		}
	}
	if err := formatter.PrintText("%-15s%t\n", "Incoming:", node.Incoming); err != nil {
		return err
	}
	if err := formatter.PrintText("%-15s%t\n", "Outgoing:", node.Outgoing); err != nil {
		return err
	}

	return nil
}
