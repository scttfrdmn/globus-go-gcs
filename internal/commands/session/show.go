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

// NewShowCmd creates the session show command.
func NewShowCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Display current authentication session",
		Long: `Display information about the current CLI authentication session.

This command shows session details including authentication method, timeouts,
consents, and other session settings.

Example:
  globus-connect-server session show \
    --endpoint example.data.globus.org

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runShow(cmd.Context(), profile, format, endpointFQDN, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runShow executes the session show command.
func runShow(ctx context.Context, profile, formatStr, endpointFQDN string, out interface{ Write([]byte) (int, error) }) error {
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

	// Get session
	session, err := gcsClient.GetSession(ctx)
	if err != nil {
		return fmt.Errorf("get session: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(session)
	}

	// Text format
	if err := formatter.Println("Session Information"); err != nil {
		return err
	}
	if err := formatter.Println("==================="); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if session.ID != "" {
		if err := formatter.PrintText("%-25s%s\n", "Session ID:", session.ID); err != nil {
			return err
		}
	}
	if session.Principal != "" {
		if err := formatter.PrintText("%-25s%s\n", "Principal:", session.Principal); err != nil {
			return err
		}
	}
	if session.AuthenticationMethod != "" {
		if err := formatter.PrintText("%-25s%s\n", "Authentication Method:", session.AuthenticationMethod); err != nil {
			return err
		}
	}
	if session.SessionTimeoutMins > 0 {
		if err := formatter.PrintText("%-25s%d minutes\n", "Session Timeout:", session.SessionTimeoutMins); err != nil {
			return err
		}
	}
	if session.InactivityTimeoutMins > 0 {
		if err := formatter.PrintText("%-25s%d minutes\n", "Inactivity Timeout:", session.InactivityTimeoutMins); err != nil {
			return err
		}
	}
	if len(session.Consents) > 0 {
		if err := formatter.PrintText("%-25s%s\n", "Consents:", strings.Join(session.Consents, ", ")); err != nil {
			return err
		}
	}
	if len(session.AllowedScopes) > 0 {
		if err := formatter.PrintText("%-25s%s\n", "Allowed Scopes:", strings.Join(session.AllowedScopes, ", ")); err != nil {
			return err
		}
	}

	return nil
}
