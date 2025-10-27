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
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewDeleteCmd())
	cmd.AddCommand(NewCheckCmd())
	cmd.AddCommand(NewBatchDeleteCmd())
	cmd.AddCommand(NewSetOwnerCmd())
	cmd.AddCommand(NewSetOwnerStringCmd())
	cmd.AddCommand(NewResetOwnerStringCmd())
	cmd.AddCommand(NewSetSubscriptionAdminVerifiedCmd())

	return cmd
}
