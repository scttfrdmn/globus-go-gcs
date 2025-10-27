// Package sharingpolicy provides commands for managing collection sharing policies.
package sharingpolicy

import "github.com/spf13/cobra"

// NewSharingPolicyCmd creates the sharing-policy command with subcommands.
func NewSharingPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sharing-policy",
		Short: "Manage collection sharing policies",
		Long: `Commands for managing collection sharing policies.

Sharing policies control who can share collections and what restrictions apply.`,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewShowCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewDeleteCmd())

	return cmd
}
