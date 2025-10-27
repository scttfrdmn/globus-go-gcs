package usercredential

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewActivescaleCreateCmd creates the activescale-create command.
func NewActivescaleCreateCmd() *cobra.Command {
	var (
		profile          string
		format           string
		endpointFQDN     string
		identityID       string
		storageGatewayID string
		username         string
	)

	cmd := &cobra.Command{
		Use:   "activescale-create",
		Short: "Create ActiveScale user credential",
		Long: `Create user credentials for ActiveScale storage.

ActiveScale is an object storage system. This command creates the necessary
credentials for users to access ActiveScale storage through Globus.

Example:
  globus-connect-server user-credential activescale-create \
    --endpoint example.data.globus.org \
    --identity abc123 \
    --storage-gateway sg-abc \
    --username jdoe

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runActivescaleCreate(cmd.Context(), profile, format, endpointFQDN, identityID,
				storageGatewayID, username, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&identityID, "identity", "", "User identity ID")
	cmd.Flags().StringVar(&storageGatewayID, "storage-gateway", "", "Storage gateway ID")
	cmd.Flags().StringVar(&username, "username", "", "ActiveScale username")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("identity")
	_ = cmd.MarkFlagRequired("storage-gateway")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}

// runActivescaleCreate executes the activescale-create command.
func runActivescaleCreate(ctx context.Context, profile, formatStr, endpointFQDN, identityID,
	storageGatewayID, username string, out interface{ Write([]byte) (int, error) }) error {
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

	// Build credential
	credential := &gcs.UserCredential{
		IdentityID:       identityID,
		StorageGatewayID: storageGatewayID,
		Username:         username,
	}

	// Create credential
	created, err := gcsClient.CreateActivescaleCredential(ctx, credential)
	if err != nil {
		return fmt.Errorf("create activescale credential: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(created)
	}

	// Text format
	if err := formatter.Println("ActiveScale credential created successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}
	if created.ID != "" {
		if err := formatter.PrintText("Credential ID: %s\n", created.ID); err != nil {
			return err
		}
	}

	return nil
}
