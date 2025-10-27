package role

import "github.com/spf13/cobra"

// NewRoleCmd creates the role command with subcommands.
func NewRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "Manage GCS roles",
		Long: `Commands for managing Globus Connect Server roles.

Roles control access permissions for collections and endpoint management.
Each role assignment grants a principal (user or group) specific permissions.

Common role types include:
  - administrator: Full endpoint management access
  - access_manager: Can manage access rules for a collection
  - owner: Collection owner with full permissions
  - activity_manager: Can monitor and manage endpoint activity`,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewShowCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewDeleteCmd())

	return cmd
}
