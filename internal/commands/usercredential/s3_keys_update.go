package usercredential

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/internal/secureinput"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewS3KeysUpdateCmd creates the s3-keys-update command.
func NewS3KeysUpdateCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		credentialID string
		accessKeyID  string
		secretStdin  bool
		secretEnv    string
	)

	cmd := &cobra.Command{
		Use:   "s3-keys-update",
		Short: "Update S3 IAM access keys",
		Long: `Update the secret access key for an existing S3 IAM key.

Use this command to rotate the secret access key for an S3 credential
while keeping the same access key ID.

🔒 SECURITY: Secret access key is read securely (not from command line).

Methods (in priority order):
  1. Environment variable: --secret-env ENV_VAR_NAME
  2. Stdin pipe: echo "secret" | ... --secret-stdin
  3. Interactive prompt (default, recommended)

Examples:
  # Interactive prompt (most secure, recommended)
  globus-connect-server user-credential s3-keys-update \
    --endpoint example.data.globus.org \
    --credential cred-abc123 \
    --access-key-id AKIAIOSFODNN7EXAMPLE
  # You will be prompted securely for the new secret

  # From stdin (for scripts)
  echo "newSecretKeyEXAMPLE" | \
    globus-connect-server user-credential s3-keys-update \
      --endpoint example.data.globus.org \
      --credential cred-abc123 \
      --access-key-id AKIAIOSFODNN7EXAMPLE \
      --secret-stdin

  # From environment variable
  export S3_SECRET="newSecretKeyEXAMPLE"
  globus-connect-server user-credential s3-keys-update \
    --endpoint example.data.globus.org \
    --credential cred-abc123 \
    --access-key-id AKIAIOSFODNN7EXAMPLE \
    --secret-env S3_SECRET

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runS3KeysUpdate(cmd.Context(), profile, format, endpointFQDN, credentialID,
				accessKeyID, secretStdin, secretEnv, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&credentialID, "credential", "", "User credential ID")
	cmd.Flags().StringVar(&accessKeyID, "access-key-id", "", "S3 access key ID")
	cmd.Flags().BoolVar(&secretStdin, "secret-stdin", false, "Read new secret access key from stdin")
	cmd.Flags().StringVar(&secretEnv, "secret-env", "", "Read new secret access key from environment variable")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("credential")
	_ = cmd.MarkFlagRequired("access-key-id")

	return cmd
}

// runS3KeysUpdate executes the s3-keys-update command.
func runS3KeysUpdate(ctx context.Context, profile, formatStr, endpointFQDN, credentialID,
	accessKeyID string, secretStdin bool, secretEnv string, out interface{ Write([]byte) (int, error) }) error {
	// Load token
	token, err := auth.LoadToken(profile)
	if err != nil {
		return fmt.Errorf("not logged in: %w (use 'login' command first)", err)
	}

	// Check if token is valid
	if !token.IsValid() {
		return fmt.Errorf("token expired, please login again")
	}

	// Read secret access key securely
	secretAccessKey, err := secureinput.ReadSecret(secureinput.ReadSecretOptions{
		PromptMessage: "Enter new S3 secret access key",
		UseStdin:      secretStdin,
		EnvVar:        secretEnv,
		AllowEmpty:    false,
	})
	if err != nil {
		return fmt.Errorf("read secret access key: %w", err)
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

	// Update S3 key
	updated, err := gcsClient.UpdateS3Key(ctx, credentialID, accessKeyID, key)
	if err != nil {
		return fmt.Errorf("update S3 key: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(updated)
	}

	// Text format
	if err := formatter.Println("S3 access key updated successfully!"); err != nil {
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
