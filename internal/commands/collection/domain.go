package collection

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewDomainCmd creates the collection domain command with subcommands.
func NewDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "Manage collection custom domain",
		Long: `Manage custom domain configuration for a collection.

Custom domains allow you to use your own domain name (e.g., data.example.org)
for accessing a specific collection instead of the default Globus domain.
This requires providing SSL certificates and configuring DNS appropriately.

Available subcommands:
  setup  - Configure a custom domain
  show   - Display current domain configuration
  delete - Remove custom domain configuration`,
	}

	// Add subcommands
	cmd.AddCommand(NewDomainSetupCmd())
	cmd.AddCommand(NewDomainShowCmd())
	cmd.AddCommand(NewDomainDeleteCmd())

	return cmd
}

// NewDomainSetupCmd creates the collection domain setup command.
func NewDomainSetupCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		domain       string
		certificate  string
		privateKey   string
	)

	cmd := &cobra.Command{
		Use:   "setup COLLECTION_ID",
		Short: "Configure custom domain for collection",
		Long: `Configure a custom domain for a collection.

This command sets up a custom domain name for your collection, allowing users
to access it via your own domain (e.g., data.example.org) instead of the
default Globus domain.

Requirements:
- A valid domain name that you control
- SSL certificate for the domain
- Private key for the SSL certificate
- DNS configured to point to the collection

Example:
  globus-connect-server collection domain setup abc123 \
    --endpoint example.data.globus.org \
    --domain data.example.org \
    --certificate /path/to/cert.pem \
    --private-key /path/to/key.pem

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			collectionID := args[0]
			return runDomainSetup(cmd.Context(), profile, format, endpointFQDN, collectionID, domain, certificate, privateKey, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&domain, "domain", "", "Custom domain name")
	cmd.Flags().StringVar(&certificate, "certificate", "", "Path to SSL certificate file")
	cmd.Flags().StringVar(&privateKey, "private-key", "", "Path to private key file")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("domain")
	_ = cmd.MarkFlagRequired("certificate")
	_ = cmd.MarkFlagRequired("private-key")

	return cmd
}

// runDomainSetup executes the collection domain setup command.
func runDomainSetup(ctx context.Context, profile, formatStr, endpointFQDN, collectionID, domain, certificate, privateKey string, out interface{ Write([]byte) (int, error) }) error {
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

	// Build domain config
	domainConfig := &gcs.DomainConfig{
		Domain:      domain,
		Certificate: certificate,
		PrivateKey:  privateKey,
	}

	// Setup domain
	if err := gcsClient.SetupCollectionDomain(ctx, collectionID, domainConfig); err != nil {
		return fmt.Errorf("setup collection domain: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		result := map[string]string{
			"status":        "success",
			"collection_id": collectionID,
			"domain":        domain,
			"message":       "Collection domain configured successfully",
		}
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.PrintText("Collection domain configured successfully.\n"); err != nil {
		return err
	}
	if err := formatter.PrintText("Collection ID: %s\n", collectionID); err != nil {
		return err
	}
	if err := formatter.PrintText("Domain: %s\n", domain); err != nil {
		return err
	}

	return nil
}

// NewDomainShowCmd creates the collection domain show command.
func NewDomainShowCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "show COLLECTION_ID",
		Short: "Display collection domain configuration",
		Long: `Display the current custom domain configuration for a collection.

This command shows the configured custom domain name and verification status.

Example:
  globus-connect-server collection domain show abc123 \
    --endpoint example.data.globus.org

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			collectionID := args[0]
			return runDomainShow(cmd.Context(), profile, format, endpointFQDN, collectionID, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runDomainShow executes the collection domain show command.
func runDomainShow(ctx context.Context, profile, formatStr, endpointFQDN, collectionID string, out interface{ Write([]byte) (int, error) }) error {
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

	// Get domain configuration
	domainConfig, err := gcsClient.GetCollectionDomain(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("get collection domain: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(domainConfig)
	}

	// Text format
	if err := formatter.Println("Collection Domain Configuration"); err != nil {
		return err
	}
	if err := formatter.Println("================================="); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if err := formatter.PrintText("%-20s%s\n", "Collection ID:", collectionID); err != nil {
		return err
	}
	if domainConfig.Domain != "" {
		if err := formatter.PrintText("%-20s%s\n", "Domain:", domainConfig.Domain); err != nil {
			return err
		}
	}
	verifiedStr := "No"
	if domainConfig.Verified {
		verifiedStr = "Yes"
	}
	if err := formatter.PrintText("%-20s%s\n", "Verified:", verifiedStr); err != nil {
		return err
	}

	return nil
}

// NewDomainDeleteCmd creates the collection domain delete command.
func NewDomainDeleteCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		force        bool
	)

	cmd := &cobra.Command{
		Use:   "delete COLLECTION_ID",
		Short: "Remove collection custom domain configuration",
		Long: `Remove the custom domain configuration from a collection.

This command removes the custom domain configuration and reverts to using
the default Globus domain for the collection.

WARNING: This action cannot be undone. The collection will revert to using
the default Globus domain name.

Example:
  globus-connect-server collection domain delete abc123 \
    --endpoint example.data.globus.org \
    --force

Use --force to skip confirmation prompt.

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			collectionID := args[0]
			return runDomainDelete(cmd.Context(), profile, format, endpointFQDN, collectionID, force, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runDomainDelete executes the collection domain delete command.
func runDomainDelete(ctx context.Context, profile, formatStr, endpointFQDN, collectionID string, force bool, out interface{ Write([]byte) (int, error) }) error {
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

	// Confirmation prompt (unless --force)
	if !force {
		if err := formatter.PrintText("WARNING: This will remove the custom domain configuration for collection %s.\n", collectionID); err != nil {
			return err
		}
		if err := formatter.Println("The collection will revert to using the default Globus domain."); err != nil {
			return err
		}
		if err := formatter.Println(); err != nil {
			return err
		}
		if err := formatter.PrintText("To proceed, use --force flag.\n"); err != nil {
			return err
		}
		return fmt.Errorf("domain delete cancelled (use --force to proceed)")
	}

	// Create GCS client
	gcsClient, err := gcs.NewClient(
		endpointFQDN,
		gcs.WithAccessToken(token.AccessToken),
	)
	if err != nil {
		return fmt.Errorf("create GCS client: %w", err)
	}

	// Delete domain
	if err := gcsClient.DeleteCollectionDomain(ctx, collectionID); err != nil {
		return fmt.Errorf("delete collection domain: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		result := map[string]string{
			"status":        "success",
			"collection_id": collectionID,
			"message":       "Collection domain configuration removed successfully",
		}
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.PrintText("Collection domain configuration removed successfully.\n"); err != nil {
		return err
	}
	if err := formatter.PrintText("Collection ID: %s\n", collectionID); err != nil {
		return err
	}

	return nil
}
