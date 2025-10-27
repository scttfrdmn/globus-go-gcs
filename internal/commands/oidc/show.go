package oidc

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

// NewShowCmd creates the oidc show command.
func NewShowCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
	)

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Display OIDC server configuration",
		Long: `Display the current OIDC server configuration.

Example:
  globus-connect-server oidc show \
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

// runShow executes the oidc show command.
func runShow(ctx context.Context, profile, formatStr, endpointFQDN string, out interface{ Write([]byte) (int, error) }) error {
	token, err := auth.LoadToken(profile)
	if err != nil {
		return fmt.Errorf("not logged in: %w (use 'login' command first)", err)
	}

	if !token.IsValid() {
		return fmt.Errorf("token expired, please login again")
	}

	formatter := output.NewFormatter(output.Format(formatStr), out)

	gcsClient, err := gcs.NewClient(endpointFQDN, gcs.WithAccessToken(token.AccessToken))
	if err != nil {
		return fmt.Errorf("create GCS client: %w", err)
	}

	server, err := gcsClient.GetOIDCServer(ctx)
	if err != nil {
		return fmt.Errorf("get OIDC server: %w", err)
	}

	if formatter.IsJSON() {
		return formatter.PrintJSON(server)
	}

	if err := formatter.Println("OIDC Server Configuration"); err != nil {
		return err
	}
	if err := formatter.Println("========================="); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if server.ID != "" {
		if err := formatter.PrintText("%-20s%s\n", "ID:", server.ID); err != nil {
			return err
		}
	}
	if server.Issuer != "" {
		if err := formatter.PrintText("%-20s%s\n", "Issuer:", server.Issuer); err != nil {
			return err
		}
	}
	if server.ClientID != "" {
		if err := formatter.PrintText("%-20s%s\n", "Client ID:", server.ClientID); err != nil {
			return err
		}
	}
	if server.Audience != "" {
		if err := formatter.PrintText("%-20s%s\n", "Audience:", server.Audience); err != nil {
			return err
		}
	}
	if len(server.Scopes) > 0 {
		if err := formatter.PrintText("%-20s%s\n", "Scopes:", strings.Join(server.Scopes, ", ")); err != nil {
			return err
		}
	}

	return nil
}
