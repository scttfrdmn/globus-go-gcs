package oidc

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewDeleteCmd creates the oidc delete command.
func NewDeleteCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		force        bool
	)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete OIDC server configuration",
		Long: `Delete the OIDC server configuration.

WARNING: This action cannot be undone and will disable OIDC authentication.

Example:
  globus-connect-server oidc delete \
    --endpoint example.data.globus.org \
    --force

Use --force to skip confirmation prompt.

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDelete(cmd.Context(), profile, format, endpointFQDN, force, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN")
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runDelete executes the oidc delete command.
func runDelete(ctx context.Context, profile, formatStr, endpointFQDN string, force bool, out interface{ Write([]byte) (int, error) }) error {
	token, err := auth.LoadToken(profile)
	if err != nil {
		return fmt.Errorf("not logged in: %w (use 'login' command first)", err)
	}

	if !token.IsValid() {
		return fmt.Errorf("token expired, please login again")
	}

	formatter := output.NewFormatter(output.Format(formatStr), out)

	if !force {
		if err := formatter.Println("WARNING: This will permanently delete the OIDC server configuration."); err != nil {
			return err
		}
		if err := formatter.Println("This action cannot be undone."); err != nil {
			return err
		}
		if err := formatter.Println(); err != nil {
			return err
		}
		if err := formatter.PrintText("To proceed, use --force flag.\n"); err != nil {
			return err
		}
		return fmt.Errorf("delete cancelled (use --force to proceed)")
	}

	gcsClient, err := gcs.NewClient(endpointFQDN, gcs.WithAccessToken(token.AccessToken))
	if err != nil {
		return fmt.Errorf("create GCS client: %w", err)
	}

	if err := gcsClient.DeleteOIDCServer(ctx); err != nil {
		return fmt.Errorf("delete OIDC server: %w", err)
	}

	if formatter.IsJSON() {
		result := map[string]string{
			"status":  "success",
			"message": "OIDC server deleted successfully",
		}
		return formatter.PrintJSON(result)
	}

	if err := formatter.Println("OIDC server deleted successfully."); err != nil {
		return err
	}

	return nil
}
