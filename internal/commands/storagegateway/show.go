package storagegateway

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewShowCmd creates the storage gateway show command.
func NewShowCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "show GATEWAY_ID",
		Short: "Display storage gateway details",
		Long: `Display detailed information about a specific storage gateway.

This command retrieves and displays the complete configuration of a storage gateway
including connector type, root path, identity mappings, and security settings.

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gatewayID := args[0]
			return runShow(cmd.Context(), profile, format, endpointFQDN, gatewayID, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runShow executes the storage gateway show command.
func runShow(ctx context.Context, profile, formatStr, endpointFQDN, gatewayID string, out interface{ Write([]byte) (int, error) }) error {
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

	// Get storage gateway
	gateway, err := gcsClient.GetStorageGateway(ctx, gatewayID)
	if err != nil {
		return fmt.Errorf("get storage gateway: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(gateway)
	}

	// Text format
	return formatGatewayText(formatter, gateway)
}

// formatGatewayText formats the storage gateway in text format.
func formatGatewayText(formatter *output.Formatter, gateway *gcs.StorageGateway) error {
	if err := formatter.Println("Storage Gateway Details:"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	// Print basic fields
	if err := printGatewayBasicFields(formatter, gateway); err != nil {
		return err
	}

	// Print identity mappings
	if err := printGatewayIdentityMappings(formatter, gateway); err != nil {
		return err
	}

	// Print security settings
	if err := printGatewaySecuritySettings(formatter, gateway); err != nil {
		return err
	}

	// Print path restrictions
	if err := printGatewayPathRestrictions(formatter, gateway); err != nil {
		return err
	}

	// Print POSIX settings
	if err := printGatewayPOSIXSettings(formatter, gateway); err != nil {
		return err
	}

	return nil
}

// printGatewayBasicFields prints the basic gateway fields.
func printGatewayBasicFields(formatter *output.Formatter, gateway *gcs.StorageGateway) error {
	printField := func(label, value string) error {
		if value != "" {
			return formatter.PrintText("%-25s%s\n", label+":", value)
		}
		return nil
	}

	if err := printField("Display Name", gateway.DisplayName); err != nil {
		return err
	}
	if err := printField("ID", gateway.ID); err != nil {
		return err
	}
	if err := printField("Connector Name", gateway.ConnectorName); err != nil {
		return err
	}
	if err := printField("Connector ID", gateway.ConnectorID); err != nil {
		return err
	}
	if err := printField("Root", gateway.Root); err != nil {
		return err
	}

	return nil
}

// printGatewayIdentityMappings prints identity mappings.
func printGatewayIdentityMappings(formatter *output.Formatter, gateway *gcs.StorageGateway) error {
	if len(gateway.IdentityMappings) == 0 {
		return nil
	}

	if err := formatter.Println(); err != nil {
		return err
	}
	if err := formatter.Println("Identity Mappings:"); err != nil {
		return err
	}

	for i, mapping := range gateway.IdentityMappings {
		if i > 0 {
			if err := formatter.Println(); err != nil {
				return err
			}
		}
		if err := formatter.PrintText("  Protocol:        %s\n", mapping.DataAccessProtocol); err != nil {
			return err
		}
		if mapping.IdentityID != "" {
			if err := formatter.PrintText("  Identity ID:     %s\n", mapping.IdentityID); err != nil {
				return err
			}
		}
		if mapping.LocalUsername != "" {
			if err := formatter.PrintText("  Local Username:  %s\n", mapping.LocalUsername); err != nil {
				return err
			}
		}
	}

	return nil
}

// printGatewaySecuritySettings prints security settings.
func printGatewaySecuritySettings(formatter *output.Formatter, gateway *gcs.StorageGateway) error {
	if err := formatter.Println(); err != nil {
		return err
	}
	if err := formatter.Println("Security Settings:"); err != nil {
		return err
	}

	if err := formatter.PrintText("  High Assurance:  %t\n", gateway.HighAssurance); err != nil {
		return err
	}
	if err := formatter.PrintText("  Require MFA:     %t\n", gateway.RequireMFA); err != nil {
		return err
	}

	if len(gateway.AllowedDomains) > 0 {
		if err := formatter.Println(); err != nil {
			return err
		}
		if err := formatter.Println("  Allowed Domains:"); err != nil {
			return err
		}
		for _, domain := range gateway.AllowedDomains {
			if err := formatter.PrintText("    - %s\n", domain); err != nil {
				return err
			}
		}
	}

	return nil
}

// printGatewayPathRestrictions prints path restrictions.
func printGatewayPathRestrictions(formatter *output.Formatter, gateway *gcs.StorageGateway) error {
	if gateway.RestrictPaths == nil {
		return nil
	}

	if err := formatter.Println(); err != nil {
		return err
	}
	if err := formatter.Println("Path Restrictions:"); err != nil {
		return err
	}

	if len(gateway.RestrictPaths.ReadOnly) > 0 {
		if err := formatter.Println("  Read-Only:"); err != nil {
			return err
		}
		for _, path := range gateway.RestrictPaths.ReadOnly {
			if err := formatter.PrintText("    - %s\n", path); err != nil {
				return err
			}
		}
	}

	if len(gateway.RestrictPaths.ReadWrite) > 0 {
		if err := formatter.Println("  Read-Write:"); err != nil {
			return err
		}
		for _, path := range gateway.RestrictPaths.ReadWrite {
			if err := formatter.PrintText("    - %s\n", path); err != nil {
				return err
			}
		}
	}

	if len(gateway.RestrictPaths.None) > 0 {
		if err := formatter.Println("  No Access:"); err != nil {
			return err
		}
		for _, path := range gateway.RestrictPaths.None {
			if err := formatter.PrintText("    - %s\n", path); err != nil {
				return err
			}
		}
	}

	return nil
}

// printGatewayPOSIXSettings prints POSIX-specific settings.
func printGatewayPOSIXSettings(formatter *output.Formatter, gateway *gcs.StorageGateway) error {
	hasPOSIXSettings := gateway.PosixStagingFolder != "" ||
		gateway.PosixUserIDMap != "" ||
		gateway.PosixGroupIDMap != ""

	if !hasPOSIXSettings {
		return nil
	}

	if err := formatter.Println(); err != nil {
		return err
	}
	if err := formatter.Println("POSIX Settings:"); err != nil {
		return err
	}

	if gateway.PosixStagingFolder != "" {
		if err := formatter.PrintText("  Staging Folder:  %s\n", gateway.PosixStagingFolder); err != nil {
			return err
		}
	}
	if gateway.PosixUserIDMap != "" {
		if err := formatter.PrintText("  User ID Map:     %s\n", gateway.PosixUserIDMap); err != nil {
			return err
		}
	}
	if gateway.PosixGroupIDMap != "" {
		if err := formatter.PrintText("  Group ID Map:    %s\n", gateway.PosixGroupIDMap); err != nil {
			return err
		}
	}

	return nil
}
