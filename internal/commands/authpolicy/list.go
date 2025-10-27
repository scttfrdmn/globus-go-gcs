package authpolicy

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewListCmd creates the auth-policy list command.
func NewListCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all authentication policies",
		Long: `List all authentication policies configured for the endpoint.

This command displays all authentication policies, including their IDs,
names, and security requirements.

Example:
  globus-connect-server auth-policy list \
    --endpoint example.data.globus.org

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runList(cmd.Context(), profile, format, endpointFQDN, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runList executes the auth-policy list command.
func runList(ctx context.Context, profile, formatStr, endpointFQDN string, out interface{ Write([]byte) (int, error) }) error {
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

	// List auth policies
	list, err := gcsClient.ListAuthPolicies(ctx)
	if err != nil {
		return fmt.Errorf("list auth policies: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(list.Data)
	}

	// Text format
	if len(list.Data) == 0 {
		return formatter.Println("No authentication policies found.")
	}

	for i, policy := range list.Data {
		if i > 0 {
			if err := formatter.Println(); err != nil {
				return err
			}
		}

		if policy.ID != "" {
			if err := formatter.PrintText("%-20s%s\n", "ID:", policy.ID); err != nil {
				return err
			}
		}
		if policy.Name != "" {
			if err := formatter.PrintText("%-20s%s\n", "Name:", policy.Name); err != nil {
				return err
			}
		}
		if policy.Description != "" {
			if err := formatter.PrintText("%-20s%s\n", "Description:", policy.Description); err != nil {
				return err
			}
		}
		if err := formatter.PrintText("%-20s%v\n", "Require MFA:", policy.RequireMFA); err != nil {
			return err
		}
		if err := formatter.PrintText("%-20s%v\n", "High Assurance:", policy.RequireHighAssurance); err != nil {
			return err
		}
	}

	return nil
}
