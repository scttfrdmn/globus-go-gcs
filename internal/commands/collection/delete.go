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

// NewDeleteCmd creates the collection delete command.
func NewDeleteCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "delete COLLECTION_ID",
		Short: "Delete a collection",
		Long: `Delete a collection from the endpoint.

WARNING: This action is permanent and cannot be undone. The collection
and all its configuration will be removed.

Example:
  globus-connect-server collection delete abc123 \
    --endpoint example.data.globus.org

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			collectionID := args[0]
			return runDelete(cmd.Context(), profile, format, endpointFQDN, collectionID, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runDelete executes the collection delete command.
func runDelete(ctx context.Context, profile, formatStr, endpointFQDN, collectionID string, out interface{ Write([]byte) (int, error) }) error {
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

	// Delete collection
	if err := gcsClient.DeleteCollection(ctx, collectionID); err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		result := map[string]string{
			"status":        "success",
			"collection_id": collectionID,
			"message":       "Collection deleted successfully",
		}
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.PrintText("Collection %s deleted successfully.\n", collectionID); err != nil {
		return err
	}

	return nil
}
