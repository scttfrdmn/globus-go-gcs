package node

import "github.com/spf13/cobra"

// NewNodeCmd creates the node command with subcommands.
func NewNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Manage GCS nodes",
		Long: `Commands for managing Globus Connect Server nodes.

Nodes are compute resources that handle data transfer operations for
an endpoint. Each node can be configured to handle incoming transfers,
outgoing transfers, or both.`,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewShowCmd())

	return cmd
}
