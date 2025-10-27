package sharingpolicy

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewListCmd creates the sharing-policy list command.
func NewListCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all sharing policies",
		Long: `List all sharing policies configured on the endpoint.

Sharing policies control who can create shares on collections and what
restrictions apply.

Example:
  globus-connect-server sharing-policy list \
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

// runList executes the sharing-policy list command.
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

	// List sharing policies
	policies, err := gcsClient.ListSharingPolicies(ctx)
	if err != nil {
		return fmt.Errorf("list sharing policies: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(policies)
	}

	// Text format
	if len(policies.Data) == 0 {
		return formatter.Println("No sharing policies found")
	}

	if err := formatter.Println("Sharing Policies"); err != nil {
		return err
	}
	if err := formatter.Println("================"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	for _, policy := range policies.Data {
		if policy.ID != "" {
			if err := formatter.PrintText("ID:          %s\n", policy.ID); err != nil {
				return err
			}
		}
		if policy.Name != "" {
			if err := formatter.PrintText("Name:        %s\n", policy.Name); err != nil {
				return err
			}
		}
		if policy.CollectionID != "" {
			if err := formatter.PrintText("Collection:  %s\n", policy.CollectionID); err != nil {
				return err
			}
		}
		if policy.Description != "" {
			if err := formatter.PrintText("Description: %s\n", policy.Description); err != nil {
				return err
			}
		}
		if err := formatter.Println(); err != nil {
			return err
		}
	}

	return nil
}
