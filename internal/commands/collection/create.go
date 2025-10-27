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

// NewCreateCmd creates the collection create command.
func NewCreateCmd() *cobra.Command {
	var (
		profile                  string
		format                   string
		endpointFQDN             string
		displayName              string
		storageGatewayID         string
		collectionBaseFolder     string
		collectionType           string
		description              string
		public                   bool
		disableAnonymousWrites   bool
		contactEmail             string
		contactInfo              string
		infoLink                 string
		keywords                 string
		organization             string
		department               string
		userMessage              string
		userMessageLink          string
		identityID               string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new collection",
		Long: `Create a new collection on the endpoint.

A collection provides access to a storage location through a storage gateway.
You must specify the display name, storage gateway ID, and base path.

Example:
  globus-connect-server collection create \
    --endpoint example.data.globus.org \
    --display-name "My Collection" \
    --storage-gateway-id abc123 \
    --collection-base-path /data/shared

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCreate(cmd.Context(), profile, format, endpointFQDN,
				displayName, storageGatewayID, collectionBaseFolder, collectionType,
				description, public, disableAnonymousWrites, contactEmail,
				contactInfo, infoLink, keywords, organization, department,
				userMessage, userMessageLink, identityID, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&displayName, "display-name", "", "Display name for the collection")
	cmd.Flags().StringVar(&storageGatewayID, "storage-gateway-id", "", "Storage gateway ID")
	cmd.Flags().StringVar(&collectionBaseFolder, "collection-base-path", "", "Base path for the collection")
	cmd.Flags().StringVar(&collectionType, "collection-type", "mapped", "Collection type (mapped, guest)")
	cmd.Flags().StringVar(&description, "description", "", "Description of the collection")
	cmd.Flags().BoolVar(&public, "public", false, "Make collection public")
	cmd.Flags().BoolVar(&disableAnonymousWrites, "disable-anonymous-writes", false, "Disable anonymous writes")
	cmd.Flags().StringVar(&contactEmail, "contact-email", "", "Contact email")
	cmd.Flags().StringVar(&contactInfo, "contact-info", "", "Contact information")
	cmd.Flags().StringVar(&infoLink, "info-link", "", "Information link URL")
	cmd.Flags().StringVar(&keywords, "keywords", "", "Comma-separated keywords")
	cmd.Flags().StringVar(&organization, "organization", "", "Organization name")
	cmd.Flags().StringVar(&department, "department", "", "Department name")
	cmd.Flags().StringVar(&userMessage, "user-message", "", "Message shown to users")
	cmd.Flags().StringVar(&userMessageLink, "user-message-link", "", "Link for user message")
	cmd.Flags().StringVar(&identityID, "identity-id", "", "Identity ID")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("display-name")
	_ = cmd.MarkFlagRequired("storage-gateway-id")
	_ = cmd.MarkFlagRequired("collection-base-path")

	return cmd
}

// runCreate executes the collection create command.
func runCreate(ctx context.Context, profile, formatStr, endpointFQDN string,
	displayName, storageGatewayID, collectionBaseFolder, collectionType,
	description string, public, disableAnonymousWrites bool,
	contactEmail, contactInfo, infoLink, keywords, organization, department,
	userMessage, userMessageLink, identityID string,
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

	// Build collection object
	collection := &gcs.Collection{
		DisplayName:            displayName,
		StorageGatewayID:       storageGatewayID,
		CollectionBaseFolder:   collectionBaseFolder,
		CollectionType:         collectionType,
		Description:            description,
		Public:                 public,
		DisableAnonymousWrites: disableAnonymousWrites,
		ContactEmail:           contactEmail,
		ContactInfo:            contactInfo,
		InfoLink:               infoLink,
		Organization:           organization,
		Department:             department,
		UserMessage:            userMessage,
		UserMessageLink:        userMessageLink,
		IdentityID:             identityID,
	}

	// Parse keywords if provided
	if keywords != "" {
		collection.Keywords = strings.Split(keywords, ",")
		// Trim whitespace from each keyword
		for i, kw := range collection.Keywords {
			collection.Keywords[i] = strings.TrimSpace(kw)
		}
	}

	// Create collection
	created, err := gcsClient.CreateCollection(ctx, collection)
	if err != nil {
		return fmt.Errorf("create collection: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(created)
	}

	// Text format
	if err := formatter.Println("Collection created successfully!"); err != nil {
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
	if created.DisplayName != "" {
		if err := formatter.PrintText("%-20s%s\n", "Display Name:", created.DisplayName); err != nil {
			return err
		}
	}
	if created.CollectionType != "" {
		if err := formatter.PrintText("%-20s%s\n", "Type:", created.CollectionType); err != nil {
			return err
		}
	}
	if created.StorageGatewayID != "" {
		if err := formatter.PrintText("%-20s%s\n", "Storage Gateway:", created.StorageGatewayID); err != nil {
			return err
		}
	}
	if created.CollectionBaseFolder != "" {
		if err := formatter.PrintText("%-20s%s\n", "Base Path:", created.CollectionBaseFolder); err != nil {
			return err
		}
	}

	return nil
}
