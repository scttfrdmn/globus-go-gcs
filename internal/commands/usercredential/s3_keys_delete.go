package usercredential

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewS3KeysDeleteCmd creates the s3-keys-delete command.
func NewS3KeysDeleteCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		credentialID string
		accessKeyID  string
		force        bool
	)

	cmd := &cobra.Command{
		Use:   "s3-keys-delete",
		Short: "Delete S3 IAM access keys",
		Long: `Delete S3 IAM access keys from an S3 credential.

This removes the specified access key from the credential. Use --force
to skip the confirmation prompt.

Example:
  globus-connect-server user-credential s3-keys-delete \
    --endpoint example.data.globus.org \
    --credential cred-abc123 \
    --access-key-id AKIAIOSFODNN7EXAMPLE

  # Skip confirmation
  globus-connect-server user-credential s3-keys-delete \
    --endpoint example.data.globus.org \
    --credential cred-abc123 \
    --access-key-id AKIAIOSFODNN7EXAMPLE \
    --force

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runS3KeysDelete(cmd.Context(), profile, format, endpointFQDN, credentialID,
				accessKeyID, force, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&credentialID, "credential", "", "User credential ID")
	cmd.Flags().StringVar(&accessKeyID, "access-key-id", "", "S3 access key ID to delete")
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("credential")
	_ = cmd.MarkFlagRequired("access-key-id")

	return cmd
}

// runS3KeysDelete executes the s3-keys-delete command.
func runS3KeysDelete(ctx context.Context, profile, formatStr, endpointFQDN, credentialID,
	accessKeyID string, force bool, out interface{ Write([]byte) (int, error) }) error {
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

	// Confirmation prompt unless --force
	if !force {
		fmt.Fprintf(os.Stderr, "Are you sure you want to delete S3 key %s? (yes/no): ", accessKeyID)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read confirmation: %w", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "yes" && response != "y" {
			return fmt.Errorf("deletion cancelled")
		}
	}

	// Create GCS client
	gcsClient, err := gcs.NewClient(
		endpointFQDN,
		gcs.WithAccessToken(token.AccessToken),
	)
	if err != nil {
		return fmt.Errorf("create GCS client: %w", err)
	}

	// Delete S3 key
	if err := gcsClient.DeleteS3Key(ctx, credentialID, accessKeyID); err != nil {
		return fmt.Errorf("delete S3 key: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(map[string]string{"status": "deleted", "access_key_id": accessKeyID})
	}

	// Text format
	return formatter.PrintText("S3 access key %s deleted successfully\n", accessKeyID)
}
