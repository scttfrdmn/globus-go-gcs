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

// NewSetOwnerCmd creates the collection set-owner command.
func NewSetOwnerCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		principalURN string
	)

	cmd := &cobra.Command{
		Use:   "set-owner COLLECTION_ID",
		Short: "Designate the owner of a collection",
		Long: `Designate the owner of a collection.

The principal is identified by their URN (Uniform Resource Name), which typically
takes the form: urn:globus:auth:identity:<uuid> for users or
urn:globus:groups:id:<uuid> for groups.

The collection owner has full administrative control over the collection,
including the ability to manage sharing policies, permissions, and other settings.

Example:
  globus-connect-server collection set-owner abc123 \
    --endpoint example.data.globus.org \
    --principal "urn:globus:auth:identity:12345678-1234-1234-1234-123456789abc"

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			collectionID := args[0]
			return runSetOwner(cmd.Context(), profile, format, endpointFQDN, collectionID, principalURN, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&principalURN, "principal", "", "Principal URN (user or group)")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("principal")

	return cmd
}

// runSetOwner executes the collection set-owner command.
func runSetOwner(ctx context.Context, profile, formatStr, endpointFQDN, collectionID, principalURN string, out interface{ Write([]byte) (int, error) }) error {
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

	// Set collection owner
	if err := gcsClient.SetCollectionOwner(ctx, collectionID, principalURN); err != nil {
		return fmt.Errorf("set collection owner: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		result := map[string]string{
			"status":        "success",
			"collection_id": collectionID,
			"principal":     principalURN,
			"message":       "Collection owner set successfully",
		}
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.PrintText("Collection owner set successfully.\n"); err != nil {
		return err
	}
	if err := formatter.PrintText("Collection ID: %s\n", collectionID); err != nil {
		return err
	}
	if err := formatter.PrintText("Principal: %s\n", principalURN); err != nil {
		return err
	}

	return nil
}
