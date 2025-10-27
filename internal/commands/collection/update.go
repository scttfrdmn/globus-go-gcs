package collection

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

// NewUpdateCmd creates the collection update command.
func NewUpdateCmd() *cobra.Command {
	var (
		profile                  string
		format                   string
		endpointFQDN             string
		displayName              string
		description              string
		public                   *bool
		disableAnonymousWrites   *bool
		contactEmail             string
		contactInfo              string
		infoLink                 string
		keywords                 string
		organization             string
		department               string
		userMessage              string
		userMessageLink          string
	)

	cmd := &cobra.Command{
		Use:   "update COLLECTION_ID",
		Short: "Update an existing collection",
		Long: `Update an existing collection's configuration.

Only the fields you specify will be updated. Other fields will remain unchanged.

Example:
  globus-connect-server collection update abc123 \
    --endpoint example.data.globus.org \
    --display-name "Updated Collection Name" \
    --description "New description"

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			collectionID := args[0]
			return runUpdate(cmd.Context(), profile, format, endpointFQDN, collectionID,
				displayName, description, public, disableAnonymousWrites,
				contactEmail, contactInfo, infoLink, keywords, organization,
				department, userMessage, userMessageLink, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&displayName, "display-name", "", "Display name for the collection")
	cmd.Flags().StringVar(&description, "description", "", "Description of the collection")

	// Boolean flags need special handling - use string and convert to *bool
	publicStr := cmd.Flags().String("public", "", "Make collection public (true/false)")
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
	cmd.Flags().StringVar(&contactEmail, "contact-email", "", "Contact email")
	cmd.Flags().StringVar(&contactInfo, "contact-info", "", "Contact information")
	cmd.Flags().StringVar(&infoLink, "info-link", "", "Information link URL")
	cmd.Flags().StringVar(&keywords, "keywords", "", "Comma-separated keywords")
	cmd.Flags().StringVar(&organization, "organization", "", "Organization name")
	cmd.Flags().StringVar(&department, "department", "", "Department name")
	cmd.Flags().StringVar(&userMessage, "user-message", "", "Message shown to users")
	cmd.Flags().StringVar(&userMessageLink, "user-message-link", "", "Link for user message")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runUpdate executes the collection update command.
func runUpdate(ctx context.Context, profile, formatStr, endpointFQDN, collectionID string,
	displayName, description string, public, disableAnonymousWrites *bool,
	contactEmail, contactInfo, infoLink, keywords, organization, department,
	userMessage, userMessageLink string,
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

	// Build collection object with only the fields that were specified
	collection := &gcs.Collection{}

	if displayName != "" {
		collection.DisplayName = displayName
	}
	if description != "" {
		collection.Description = description
	}
	if public != nil {
		collection.Public = *public
	}
	if disableAnonymousWrites != nil {
		collection.DisableAnonymousWrites = *disableAnonymousWrites
	}
	if contactEmail != "" {
		collection.ContactEmail = contactEmail
	}
	if contactInfo != "" {
		collection.ContactInfo = contactInfo
	}
	if infoLink != "" {
		collection.InfoLink = infoLink
	}
	if organization != "" {
		collection.Organization = organization
	}
	if department != "" {
		collection.Department = department
	}
	if userMessage != "" {
		collection.UserMessage = userMessage
	}
	if userMessageLink != "" {
		collection.UserMessageLink = userMessageLink
	}

	// Parse keywords if provided
	if keywords != "" {
		collection.Keywords = strings.Split(keywords, ",")
		// Trim whitespace from each keyword
		for i, kw := range collection.Keywords {
			collection.Keywords[i] = strings.TrimSpace(kw)
		}
	}

	// Update collection
	updated, err := gcsClient.UpdateCollection(ctx, collectionID, collection)
	if err != nil {
		return fmt.Errorf("update collection: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(updated)
	}

	// Text format
	if err := formatter.Println("Collection updated successfully!"); err != nil {
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
	if updated.CollectionType != "" {
		if err := formatter.PrintText("%-20s%s\n", "Type:", updated.CollectionType); err != nil {
			return err
		}
	}

	return nil
}
