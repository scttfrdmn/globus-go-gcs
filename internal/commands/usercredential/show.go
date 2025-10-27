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

// NewShowCmd creates the show command.
func NewShowCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "show CREDENTIAL_ID",
		Short: "Display details of a user credential",
		Long: `Display detailed information about a specific user credential.

Shows credential type, identity, storage gateway, and type-specific details
like S3 keys or OAuth tokens.

Example:
  globus-connect-server user-credential show cred-abc123 \
    --endpoint example.data.globus.org

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShow(cmd.Context(), profile, format, endpointFQDN, args[0], cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runShow executes the show command.
func runShow(ctx context.Context, profile, formatStr, endpointFQDN, credentialID string, out interface{ Write([]byte) (int, error) }) error {
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

	// Get credential
	credential, err := gcsClient.GetUserCredential(ctx, credentialID)
	if err != nil {
		return fmt.Errorf("get user credential: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(credential)
	}

	// Text format
	if err := formatter.Println("User Credential Details"); err != nil {
		return err
	}
	if err := formatter.Println("======================="); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if credential.ID != "" {
		if err := formatter.PrintText("%-20s%s\n", "ID:", credential.ID); err != nil {
			return err
		}
	}
	if credential.Type != "" {
		if err := formatter.PrintText("%-20s%s\n", "Type:", credential.Type); err != nil {
			return err
		}
	}
	if credential.IdentityID != "" {
		if err := formatter.PrintText("%-20s%s\n", "Identity ID:", credential.IdentityID); err != nil {
			return err
		}
	}
	if credential.StorageGatewayID != "" {
		if err := formatter.PrintText("%-20s%s\n", "Storage Gateway ID:", credential.StorageGatewayID); err != nil {
			return err
		}
	}
	if credential.Username != "" {
		if err := formatter.PrintText("%-20s%s\n", "Username:", credential.Username); err != nil {
			return err
		}
	}

	// Type-specific details
	if len(credential.S3Keys) > 0 {
		if err := formatter.Println(); err != nil {
			return err
		}
		if err := formatter.Println("S3 Access Keys:"); err != nil {
			return err
		}
		for _, key := range credential.S3Keys {
			if err := formatter.PrintText("  - %s (created: %s)\n", key.AccessKeyID, key.CreatedAt.Format("2006-01-02 15:04:05")); err != nil {
				return err
			}
		}
	}

	if credential.OAuthToken != "" {
		if err := formatter.Println(); err != nil {
			return err
		}
		if err := formatter.PrintText("OAuth Token: %s\n", credential.OAuthToken[:20]+"..."); err != nil {
			return err
		}
	}

	return nil
}
