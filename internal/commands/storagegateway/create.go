package storagegateway

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

// NewCreateCmd creates the storage gateway create command.
func NewCreateCmd() *cobra.Command {
	var (
		profile            string
		format             string
		endpointFQDN       string
		displayName        string
		connectorID        string
		root               string
		allowedDomains     string
		highAssurance      bool
		requireMFA         bool
		posixStagingFolder string
		posixUserIDMap     string
		posixGroupIDMap    string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new storage gateway",
		Long: `Create a new storage gateway on the endpoint.

A storage gateway connects collections to storage backends.
You must specify the display name, connector ID, and root path.

Common connector IDs:
  - posix: For POSIX filesystems
  - blackpearl: For Spectra Logic BlackPearl
  - azure-blob: For Azure Blob Storage
  - s3: For Amazon S3

Example:
  globus-connect-server storagegateway create \
    --endpoint example.data.globus.org \
    --display-name "My Storage" \
    --connector-id posix \
    --root /data

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCreate(cmd.Context(), profile, format, endpointFQDN,
				displayName, connectorID, root, allowedDomains,
				highAssurance, requireMFA, posixStagingFolder,
				posixUserIDMap, posixGroupIDMap, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&displayName, "display-name", "", "Display name for the storage gateway")
	cmd.Flags().StringVar(&connectorID, "connector-id", "posix", "Connector ID (posix, s3, azure-blob, etc.)")
	cmd.Flags().StringVar(&root, "root", "", "Root path for the storage gateway")
	cmd.Flags().StringVar(&allowedDomains, "allowed-domains", "", "Comma-separated allowed authentication domains")
	cmd.Flags().BoolVar(&highAssurance, "high-assurance", false, "Require high assurance for access")
	cmd.Flags().BoolVar(&requireMFA, "require-mfa", false, "Require multi-factor authentication")
	cmd.Flags().StringVar(&posixStagingFolder, "posix-staging-path", "", "POSIX staging folder path")
	cmd.Flags().StringVar(&posixUserIDMap, "posix-user-id-map", "", "POSIX user ID mapping")
	cmd.Flags().StringVar(&posixGroupIDMap, "posix-group-id-map", "", "POSIX group ID mapping")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("display-name")
	_ = cmd.MarkFlagRequired("root")

	return cmd
}

// runCreate executes the storage gateway create command.
func runCreate(ctx context.Context, profile, formatStr, endpointFQDN string,
	displayName, connectorID, root, allowedDomains string,
	highAssurance, requireMFA bool,
	posixStagingFolder, posixUserIDMap, posixGroupIDMap string,
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

	// Build storage gateway object
	gateway := &gcs.StorageGateway{
		DisplayName:        displayName,
		ConnectorID:        connectorID,
		Root:               root,
		HighAssurance:      highAssurance,
		RequireMFA:         requireMFA,
		PosixStagingFolder: posixStagingFolder,
		PosixUserIDMap:     posixUserIDMap,
		PosixGroupIDMap:    posixGroupIDMap,
	}

	// Parse allowed domains if provided
	if allowedDomains != "" {
		gateway.AllowedDomains = strings.Split(allowedDomains, ",")
		// Trim whitespace from each domain
		for i, domain := range gateway.AllowedDomains {
			gateway.AllowedDomains[i] = strings.TrimSpace(domain)
		}
	}

	// Create storage gateway
	created, err := gcsClient.CreateStorageGateway(ctx, gateway)
	if err != nil {
		return fmt.Errorf("create storage gateway: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(created)
	}

	// Text format
	if err := formatter.Println("Storage gateway created successfully!"); err != nil {
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
	if created.ConnectorID != "" {
		if err := formatter.PrintText("%-20s%s\n", "Connector:", created.ConnectorID); err != nil {
			return err
		}
	}
	if created.Root != "" {
		if err := formatter.PrintText("%-20s%s\n", "Root:", created.Root); err != nil {
			return err
		}
	}

	return nil
}
