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

// NewOAuthCreateCmd creates the oauth-create command.
func NewOAuthCreateCmd() *cobra.Command {
	var (
		profile          string
		format           string
		endpointFQDN     string
		identityID       string
		storageGatewayID string
		oauthToken       string
	)

	cmd := &cobra.Command{
		Use:   "oauth-create",
		Short: "Create OAuth2 user credential",
		Long: `Create user credentials for OAuth2-based storage.

OAuth2 credentials enable users to access storage systems that use
OAuth2 for authentication.

Example:
  globus-connect-server user-credential oauth-create \
    --endpoint example.data.globus.org \
    --identity abc123 \
    --storage-gateway sg-abc \
    --oauth-token "token_string"

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runOAuthCreate(cmd.Context(), profile, format, endpointFQDN, identityID,
				storageGatewayID, oauthToken, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&identityID, "identity", "", "User identity ID")
	cmd.Flags().StringVar(&storageGatewayID, "storage-gateway", "", "Storage gateway ID")
	cmd.Flags().StringVar(&oauthToken, "oauth-token", "", "OAuth2 access token")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("identity")
	_ = cmd.MarkFlagRequired("storage-gateway")
	_ = cmd.MarkFlagRequired("oauth-token")

	return cmd
}

// runOAuthCreate executes the oauth-create command.
func runOAuthCreate(ctx context.Context, profile, formatStr, endpointFQDN, identityID,
	storageGatewayID, oauthToken string, out interface{ Write([]byte) (int, error) }) error {
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
		OAuthToken:       oauthToken,
	}

	// Create credential
	created, err := gcsClient.CreateOAuthCredential(ctx, credential)
	if err != nil {
		return fmt.Errorf("create oauth credential: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(created)
	}

	// Text format
	if err := formatter.Println("OAuth2 credential created successfully!"); err != nil {
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
