// Package oidc provides commands for managing OpenID Connect server configuration.
package oidc

import "github.com/spf13/cobra"

// NewOIDCCmd creates the oidc command with subcommands.
func NewOIDCCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oidc",
		Short: "Manage OpenID Connect (OIDC) server configuration",
		Long: `Commands for managing OpenID Connect server configuration.

OIDC integration allows you to use your own identity provider for
authentication with Globus Connect Server.`,
	}

	// Add subcommands
	cmd.AddCommand(NewShowCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewRegisterCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewDeleteCmd())

	return cmd
}
