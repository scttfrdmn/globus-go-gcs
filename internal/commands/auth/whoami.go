package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	globusauth "github.com/scttfrdmn/globus-go-sdk/v3/pkg/services/auth"
	"github.com/spf13/cobra"
)

// NewWhoamiCmd creates the whoami command.
func NewWhoamiCmd() *cobra.Command {
	var (
		profile string
		format  string
	)

	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Show current authenticated identity",
		Long: `Show information about the current authenticated identity.

This command displays your Globus identity information including
username, email, and organization.

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runWhoami(cmd.Context(), profile, format, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")

	return cmd
}

// runWhoami executes the whoami command.
func runWhoami(ctx context.Context, profile, formatStr string, out interface{ Write([]byte) (int, error) }) error {
	// Load token
	token, err := auth.LoadToken(profile)
	if err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	// Check if token is valid
	if !token.IsValid() {
		return fmt.Errorf("token expired, please login again")
	}

	// Create output formatter
	formatter := output.NewFormatter(output.Format(formatStr), out)

	// Create auth client (for potential future API calls)
	cfg, err := config.LoadClientConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	authClient, err := globusauth.NewClient(
		globusauth.WithClientID(cfg.ClientID),
		globusauth.WithClientSecret(cfg.ClientSecret),
	)
	if err != nil {
		return fmt.Errorf("create auth client: %w", err)
	}

	// Introspect token to get user info
	introspectResp, err := authClient.IntrospectToken(ctx, token.AccessToken)
	if err != nil {
		return fmt.Errorf("introspect token: %w", err)
	}

	if !introspectResp.Active {
		return fmt.Errorf("token is not active, please login again")
	}

	// Prepare output data
	info := map[string]interface{}{
		"username":        introspectResp.Username,
		"email":           introspectResp.Email,
		"name":            introspectResp.Name,
		"sub":             introspectResp.Subject,
		"profile":         profile,
		"expires_at":      token.ExpiresAt.Format(time.RFC3339),
		"scopes":          token.Scopes,
		"resource_server": token.ResourceServer,
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(info)
	}

	// Text format
	if err := formatter.Println("Authenticated User Information:"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if introspectResp.Name != "" {
		if err := formatter.PrintText("Name:     %s\n", introspectResp.Name); err != nil {
			return err
		}
	}

	if introspectResp.Username != "" {
		if err := formatter.PrintText("Username: %s\n", introspectResp.Username); err != nil {
			return err
		}
	}

	if introspectResp.Email != "" {
		if err := formatter.PrintText("Email:    %s\n", introspectResp.Email); err != nil {
			return err
		}
	}

	if err := formatter.PrintText("ID:       %s\n", introspectResp.Subject); err != nil {
		return err
	}

	if err := formatter.PrintText("Profile:  %s\n", profile); err != nil {
		return err
	}

	if err := formatter.PrintText("Expires:  %s\n", token.ExpiresAt.Format(time.RFC3339)); err != nil {
		return err
	}

	return nil
}
