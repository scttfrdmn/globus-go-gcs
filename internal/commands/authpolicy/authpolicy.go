// Package authpolicy provides commands for managing authentication policies.
package authpolicy

import "github.com/spf13/cobra"

// NewAuthPolicyCmd creates the auth-policy command with subcommands.
func NewAuthPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth-policy",
		Short: "Manage authentication policies",
		Long: `Commands for managing authentication policies.

Authentication policies define security requirements and restrictions
for accessing the endpoint, such as MFA requirements, high assurance,
and domain-based access controls.`,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewShowCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewDeleteCmd())

	return cmd
}
