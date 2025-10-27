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

// NewUpdateCmd creates the storage gateway update command.
func NewUpdateCmd() *cobra.Command {
	var (
		profile            string
		format             string
		endpointFQDN       string
		displayName        string
		allowedDomains     string
		highAssurance      *bool
		requireMFA         *bool
		posixStagingFolder string
		posixUserIDMap     string
		posixGroupIDMap    string
	)

	cmd := &cobra.Command{
		Use:   "update GATEWAY_ID",
		Short: "Update an existing storage gateway",
		Long: `Update an existing storage gateway's configuration.

Only the fields you specify will be updated. Other fields will remain unchanged.

Example:
  globus-connect-server storagegateway update abc123 \
    --endpoint example.data.globus.org \
    --display-name "Updated Storage Gateway" \
    --high-assurance true

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gatewayID := args[0]
			return runUpdate(cmd.Context(), profile, format, endpointFQDN, gatewayID,
				displayName, allowedDomains, highAssurance, requireMFA,
				posixStagingFolder, posixUserIDMap, posixGroupIDMap, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&displayName, "display-name", "", "Display name for the storage gateway")
	cmd.Flags().StringVar(&allowedDomains, "allowed-domains", "", "Comma-separated allowed authentication domains")

	// Boolean flags need special handling
	highAssuranceStr := cmd.Flags().String("high-assurance", "", "Require high assurance for access (true/false)")
	cmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		if *highAssuranceStr != "" {
			val := *highAssuranceStr == "true"
			highAssurance = &val
		}
		requireMFAStr, _ := cmd.Flags().GetString("require-mfa")
		if requireMFAStr != "" {
			val := requireMFAStr == "true"
			requireMFA = &val
		}
		return nil
	}

	cmd.Flags().String("require-mfa", "", "Require multi-factor authentication (true/false)")
	cmd.Flags().StringVar(&posixStagingFolder, "posix-staging-path", "", "POSIX staging folder path")
	cmd.Flags().StringVar(&posixUserIDMap, "posix-user-id-map", "", "POSIX user ID mapping")
	cmd.Flags().StringVar(&posixGroupIDMap, "posix-group-id-map", "", "POSIX group ID mapping")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runUpdate executes the storage gateway update command.
func runUpdate(ctx context.Context, profile, formatStr, endpointFQDN, gatewayID string,
	displayName, allowedDomains string, highAssurance, requireMFA *bool,
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

	// Build storage gateway object with only the fields that were specified
	gateway := &gcs.StorageGateway{}

	if displayName != "" {
		gateway.DisplayName = displayName
	}
	if highAssurance != nil {
		gateway.HighAssurance = *highAssurance
	}
	if requireMFA != nil {
		gateway.RequireMFA = *requireMFA
	}
	if posixStagingFolder != "" {
		gateway.PosixStagingFolder = posixStagingFolder
	}
	if posixUserIDMap != "" {
		gateway.PosixUserIDMap = posixUserIDMap
	}
	if posixGroupIDMap != "" {
		gateway.PosixGroupIDMap = posixGroupIDMap
	}

	// Parse allowed domains if provided
	if allowedDomains != "" {
		gateway.AllowedDomains = strings.Split(allowedDomains, ",")
		// Trim whitespace from each domain
		for i, domain := range gateway.AllowedDomains {
			gateway.AllowedDomains[i] = strings.TrimSpace(domain)
		}
	}

	// Update storage gateway
	updated, err := gcsClient.UpdateStorageGateway(ctx, gatewayID, gateway)
	if err != nil {
		return fmt.Errorf("update storage gateway: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(updated)
	}

	// Text format
	if err := formatter.Println("Storage gateway updated successfully!"); err != nil {
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
	if updated.ConnectorID != "" {
		if err := formatter.PrintText("%-20s%s\n", "Connector:", updated.ConnectorID); err != nil {
			return err
		}
	}

	return nil
}
