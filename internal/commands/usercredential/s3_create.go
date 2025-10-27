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

// NewS3CreateCmd creates the s3-create command.
func NewS3CreateCmd() *cobra.Command {
	var (
		profile          string
		format           string
		endpointFQDN     string
		identityID       string
		storageGatewayID string
		username         string
	)

	cmd := &cobra.Command{
		Use:   "s3-create",
		Short: "Create S3 user credential",
		Long: `Create user credentials for S3-compatible storage.

S3 credentials enable users to access S3-compatible storage systems.
After creating the credential, use s3-keys-add to add IAM access keys.

Example:
  globus-connect-server user-credential s3-create \
    --endpoint example.data.globus.org \
    --identity abc123 \
    --storage-gateway sg-abc \
    --username jdoe

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runS3Create(cmd.Context(), profile, format, endpointFQDN, identityID,
				storageGatewayID, username, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&identityID, "identity", "", "User identity ID")
	cmd.Flags().StringVar(&storageGatewayID, "storage-gateway", "", "Storage gateway ID")
	cmd.Flags().StringVar(&username, "username", "", "S3 username")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("identity")
	_ = cmd.MarkFlagRequired("storage-gateway")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}

// runS3Create executes the s3-create command.
func runS3Create(ctx context.Context, profile, formatStr, endpointFQDN, identityID,
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
	created, err := gcsClient.CreateS3Credential(ctx, credential)
	if err != nil {
		return fmt.Errorf("create s3 credential: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(created)
	}

	// Text format
	if err := formatter.Println("S3 credential created successfully!"); err != nil {
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
	if err := formatter.Println(); err != nil {
		return err
	}
	if err := formatter.Println("Use 's3-keys-add' to add IAM access keys to this credential."); err != nil {
		return err
	}

	return nil
}
