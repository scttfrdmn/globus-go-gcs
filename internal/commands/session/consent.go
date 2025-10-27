package session

import (
	"context"
	"fmt"
	"strings"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewConsentCmd creates the session consent command.
func NewConsentCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		consents     string
	)

	cmd := &cobra.Command{
		Use:   "consent",
		Short: "Update session consents",
		Long: `Update the consents for the current CLI authentication session.

Consents control what operations and data access the session permits.

Example:
  globus-connect-server session consent \
    --endpoint example.data.globus.org \
    --consents "data_access,transfer,sharing"

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runConsent(cmd.Context(), profile, format, endpointFQDN, consents, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&consents, "consents", "", "Comma-separated list of consents")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("consents")

	return cmd
}

// runConsent executes the session consent command.
func runConsent(ctx context.Context, profile, formatStr, endpointFQDN, consents string, out interface{ Write([]byte) (int, error) }) error {
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

	// Parse consents
	consentList := strings.Split(consents, ",")
	for i, c := range consentList {
		consentList[i] = strings.TrimSpace(c)
	}

	// Update consents
	updated, err := gcsClient.UpdateSessionConsents(ctx, consentList)
	if err != nil {
		return fmt.Errorf("update session consents: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(updated)
	}

	// Text format
	if err := formatter.Println("Session consents updated successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}
	if len(updated.Consents) > 0 {
		if err := formatter.PrintText("Consents: %s\n", strings.Join(updated.Consents, ", ")); err != nil {
			return err
		}
	}

	return nil
}
