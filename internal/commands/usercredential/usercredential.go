// Package usercredential provides commands for managing user storage credentials.
package usercredential

import "github.com/spf13/cobra"

// NewUserCredentialCmd creates the user-credential command with subcommands.
func NewUserCredentialCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-credential",
		Short: "Manage user storage credentials",
		Long: `Commands for managing user storage credentials.

User credentials provide authentication for storage backends including
ActiveScale, OAuth2, and S3-compatible storage.`,
	}

	// Add subcommands
	cmd.AddCommand(NewActivescaleCreateCmd())
	cmd.AddCommand(NewOAuthCreateCmd())
	cmd.AddCommand(NewS3CreateCmd())
	cmd.AddCommand(NewS3KeysAddCmd())
	cmd.AddCommand(NewS3KeysUpdateCmd())
	cmd.AddCommand(NewS3KeysDeleteCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewShowCmd())
	cmd.AddCommand(NewDeleteCmd())

	return cmd
}
