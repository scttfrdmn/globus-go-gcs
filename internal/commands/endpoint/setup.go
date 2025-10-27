package endpoint

import (
	"context"
	"fmt"
	"strings"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewSetupCmd creates the endpoint setup command.
func NewSetupCmd() *cobra.Command {
	var (
		profile              string
		format               string
		endpointFQDN         string
		displayName          string
		organization         string
		department           string
		description          string
		contactEmail         string
		contactInfo          string
		infoLink             string
		public               bool
		keywords             string
	)

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Create and initialize a new GCS endpoint",
		Long: `Create and initialize a new Globus Connect Server endpoint.

This command sets up a new endpoint with the specified configuration.
The endpoint will be registered with Globus and ready for use.

Example:
  globus-connect-server endpoint setup \
    --endpoint example.data.globus.org \
    --display-name "My Research Endpoint" \
    --organization "Example University" \
    --contact-email "support@example.edu"

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSetup(cmd.Context(), profile, format, endpointFQDN,
				displayName, organization, department, description,
				contactEmail, contactInfo, infoLink, public, keywords, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&displayName, "display-name", "", "Display name for the endpoint")
	cmd.Flags().StringVar(&organization, "organization", "", "Organization name")
	cmd.Flags().StringVar(&department, "department", "", "Department name")
	cmd.Flags().StringVar(&description, "description", "", "Description of the endpoint")
	cmd.Flags().StringVar(&contactEmail, "contact-email", "", "Contact email")
	cmd.Flags().StringVar(&contactInfo, "contact-info", "", "Contact information")
	cmd.Flags().StringVar(&infoLink, "info-link", "", "Information link URL")
	cmd.Flags().BoolVar(&public, "public", false, "Make endpoint public")
	cmd.Flags().StringVar(&keywords, "keywords", "", "Comma-separated keywords")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("display-name")

	return cmd
}

// runSetup executes the endpoint setup command.
func runSetup(ctx context.Context, profile, formatStr, endpointFQDN string,
	displayName, organization, department, description,
	contactEmail, contactInfo, infoLink string, public bool, keywords string,
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

	// Build endpoint configuration
	endpoint := &gcs.Endpoint{
		DisplayName:  displayName,
		Organization: organization,
		Department:   department,
		Description:  description,
		ContactEmail: contactEmail,
		ContactInfo:  contactInfo,
		InfoLink:     infoLink,
		Public:       public,
	}

	// Parse keywords if provided
	if keywords != "" {
		endpoint.Keywords = strings.Split(keywords, ",")
		// Trim whitespace from each keyword
		for i, kw := range endpoint.Keywords {
			endpoint.Keywords[i] = strings.TrimSpace(kw)
		}
	}

	// Setup endpoint
	created, err := gcsClient.SetupEndpoint(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("setup endpoint: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(created)
	}

	// Text format
	if err := formatter.Println("Endpoint setup completed successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if created.ID != "" {
		if err := formatter.PrintText("%-20s%s\n", "Endpoint ID:", created.ID); err != nil {
			return err
		}
	}
	if created.DisplayName != "" {
		if err := formatter.PrintText("%-20s%s\n", "Display Name:", created.DisplayName); err != nil {
			return err
		}
	}
	if created.Organization != "" {
		if err := formatter.PrintText("%-20s%s\n", "Organization:", created.Organization); err != nil {
			return err
		}
	}

	return nil
}
