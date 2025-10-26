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

// NewShowCmd creates the collection show command.
func NewShowCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "show COLLECTION_ID",
		Short: "Display collection details",
		Long: `Display detailed information about a specific collection.

This command retrieves and displays the complete configuration of a collection
including display name, type, storage gateway, paths, and policies.

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			collectionID := args[0]
			return runShow(cmd.Context(), profile, format, endpointFQDN, collectionID, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runShow executes the collection show command.
func runShow(ctx context.Context, profile, formatStr, endpointFQDN, collectionID string, out interface{ Write([]byte) (int, error) }) error {
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

	// Get collection
	collection, err := gcsClient.GetCollection(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("get collection: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(collection)
	}

	// Text format
	return formatCollectionText(formatter, collection)
}

// formatCollectionText formats the collection in text format.
func formatCollectionText(formatter *output.Formatter, collection *gcs.Collection) error {
	if err := formatter.Println("Collection Details:"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	// Print basic fields
	if err := printCollectionBasicFields(formatter, collection); err != nil {
		return err
	}

	// Print optional fields
	if err := printCollectionOptionalFields(formatter, collection); err != nil {
		return err
	}

	// Print keywords
	if err := printCollectionKeywords(formatter, collection); err != nil {
		return err
	}

	// Print policies
	if err := printCollectionPolicies(formatter, collection); err != nil {
		return err
	}

	return nil
}

// printCollectionBasicFields prints the basic collection fields.
func printCollectionBasicFields(formatter *output.Formatter, collection *gcs.Collection) error {
	printField := func(label, value string) error {
		if value != "" {
			return formatter.PrintText("%-25s%s\n", label+":", value)
		}
		return nil
	}

	if err := printField("Display Name", collection.DisplayName); err != nil {
		return err
	}
	if err := printField("ID", collection.ID); err != nil {
		return err
	}
	if err := printField("Type", collection.CollectionType); err != nil {
		return err
	}
	if err := printField("Storage Gateway ID", collection.StorageGatewayID); err != nil {
		return err
	}
	if err := printField("Base Path", collection.CollectionBaseFolder); err != nil {
		return err
	}

	// Boolean fields
	if err := formatter.PrintText("%-25s%t\n", "Public:", collection.Public); err != nil {
		return err
	}
	return formatter.PrintText("%-25s%t\n", "Disable Anonymous Writes:", collection.DisableAnonymousWrites)
}

// printCollectionOptionalFields prints optional collection fields.
func printCollectionOptionalFields(formatter *output.Formatter, collection *gcs.Collection) error {
	printField := func(label, value string) error {
		if value != "" {
			return formatter.PrintText("%-25s%s\n", label+":", value)
		}
		return nil
	}

	fields := []struct{ label, value string }{
		{"Organization", collection.Organization},
		{"Department", collection.Department},
		{"Description", collection.Description},
		{"Contact Email", collection.ContactEmail},
		{"Contact Info", collection.ContactInfo},
		{"Info Link", collection.InfoLink},
		{"User Message", collection.UserMessage},
		{"User Message Link", collection.UserMessageLink},
		{"Identity ID", collection.IdentityID},
	}

	for _, f := range fields {
		if err := printField(f.label, f.value); err != nil {
			return err
		}
	}

	return nil
}

// printCollectionKeywords prints collection keywords.
func printCollectionKeywords(formatter *output.Formatter, collection *gcs.Collection) error {
	if len(collection.Keywords) == 0 {
		return nil
	}

	if err := formatter.Println(); err != nil {
		return err
	}
	if err := formatter.Println("Keywords:"); err != nil {
		return err
	}
	for _, keyword := range collection.Keywords {
		if err := formatter.PrintText("  - %s\n", keyword); err != nil {
			return err
		}
	}
	return nil
}

// printCollectionPolicies prints collection policies.
func printCollectionPolicies(formatter *output.Formatter, collection *gcs.Collection) error {
	if collection.Policies == nil {
		return nil
	}

	if err := formatter.Println(); err != nil {
		return err
	}
	if err := formatter.Println("Policies:"); err != nil {
		return err
	}

	if collection.Policies.AuthenticationTimeoutMins > 0 {
		if err := formatter.PrintText("  Authentication Timeout: %d minutes\n", collection.Policies.AuthenticationTimeoutMins); err != nil {
			return err
		}
	}
	if collection.Policies.SharingRestrict != "" {
		if err := formatter.PrintText("  Sharing Restrict:       %s\n", collection.Policies.SharingRestrict); err != nil {
			return err
		}
	}

	return nil
}
