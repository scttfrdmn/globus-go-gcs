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

// NewSetSubscriptionIDCmd creates the endpoint set-subscription-id command.
func NewSetSubscriptionIDCmd() *cobra.Command {
	var (
		profile        string
		format         string
		endpointFQDN   string
		subscriptionID string
	)

	cmd := &cobra.Command{
		Use:   "set-subscription-id",
		Short: "Update subscription assignment for endpoint",
		Long: `Update the subscription assignment for the endpoint.

This command associates the endpoint with a specific Globus subscription ID.
Subscriptions provide access to premium features and determine billing
and usage tracking for the endpoint.

Example:
  globus-connect-server endpoint set-subscription-id \
    --endpoint example.data.globus.org \
    --subscription-id "12345678-1234-1234-1234-123456789abc"

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSetSubscriptionID(cmd.Context(), profile, format, endpointFQDN, subscriptionID, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&subscriptionID, "subscription-id", "", "Subscription UUID")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("subscription-id")

	return cmd
}

// runSetSubscriptionID executes the endpoint set-subscription-id command.
func runSetSubscriptionID(ctx context.Context, profile, formatStr, endpointFQDN, subscriptionID string, out interface{ Write([]byte) (int, error) }) error {
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

	// Set subscription ID
	if err := gcsClient.SetSubscriptionID(ctx, subscriptionID); err != nil {
		return fmt.Errorf("set subscription ID: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		result := map[string]string{
			"status":          "success",
			"subscription_id": subscriptionID,
			"message":         "Subscription ID set successfully",
		}
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.PrintText("Subscription ID set successfully.\n"); err != nil {
		return err
	}
	if err := formatter.PrintText("Subscription ID: %s\n", subscriptionID); err != nil {
		return err
	}

	return nil
}
