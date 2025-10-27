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

// NewCleanupCmd creates the node cleanup command.
func NewCleanupCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		force        bool
	)

	cmd := &cobra.Command{
		Use:   "cleanup NODE_ID",
		Short: "Remove node and its configuration",
		Long: `Remove a node and its configuration from the endpoint.

WARNING: This action is permanent and cannot be undone.

Example:
  globus-connect-server node cleanup abc123 \
    --endpoint example.data.globus.org

Use --force to skip confirmation prompt.

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeID := args[0]
			return runCleanup(cmd.Context(), profile, format, endpointFQDN, nodeID, force, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runCleanup executes the node cleanup command.
func runCleanup(ctx context.Context, profile, formatStr, endpointFQDN, nodeID string, force bool, out interface{ Write([]byte) (int, error) }) error {
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
		if err := formatter.PrintText("WARNING: This will permanently remove node %s.\n", nodeID); err != nil {
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

	// Cleanup node
	if err := gcsClient.CleanupNode(ctx, nodeID); err != nil {
		return fmt.Errorf("cleanup node: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		result := map[string]string{
			"status":  "success",
			"node_id": nodeID,
			"message": "Node cleaned up successfully",
		}
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.PrintText("Node %s cleaned up successfully.\n", nodeID); err != nil {
		return err
	}

	return nil
}
