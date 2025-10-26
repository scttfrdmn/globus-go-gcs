package auth

import (
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/spf13/cobra"
)

// NewLogoutCmd creates the logout command.
func NewLogoutCmd() *cobra.Command {
	var profile string

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Remove stored authentication tokens",
		Long: `Remove stored authentication tokens for the specified profile.

This command deletes the locally stored tokens, effectively logging you out.
You will need to login again to use authenticated commands.

The token file is removed from: ~/.globus-connect-server/tokens/<profile>.json`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runLogout(profile)
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")

	return cmd
}

// runLogout executes the logout flow.
func runLogout(profile string) error {
	// Check if token exists
	_, err := auth.LoadToken(profile)
	if err != nil {
		return fmt.Errorf("no active session for profile %q", profile)
	}

	// Delete token
	if err := auth.DeleteToken(profile); err != nil {
		return fmt.Errorf("delete token: %w", err)
	}

	fmt.Printf("✓ Logged out from profile: %s\n", profile)
	return nil
}
