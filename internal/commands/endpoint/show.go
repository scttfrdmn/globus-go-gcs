// Package endpoint provides endpoint management commands for the GCS CLI.
package endpoint

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewShowCmd creates the endpoint show command.
func NewShowCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Display endpoint configuration",
		Long: `Display the configuration of a Globus Connect Server endpoint.

This command retrieves and displays the current configuration of a GCS endpoint
including display name, organization, contact information, and other settings.

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runShow(cmd.Context(), profile, format, endpointFQDN, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runShow executes the endpoint show command.
func runShow(ctx context.Context, profile, formatStr, endpointFQDN string, out interface{ Write([]byte) (int, error) }) error {
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

	// Get endpoint configuration
	endpoint, err := gcsClient.GetEndpoint(ctx)
	if err != nil {
		return fmt.Errorf("get endpoint: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(endpoint)
	}

	// Text format
	return formatEndpointText(formatter, endpoint)
}

// formatEndpointText formats the endpoint in text format.
func formatEndpointText(formatter *output.Formatter, endpoint *gcs.Endpoint) error {
	if err := formatter.Println("Endpoint Configuration:"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	// Helper function to print optional fields
	printField := func(label, value string) error {
		if value != "" {
			return formatter.PrintText("%-20s%s\n", label+":", value)
		}
		return nil
	}

	// Print fields
	if err := printField("Display Name", endpoint.DisplayName); err != nil {
		return err
	}
	if err := printField("ID", endpoint.ID); err != nil {
		return err
	}
	if err := printField("Organization", endpoint.Organization); err != nil {
		return err
	}
	if err := printField("Department", endpoint.Department); err != nil {
		return err
	}
	if err := printField("Description", endpoint.Description); err != nil {
		return err
	}
	if err := printField("Contact Email", endpoint.ContactEmail); err != nil {
		return err
	}
	if err := printField("Contact Info", endpoint.ContactInfo); err != nil {
		return err
	}
	if err := printField("Info Link", endpoint.InfoLink); err != nil {
		return err
	}

	// Boolean field (always print)
	if err := formatter.PrintText("%-20s%t\n", "Public:", endpoint.Public); err != nil {
		return err
	}

	if err := printField("Default Directory", endpoint.DefaultDirectory); err != nil {
		return err
	}
	if err := printField("Network Use", endpoint.NetworkUse); err != nil {
		return err
	}

	// Numeric fields (only if > 0)
	if endpoint.MaxConcurrency > 0 {
		if err := formatter.PrintText("%-20s%d\n", "Max Concurrency:", endpoint.MaxConcurrency); err != nil {
			return err
		}
	}
	if endpoint.PreferredConcurrency > 0 {
		if err := formatter.PrintText("%-20s%d\n", "Preferred Concurrency:", endpoint.PreferredConcurrency); err != nil {
			return err
		}
	}

	// Keywords
	if len(endpoint.Keywords) > 0 {
		if err := formatter.Println(); err != nil {
			return err
		}
		if err := formatter.Println("Keywords:"); err != nil {
			return err
		}
		for _, keyword := range endpoint.Keywords {
			if err := formatter.PrintText("  - %s\n", keyword); err != nil {
				return err
			}
		}
	}

	return nil
}
