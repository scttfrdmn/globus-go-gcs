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

// NewBatchDeleteCmd creates the collection batch-delete command.
func NewBatchDeleteCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		force        bool
	)

	cmd := &cobra.Command{
		Use:   "batch-delete COLLECTION_ID [COLLECTION_ID...]",
		Short: "Delete multiple collections in one operation",
		Long: `Delete multiple collections in a single batch operation.

This command allows you to delete several collections at once. If some
deletions fail, the command will continue and report which succeeded
and which failed.

WARNING: This action is permanent and cannot be undone.

Example:
  globus-connect-server collection batch-delete abc123 def456 ghi789 \
    --endpoint example.data.globus.org \
    --force

Use --force to skip confirmation prompt.

Requires an active authentication session (use 'login' first).`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			collectionIDs := args
			return runBatchDelete(cmd.Context(), profile, format, endpointFQDN, collectionIDs, force, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runBatchDelete executes the collection batch-delete command.
func runBatchDelete(ctx context.Context, profile, formatStr, endpointFQDN string, collectionIDs []string, force bool, out interface{ Write([]byte) (int, error) }) error {
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

	// Confirmation prompt (unless --force)
	if !force {
		if err := formatter.PrintText("WARNING: This will permanently delete %d collection(s).\\n", len(collectionIDs)); err != nil {
			return err
		}
		if err := formatter.Println("This action cannot be undone."); err != nil {
			return err
		}
		if err := formatter.Println(); err != nil {
			return err
		}
		if err := formatter.PrintText("Collections to delete:\\n"); err != nil {
			return err
		}
		for _, id := range collectionIDs {
			if err := formatter.PrintText("  - %s\\n", id); err != nil {
				return err
			}
		}
		if err := formatter.Println(); err != nil {
			return err
		}
		if err := formatter.PrintText("To proceed, use --force flag.\\n"); err != nil {
			return err
		}
		return fmt.Errorf("batch delete cancelled (use --force to proceed)")
	}

	// Create GCS client
	gcsClient, err := gcs.NewClient(
		endpointFQDN,
		gcs.WithAccessToken(token.AccessToken),
	)
	if err != nil {
		return fmt.Errorf("create GCS client: %w", err)
	}

	// Batch delete collections
	result, err := gcsClient.BatchDeleteCollections(ctx, collectionIDs)
	if err != nil {
		return fmt.Errorf("batch delete collections: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(result)
	}

	// Text format
	if err := formatter.Println("Batch Delete Results"); err != nil {
		return err
	}
	if err := formatter.Println("===================="); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	// Show successful deletions
	if len(result.Deleted) > 0 {
		if err := formatter.PrintText("Successfully deleted %d collection(s):\\n", len(result.Deleted)); err != nil {
			return err
		}
		for _, id := range result.Deleted {
			if err := formatter.PrintText("  ✓ %s\\n", id); err != nil {
				return err
			}
		}
		if err := formatter.Println(); err != nil {
			return err
		}
	}

	// Show failed deletions
	if len(result.Failed) > 0 {
		if err := formatter.PrintText("Failed to delete %d collection(s):\\n", len(result.Failed)); err != nil {
			return err
		}
		for _, f := range result.Failed {
			if err := formatter.PrintText("  ✗ %s: %s\\n", f.CollectionID, f.Error); err != nil {
				return err
			}
		}
		if err := formatter.Println(); err != nil {
			return err
		}
	}

	if len(result.Failed) > 0 {
		return fmt.Errorf("batch delete completed with %d failure(s)", len(result.Failed))
	}

	return nil
}
