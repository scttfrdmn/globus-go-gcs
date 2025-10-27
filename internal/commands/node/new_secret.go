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

// NewNewSecretCmd creates the node new-secret command.
func NewNewSecretCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "new-secret NODE_ID",
		Short: "Generate a new authentication secret for a node",
		Long: `Generate a new authentication secret for a node.

This command creates a new secret that the node uses to authenticate
with the GCS Manager API. The old secret will be invalidated.

IMPORTANT: Save the new secret securely immediately. It will not be
displayed again and cannot be recovered.

Example:
  globus-connect-server node new-secret abc123 \
    --endpoint example.data.globus.org

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeID := args[0]
			return runNewSecret(cmd.Context(), profile, format, endpointFQDN, nodeID, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runNewSecret executes the node new-secret command.
func runNewSecret(ctx context.Context, profile, formatStr, endpointFQDN, nodeID string, out interface{ Write([]byte) (int, error) }) error {
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

	// Generate new secret
	result, err := gcsClient.GenerateNodeSecret(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("generate node secret: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.Println("New node secret generated successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if result.NodeID != "" {
		if err := formatter.PrintText("%-20s%s\n", "Node ID:", result.NodeID); err != nil {
			return err
		}
	}
	if result.Secret != "" {
		if err := formatter.PrintText("%-20s%s\n", "Secret:", result.Secret); err != nil {
			return err
		}
	}

	if err := formatter.Println(); err != nil {
		return err
	}
	if err := formatter.Println("IMPORTANT: Save this secret securely. The old secret is now invalid."); err != nil {
		return err
	}
	if err := formatter.Println("This secret will not be shown again."); err != nil {
		return err
	}

	return nil
}
