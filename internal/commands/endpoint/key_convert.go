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

// NewKeyConvertCmd creates the endpoint key-convert command.
func NewKeyConvertCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		oldKey       string
	)

	cmd := &cobra.Command{
		Use:   "key-convert",
		Short: "Convert deployment key to new format",
		Long: `Convert an old deployment key to a new format.

This command is used when migrating from an older version of GCS
or when key rotation is required for security purposes.

Example:
  globus-connect-server endpoint key-convert \
    --endpoint example.data.globus.org \
    --old-key "old-deployment-key-string"

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runKeyConvert(cmd.Context(), profile, format, endpointFQDN, oldKey, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&oldKey, "old-key", "", "Old deployment key to convert")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("old-key")

	return cmd
}

// runKeyConvert executes the endpoint key-convert command.
func runKeyConvert(ctx context.Context, profile, formatStr, endpointFQDN, oldKey string, out interface{ Write([]byte) (int, error) }) error {
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

	// Convert key
	result, err := gcsClient.ConvertDeploymentKey(ctx, oldKey)
	if err != nil {
		return fmt.Errorf("convert deployment key: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.Println("Deployment key converted successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if result.NewKey != "" {
		if err := formatter.PrintText("%-20s%s\n", "New Key:", result.NewKey); err != nil {
			return err
		}
	}

	if err := formatter.Println(); err != nil {
		return err
	}
	if err := formatter.Println("IMPORTANT: Save the new key securely. The old key is now invalid."); err != nil {
		return err
	}

	return nil
}
