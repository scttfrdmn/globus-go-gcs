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

// NewS3KeysAddCmd creates the s3-keys-add command.
func NewS3KeysAddCmd() *cobra.Command {
	var (
		profile         string
		format          string
		endpointFQDN    string
		credentialID    string
		accessKeyID     string
		secretAccessKey string
	)

	cmd := &cobra.Command{
		Use:   "s3-keys-add",
		Short: "Add S3 IAM access keys",
		Long: `Add S3 IAM access keys to an S3 credential.

S3-compatible storage requires IAM access keys for authentication.
This command adds a key pair to an existing S3 credential.

Example:
  globus-connect-server user-credential s3-keys-add \
    --endpoint example.data.globus.org \
    --credential cred-abc123 \
    --access-key-id AKIAIOSFODNN7EXAMPLE \
    --secret-access-key wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runS3KeysAdd(cmd.Context(), profile, format, endpointFQDN, credentialID,
				accessKeyID, secretAccessKey, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&credentialID, "credential", "", "User credential ID")
	cmd.Flags().StringVar(&accessKeyID, "access-key-id", "", "S3 access key ID")
	cmd.Flags().StringVar(&secretAccessKey, "secret-access-key", "", "S3 secret access key")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("credential")
	_ = cmd.MarkFlagRequired("access-key-id")
	_ = cmd.MarkFlagRequired("secret-access-key")

	return cmd
}

// runS3KeysAdd executes the s3-keys-add command.
func runS3KeysAdd(ctx context.Context, profile, formatStr, endpointFQDN, credentialID,
	accessKeyID, secretAccessKey string, out interface{ Write([]byte) (int, error) }) error {
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

	// Build S3 key
	key := &gcs.S3Key{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}

	// Add S3 key
	updated, err := gcsClient.AddS3Key(ctx, credentialID, key)
	if err != nil {
		return fmt.Errorf("add S3 key: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(updated)
	}

	// Text format
	if err := formatter.Println("S3 access key added successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}
	if err := formatter.PrintText("Access Key ID: %s\n", accessKeyID); err != nil {
		return err
	}

	return nil
}
