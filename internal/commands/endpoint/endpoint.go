package endpoint

import "github.com/spf13/cobra"

// NewEndpointCmd creates the endpoint command with subcommands.
func NewEndpointCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "endpoint",
		Short: "Manage GCS endpoints",
		Long: `Commands for managing Globus Connect Server endpoints.

Endpoints are the top-level configuration for a GCS installation,
defining the service's identity and general settings.`,
	}

	// Add subcommands
	cmd.AddCommand(NewShowCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewSetupCmd())
	cmd.AddCommand(NewCleanupCmd())
	cmd.AddCommand(NewKeyConvertCmd())
	cmd.AddCommand(NewSetOwnerCmd())
	cmd.AddCommand(NewSetOwnerStringCmd())
	cmd.AddCommand(NewResetOwnerStringCmd())
	cmd.AddCommand(NewSetSubscriptionIDCmd())
	cmd.AddCommand(NewDomainCmd())
	cmd.AddCommand(NewUpgradeCmd())

	return cmd
}
