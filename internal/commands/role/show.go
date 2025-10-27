package role

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewShowCmd creates the role show command.
func NewShowCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "show ROLE_ID",
		Short: "Display role details",
		Long: `Display detailed information about a specific role assignment.

This command retrieves and displays the complete configuration of a role
assignment including the collection, principal (user/group), and role type.

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			roleID := args[0]
			return runShow(cmd.Context(), profile, format, endpointFQDN, roleID, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runShow executes the role show command.
func runShow(ctx context.Context, profile, formatStr, endpointFQDN, roleID string, out interface{ Write([]byte) (int, error) }) error {
	// Load token
	token, err := auth.LoadToken(profile)
	if err != nil {
		return fmt.Errorf("not logged in: %w (use 'login' command first)", err)
	}

	// Check if token is valid
	if !token.IsValid() {
		return fmt.Errorf("token expired, please login again")
	}

	// Create output formatter
	formatter := output.NewFormatter(output.Format(formatStr), out)

	// Create GCS client
	gcsClient, err := gcs.NewClient(
		endpointFQDN,
		gcs.WithAccessToken(token.AccessToken),
	)
	if err != nil {
		return fmt.Errorf("create GCS client: %w", err)
	}

	// Get role
	role, err := gcsClient.GetRole(ctx, roleID)
	if err != nil {
		return fmt.Errorf("get role: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(role)
	}

	// Text format
	if err := formatter.Println("Role Details:"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if role.ID != "" {
		if err := formatter.PrintText("%-15s%s\n", "ID:", role.ID); err != nil {
			return err
		}
	}
	if role.Collection != "" {
		if err := formatter.PrintText("%-15s%s\n", "Collection:", role.Collection); err != nil {
			return err
		}
	}
	if role.Principal != "" {
		if err := formatter.PrintText("%-15s%s\n", "Principal:", role.Principal); err != nil {
			return err
		}
	}
	if role.Role != "" {
		if err := formatter.PrintText("%-15s%s\n", "Role:", role.Role); err != nil {
			return err
		}
	}

	return nil
}
