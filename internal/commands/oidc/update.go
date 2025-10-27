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

// NewUpdateCmd creates the oidc update command.
func NewUpdateCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		issuer       string
		clientID     string
		clientSecret string
		audience     string
		scopes       string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update OIDC server configuration",
		Long: `Update the OIDC server configuration.

Only specified fields will be updated.

Example:
  globus-connect-server oidc update \
    --endpoint example.data.globus.org \
    --client-secret "new-secret"

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runUpdate(cmd.Context(), profile, format, endpointFQDN, issuer, clientID, clientSecret, audience, scopes, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN")
	cmd.Flags().StringVar(&issuer, "issuer", "", "OIDC issuer URL")
	cmd.Flags().StringVar(&clientID, "client-id", "", "OAuth2 client ID")
	cmd.Flags().StringVar(&clientSecret, "client-secret", "", "OAuth2 client secret")
	cmd.Flags().StringVar(&audience, "audience", "", "OAuth2 audience")
	cmd.Flags().StringVar(&scopes, "scopes", "", "Comma-separated list of scopes")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runUpdate executes the oidc update command.
func runUpdate(ctx context.Context, profile, formatStr, endpointFQDN, issuer, clientID, clientSecret, audience, scopes string, out interface{ Write([]byte) (int, error) }) error {
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

	server := &gcs.OIDCServer{}

	if issuer != "" {
		server.Issuer = issuer
	}
	if clientID != "" {
		server.ClientID = clientID
	}
	if clientSecret != "" {
		server.ClientSecret = clientSecret
	}
	if audience != "" {
		server.Audience = audience
	}
	if scopes != "" {
		scopeList := strings.Split(scopes, ",")
		for i, s := range scopeList {
			scopeList[i] = strings.TrimSpace(s)
		}
		server.Scopes = scopeList
	}

	updated, err := gcsClient.UpdateOIDCServer(ctx, server)
	if err != nil {
		return fmt.Errorf("update OIDC server: %w", err)
	}

	if formatter.IsJSON() {
		return formatter.PrintJSON(updated)
	}

	if err := formatter.Println("OIDC server updated successfully!"); err != nil {
		return err
	}
	if updated.ID != "" {
		if err := formatter.PrintText("ID: %s\n", updated.ID); err != nil {
			return err
		}
	}

	return nil
}
