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

// NewSetupCmd creates the node setup command.
func NewSetupCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		name         string
		incoming     bool
		outgoing     bool
	)

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Configure and initialize a new node",
		Long: `Configure and initialize a new data transfer node.

A node represents a compute resource that handles data transfers
for the endpoint. This command sets up a new node with the specified configuration.

Example:
  globus-connect-server node setup \
    --endpoint example.data.globus.org \
    --name "transfer-node-1" \
    --incoming \
    --outgoing

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSetup(cmd.Context(), profile, format, endpointFQDN,
				name, incoming, outgoing, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&name, "name", "", "Name for the node")
	cmd.Flags().BoolVar(&incoming, "incoming", false, "Enable incoming transfers")
	cmd.Flags().BoolVar(&outgoing, "outgoing", false, "Enable outgoing transfers")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

// runSetup executes the node setup command.
func runSetup(ctx context.Context, profile, formatStr, endpointFQDN string,
	name string, incoming, outgoing bool,
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

	// Build node configuration
	node := &gcs.Node{
		Name:     name,
		Incoming: incoming,
		Outgoing: outgoing,
	}

	// Setup node
	created, err := gcsClient.SetupNode(ctx, node)
	if err != nil {
		return fmt.Errorf("setup node: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(created)
	}

	// Text format
	if err := formatter.Println("Node setup completed successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if created.ID != "" {
		if err := formatter.PrintText("%-20s%s\n", "Node ID:", created.ID); err != nil {
			return err
		}
	}
	if created.Name != "" {
		if err := formatter.PrintText("%-20s%s\n", "Name:", created.Name); err != nil {
			return err
		}
	}
	if err := formatter.PrintText("%-20s%v\n", "Incoming:", created.Incoming); err != nil {
		return err
	}
	if err := formatter.PrintText("%-20s%v\n", "Outgoing:", created.Outgoing); err != nil {
		return err
	}

	return nil
}
