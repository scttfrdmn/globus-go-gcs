package session

import (
	"context"
	"fmt"
	"strconv"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewUpdateCmd creates the session update command.
func NewUpdateCmd() *cobra.Command {
	var (
		profile             string
		format              string
		endpointFQDN        string
		sessionTimeout      string
		inactivityTimeout   string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update session settings",
		Long: `Update settings for the current CLI authentication session.

You can modify session timeout and inactivity timeout settings.

Example:
  globus-connect-server session update \
    --endpoint example.data.globus.org \
    --session-timeout 60 \
    --inactivity-timeout 30

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runUpdate(cmd.Context(), profile, format, endpointFQDN, sessionTimeout, inactivityTimeout, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&sessionTimeout, "session-timeout", "", "Session timeout in minutes")
	cmd.Flags().StringVar(&inactivityTimeout, "inactivity-timeout", "", "Inactivity timeout in minutes")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runUpdate executes the session update command.
func runUpdate(ctx context.Context, profile, formatStr, endpointFQDN, sessionTimeout, inactivityTimeout string, out interface{ Write([]byte) (int, error) }) error {
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

	// Build session update
	session := &gcs.Session{}

	if sessionTimeout != "" {
		timeout, err := strconv.Atoi(sessionTimeout)
		if err != nil {
			return fmt.Errorf("invalid session timeout: %w", err)
		}
		session.SessionTimeoutMins = timeout
	}

	if inactivityTimeout != "" {
		timeout, err := strconv.Atoi(inactivityTimeout)
		if err != nil {
			return fmt.Errorf("invalid inactivity timeout: %w", err)
		}
		session.InactivityTimeoutMins = timeout
	}

	// Update session
	updated, err := gcsClient.UpdateSession(ctx, session)
	if err != nil {
		return fmt.Errorf("update session: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(updated)
	}

	// Text format
	if err := formatter.Println("Session updated successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}
	if updated.SessionTimeoutMins > 0 {
		if err := formatter.PrintText("Session Timeout:    %d minutes\n", updated.SessionTimeoutMins); err != nil {
			return err
		}
	}
	if updated.InactivityTimeoutMins > 0 {
		if err := formatter.PrintText("Inactivity Timeout: %d minutes\n", updated.InactivityTimeoutMins); err != nil {
			return err
		}
	}

	return nil
}
