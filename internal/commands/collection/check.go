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

// NewCheckCmd creates the collection check command.
func NewCheckCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "check COLLECTION_ID",
		Short: "Validate collection configuration",
		Long: `Validate a collection's configuration and check for issues.

This command performs comprehensive validation of a collection's settings,
including storage gateway connectivity, path accessibility, and permission
configuration. It returns any errors or warnings found.

Example:
  globus-connect-server collection check abc123 \
    --endpoint example.data.globus.org

Requires an active authentication session (use 'login' first).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			collectionID := args[0]
			return runCheck(cmd.Context(), profile, format, endpointFQDN, collectionID, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runCheck executes the collection check command.
func runCheck(ctx context.Context, profile, formatStr, endpointFQDN, collectionID string, out interface{ Write([]byte) (int, error) }) error {
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

	// Check collection
	result, err := gcsClient.CheckCollection(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("check collection: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(result)
	}

	return formatCheckResults(formatter, result)
}

// formatCheckResults formats the validation results in text format.
func formatCheckResults(formatter *output.Formatter, result *gcs.CollectionValidation) error {
	// Header
	if err := formatter.Println("Collection Validation Results"); err != nil {
		return err
	}
	if err := formatter.Println("============================"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	// Status
	if err := formatter.PrintText("%-20s%s\n", "Collection ID:", result.CollectionID); err != nil {
		return err
	}
	status := "Valid"
	if !result.Valid {
		status = "Invalid"
	}
	if err := formatter.PrintText("%-20s%s\n", "Status:", status); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	// Display errors
	if err := formatValidationIssues(formatter, "Errors:", result.Errors); err != nil {
		return err
	}

	// Display warnings
	if err := formatValidationIssues(formatter, "Warnings:", result.Warnings); err != nil {
		return err
	}

	// Success message
	if result.Valid && len(result.Warnings) == 0 {
		if err := formatter.Println("No issues found. Collection is properly configured."); err != nil {
			return err
		}
	}

	return nil
}

// formatValidationIssues formats a list of validation errors or warnings.
func formatValidationIssues(formatter *output.Formatter, header string, issues []gcs.ValidationError) error {
	if len(issues) == 0 {
		return nil
	}

	if err := formatter.Println(header); err != nil {
		return err
	}

	for i, issue := range issues {
		prefix := fmt.Sprintf("  %d. ", i+1)
		if err := formatter.PrintText("%s[%s] %s", prefix, issue.Code, issue.Message); err != nil {
			return err
		}
		if issue.Field != "" {
			if err := formatter.PrintText(" (Field: %s)", issue.Field); err != nil {
				return err
			}
		}
		if err := formatter.Println(); err != nil {
			return err
		}
	}

	return formatter.Println()
}
