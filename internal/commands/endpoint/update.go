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

// NewUpdateCmd creates the endpoint update command.
func NewUpdateCmd() *cobra.Command {
	var (
		profile                  string
		format                   string
		endpointFQDN             string
		displayName              string
		organization             string
		department               string
		description              string
		contactEmail             string
		contactInfo              string
		infoLink                 string
		public                   *bool
		defaultDirectory         string
		networkUse               string
		maxConcurrency           int
		preferredConcurrency     int
		disableAnonymousWrites   *bool
		keywords                 string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update endpoint configuration",
		Long: `Update the configuration of a Globus Connect Server endpoint.

Only the fields you specify will be updated. Other fields will remain unchanged.

Example:
  globus-connect-server endpoint update \
    --endpoint example.data.globus.org \
    --display-name "My Updated Endpoint" \
    --organization "Example University" \
    --contact-email "support@example.edu"

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runUpdate(cmd.Context(), profile, format, endpointFQDN,
				displayName, organization, department, description,
				contactEmail, contactInfo, infoLink, public, defaultDirectory,
				networkUse, maxConcurrency, preferredConcurrency,
				disableAnonymousWrites, keywords, cmd.OutOrStdout())
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

	// Boolean flags need special handling
	publicStr := cmd.Flags().String("public", "", "Make endpoint public (true/false)")
	cmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		if *publicStr != "" {
			val := *publicStr == "true"
			public = &val
		}
		disableAnonStr, _ := cmd.Flags().GetString("disable-anonymous-writes")
		if disableAnonStr != "" {
			val := disableAnonStr == "true"
			disableAnonymousWrites = &val
		}
		return nil
	}

	cmd.Flags().String("disable-anonymous-writes", "", "Disable anonymous writes (true/false)")
	cmd.Flags().StringVar(&defaultDirectory, "default-directory", "", "Default directory for transfers")
	cmd.Flags().StringVar(&networkUse, "network-use", "", "Network use policy (normal, minimal, aggressive, custom)")
	cmd.Flags().IntVar(&maxConcurrency, "max-concurrency", 0, "Maximum transfer concurrency")
	cmd.Flags().IntVar(&preferredConcurrency, "preferred-concurrency", 0, "Preferred transfer concurrency")
	cmd.Flags().StringVar(&keywords, "keywords", "", "Comma-separated keywords")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runUpdate executes the endpoint update command.
func runUpdate(ctx context.Context, profile, formatStr, endpointFQDN string,
	displayName, organization, department, description,
	contactEmail, contactInfo, infoLink string, public *bool, defaultDirectory,
	networkUse string, maxConcurrency, preferredConcurrency int,
	disableAnonymousWrites *bool, keywords string,
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

	// Build endpoint object
	endpoint := buildEndpointUpdate(displayName, organization, department, description,
		contactEmail, contactInfo, infoLink, public, defaultDirectory,
		networkUse, maxConcurrency, preferredConcurrency,
		disableAnonymousWrites, keywords)

	// Update endpoint
	updated, err := gcsClient.UpdateEndpoint(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("update endpoint: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(updated)
	}

	// Text format
	return formatUpdateSuccess(formatter, updated)
}

// buildEndpointUpdate constructs an endpoint object with specified fields.
func buildEndpointUpdate(displayName, organization, department, description,
	contactEmail, contactInfo, infoLink string, public *bool, defaultDirectory,
	networkUse string, maxConcurrency, preferredConcurrency int,
	disableAnonymousWrites *bool, keywords string) *gcs.Endpoint {

	endpoint := &gcs.Endpoint{}

	if displayName != "" {
		endpoint.DisplayName = displayName
	}
	if organization != "" {
		endpoint.Organization = organization
	}
	if department != "" {
		endpoint.Department = department
	}
	if description != "" {
		endpoint.Description = description
	}
	if contactEmail != "" {
		endpoint.ContactEmail = contactEmail
	}
	if contactInfo != "" {
		endpoint.ContactInfo = contactInfo
	}
	if infoLink != "" {
		endpoint.InfoLink = infoLink
	}
	if public != nil {
		endpoint.Public = *public
	}
	if defaultDirectory != "" {
		endpoint.DefaultDirectory = defaultDirectory
	}
	if networkUse != "" {
		endpoint.NetworkUse = networkUse
	}
	if maxConcurrency > 0 {
		endpoint.MaxConcurrency = maxConcurrency
	}
	if preferredConcurrency > 0 {
		endpoint.PreferredConcurrency = preferredConcurrency
	}
	if disableAnonymousWrites != nil {
		endpoint.DisableAnonymousWrites = *disableAnonymousWrites
	}

	// Parse keywords if provided
	if keywords != "" {
		endpoint.Keywords = strings.Split(keywords, ",")
		// Trim whitespace from each keyword
		for i, kw := range endpoint.Keywords {
			endpoint.Keywords[i] = strings.TrimSpace(kw)
		}
	}

	return endpoint
}

// formatUpdateSuccess formats the successful update response.
func formatUpdateSuccess(formatter *output.Formatter, updated *gcs.Endpoint) error {
	if err := formatter.Println("Endpoint updated successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if updated.ID != "" {
		if err := formatter.PrintText("%-20s%s\n", "ID:", updated.ID); err != nil {
			return err
		}
	}
	if updated.DisplayName != "" {
		if err := formatter.PrintText("%-20s%s\n", "Display Name:", updated.DisplayName); err != nil {
			return err
		}
	}
	if updated.Organization != "" {
		if err := formatter.PrintText("%-20s%s\n", "Organization:", updated.Organization); err != nil {
			return err
		}
	}
	if updated.Department != "" {
		if err := formatter.PrintText("%-20s%s\n", "Department:", updated.Department); err != nil {
			return err
		}
	}

	return nil
}
