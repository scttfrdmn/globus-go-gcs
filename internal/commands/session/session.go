// Package session provides commands for managing CLI authentication sessions.
package session

import "github.com/spf13/cobra"

// NewSessionCmd creates the session command with subcommands.
func NewSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage CLI authentication session",
		Long: `Commands for managing the current CLI authentication session.

Sessions manage authentication state, timeout settings, and user consents
for CLI operations.`,
	}

	// Add subcommands
	cmd.AddCommand(NewShowCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewConsentCmd())

	return cmd
}
