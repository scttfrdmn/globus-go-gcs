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

// NewUpdateCmd creates the node update command.
func NewUpdateCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		name         string
		incoming     *bool
		outgoing     *bool
	)

	cmd := &cobra.Command{
		Use:   "update NODE_ID",
		Short: "Update an existing node",
		Long: `Update an existing node's configuration.

Only the fields you specify will be updated. Other fields will remain unchanged.

Example:
  globus-connect-server node update abc123 \
    --endpoint example.data.globus.org \
    --name "updated-node-name" \
    --incoming true

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeID := args[0]
			return runUpdate(cmd.Context(), profile, format, endpointFQDN, nodeID,
				name, incoming, outgoing, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&name, "name", "", "Name for the node")

	// Boolean flags need special handling
	incomingStr := cmd.Flags().String("incoming", "", "Enable incoming transfers (true/false)")
	cmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		if *incomingStr != "" {
			val := *incomingStr == "true"
			incoming = &val
		}
		outgoingStr, _ := cmd.Flags().GetString("outgoing")
		if outgoingStr != "" {
			val := outgoingStr == "true"
			outgoing = &val
		}
		return nil
	}

	cmd.Flags().String("outgoing", "", "Enable outgoing transfers (true/false)")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runUpdate executes the node update command.
func runUpdate(ctx context.Context, profile, formatStr, endpointFQDN, nodeID string,
	name string, incoming, outgoing *bool,
	out interface{ Write([]byte) (int, error) }) error {

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

	// Build node object with only the fields that were specified
	node := &gcs.Node{}

	if name != "" {
		node.Name = name
	}
	if incoming != nil {
		node.Incoming = *incoming
	}
	if outgoing != nil {
		node.Outgoing = *outgoing
	}

	// Update node
	updated, err := gcsClient.UpdateNode(ctx, nodeID, node)
	if err != nil {
		return fmt.Errorf("update node: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(updated)
	}

	// Text format
	if err := formatter.Println("Node updated successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if updated.ID != "" {
		if err := formatter.PrintText("%-20s%s\n", "ID:", updated.ID); err != nil {
			return err
		}
	}
	if updated.Name != "" {
		if err := formatter.PrintText("%-20s%s\n", "Name:", updated.Name); err != nil {
			return err
		}
	}
	if err := formatter.PrintText("%-20s%v\n", "Incoming:", updated.Incoming); err != nil {
		return err
	}
	if err := formatter.PrintText("%-20s%v\n", "Outgoing:", updated.Outgoing); err != nil {
		return err
	}

	return nil
}
