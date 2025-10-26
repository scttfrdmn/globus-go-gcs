package collection

import "github.com/spf13/cobra"

// NewCollectionCmd creates the collection command with subcommands.
func NewCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection",
		Short: "Manage GCS collections",
		Long: `Commands for managing Globus Connect Server collections.

Collections define access points to data on an endpoint, including
both mapped collections (direct storage access) and guest collections
(shared subdirectories).`,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewShowCmd())

	return cmd
}
