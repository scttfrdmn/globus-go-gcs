package storagegateway

import "github.com/spf13/cobra"

// NewStorageGatewayCmd creates the storage-gateway command with subcommands.
func NewStorageGatewayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage-gateway",
		Short: "Manage GCS storage gateways",
		Long: `Commands for managing Globus Connect Server storage gateways.

Storage gateways connect GCS to storage backends such as POSIX filesystems,
Amazon S3, Azure Blob Storage, Google Cloud Storage, and others.`,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewShowCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewDeleteCmd())

	return cmd
}
