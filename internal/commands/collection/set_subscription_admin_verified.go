package collection

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewSetSubscriptionAdminVerifiedCmd creates the collection set-subscription-admin-verified command.
func NewSetSubscriptionAdminVerifiedCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		verified     bool
	)

	cmd := &cobra.Command{
		Use:   "set-subscription-admin-verified COLLECTION_ID",
		Short: "Set subscription admin verification status for collection",
		Long: `Set the subscription admin verification status for a collection.

This command controls whether a collection has been verified by the
subscription administrator. Verified collections typically have access
to premium features and resources associated with the subscription.

Example:
  # Mark collection as verified
  globus-connect-server collection set-subscription-admin-verified abc123 \
    --endpoint example.data.globus.org \
    --verified

  # Mark collection as not verified
  globus-connect-server collection set-subscription-admin-verified abc123 \
    --endpoint example.data.globus.org \
    --verified=false

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			collectionID := args[0]
			return runSetSubscriptionAdminVerified(cmd.Context(), profile, format, endpointFQDN, collectionID, verified, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().BoolVar(&verified, "verified", false, "Verification status (true or false)")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runSetSubscriptionAdminVerified executes the collection set-subscription-admin-verified command.
func runSetSubscriptionAdminVerified(ctx context.Context, profile, formatStr, endpointFQDN, collectionID string, verified bool, out interface{ Write([]byte) (int, error) }) error {
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

	// Set subscription admin verified
	if err := gcsClient.SetSubscriptionAdminVerified(ctx, collectionID, verified); err != nil {
		return fmt.Errorf("set subscription admin verified: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		result := map[string]interface{}{
			"status":        "success",
			"collection_id": collectionID,
			"verified":      verified,
			"message":       "Subscription admin verification status set successfully",
		}
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.PrintText("Subscription admin verification status set successfully.\n"); err != nil {
		return err
	}
	if err := formatter.PrintText("Collection ID: %s\n", collectionID); err != nil {
		return err
	}
	verifiedStr := "verified"
	if !verified {
		verifiedStr = "not verified"
	}
	if err := formatter.PrintText("Status: %s\n", verifiedStr); err != nil {
		return err
	}

	return nil
}
