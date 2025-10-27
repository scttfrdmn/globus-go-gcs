// Package role provides role management commands for the GCS CLI.
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

// NewListCmd creates the role list command.
func NewListCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		collection   string
		principal    string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List roles on an endpoint",
		Long: `List all role assignments on a Globus Connect Server endpoint.

Roles control access to collections and endpoint management. Each role
assignment grants a principal (user or group) specific permissions on
a collection or the endpoint itself.

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runList(cmd.Context(), profile, format, endpointFQDN, collection, principal, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVarP(&collection, "collection", "c", "", "Filter roles by collection ID")
	cmd.Flags().StringVar(&principal, "principal", "", "Filter roles by principal (identity)")
	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runList executes the role list command.
func runList(ctx context.Context, profile, formatStr, endpointFQDN, collection, principal string, out interface{ Write([]byte) (int, error) }) error {
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

	// Build list options
	opts := &gcs.ListRolesOptions{}
	if collection != "" {
		opts.Collection = collection
	}
	if principal != "" {
		opts.Principal = principal
	}

	// Get roles
	list, err := gcsClient.ListRoles(ctx, opts)
	if err != nil {
		return fmt.Errorf("list roles: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(list)
	}

	// Text format
	if len(list.Data) == 0 {
		if err := formatter.Println("No roles found."); err != nil {
			return err
		}
		return nil
	}

	if err := formatter.PrintText("Roles (%d):\n\n", len(list.Data)); err != nil {
		return err
	}

	for i, role := range list.Data {
		if i > 0 {
			if err := formatter.Println(); err != nil {
				return err
			}
		}

		if err := formatter.PrintText("  ID:           %s\n", role.ID); err != nil {
			return err
		}
		if role.Collection != "" {
			if err := formatter.PrintText("  Collection:   %s\n", role.Collection); err != nil {
				return err
			}
		}
		if role.Principal != "" {
			if err := formatter.PrintText("  Principal:    %s\n", role.Principal); err != nil {
				return err
			}
		}
		if role.Role != "" {
			if err := formatter.PrintText("  Role:         %s\n", role.Role); err != nil {
				return err
			}
		}
	}

	return nil
}
