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

// NewCreateCmd creates the role create command.
func NewCreateCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		collection   string
		principal    string
		role         string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new role assignment",
		Long: `Create a new role assignment for a principal on a collection or endpoint.

A role grants a principal (user or group identity) specific permissions.

Common role types:
  - administrator: Full endpoint management access
  - access_manager: Can manage access rules for a collection
  - owner: Collection owner with full permissions
  - activity_manager: Can monitor and manage endpoint activity

Example:
  globus-connect-server role create \
    --endpoint example.data.globus.org \
    --collection abc123 \
    --principal "user@globusid.org" \
    --role owner

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCreate(cmd.Context(), profile, format, endpointFQDN,
				collection, principal, role, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVarP(&collection, "collection", "c", "", "Collection ID")
	cmd.Flags().StringVar(&principal, "principal", "", "Principal identity (user or group)")
	cmd.Flags().StringVar(&role, "role", "", "Role type (administrator, owner, access_manager, etc.)")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("collection")
	_ = cmd.MarkFlagRequired("principal")
	_ = cmd.MarkFlagRequired("role")

	return cmd
}

// runCreate executes the role create command.
func runCreate(ctx context.Context, profile, formatStr, endpointFQDN string,
	collection, principal, role string,
	out interface{ Write([]byte) (int, error) }) error {

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

	// Build role object
	roleObj := &gcs.Role{
		Collection: collection,
		Principal:  principal,
		Role:       role,
	}

	// Create role
	created, err := gcsClient.CreateRole(ctx, roleObj)
	if err != nil {
		return fmt.Errorf("create role: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(created)
	}

	// Text format
	if err := formatter.Println("Role created successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if created.ID != "" {
		if err := formatter.PrintText("%-20s%s\n", "ID:", created.ID); err != nil {
			return err
		}
	}
	if created.Collection != "" {
		if err := formatter.PrintText("%-20s%s\n", "Collection:", created.Collection); err != nil {
			return err
		}
	}
	if created.Principal != "" {
		if err := formatter.PrintText("%-20s%s\n", "Principal:", created.Principal); err != nil {
			return err
		}
	}
	if created.Role != "" {
		if err := formatter.PrintText("%-20s%s\n", "Role:", created.Role); err != nil {
			return err
		}
	}

	return nil
}
