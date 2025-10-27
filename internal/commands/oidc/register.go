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

// NewRegisterCmd creates the oidc register command.
func NewRegisterCmd() *cobra.Command {
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
		Use:   "register",
		Short: "Register existing OIDC server",
		Long: `Register an existing OIDC server with the endpoint.

Example:
  globus-connect-server oidc register \
    --endpoint example.data.globus.org \
    --issuer "https://id.example.org" \
    --client-id "my-client-id" \
    --client-secret "my-secret"

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runRegister(cmd.Context(), profile, format, endpointFQDN, issuer, clientID, clientSecret, audience, scopes, cmd.OutOrStdout())
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
	_ = cmd.MarkFlagRequired("issuer")
	_ = cmd.MarkFlagRequired("client-id")
	_ = cmd.MarkFlagRequired("client-secret")

	return cmd
}

// runRegister executes the oidc register command.
func runRegister(ctx context.Context, profile, formatStr, endpointFQDN, issuer, clientID, clientSecret, audience, scopes string, out interface{ Write([]byte) (int, error) }) error {
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

	server := &gcs.OIDCServer{
		Issuer:       issuer,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     audience,
	}

	if scopes != "" {
		scopeList := strings.Split(scopes, ",")
		for i, s := range scopeList {
			scopeList[i] = strings.TrimSpace(s)
		}
		server.Scopes = scopeList
	}

	registered, err := gcsClient.RegisterOIDCServer(ctx, server)
	if err != nil {
		return fmt.Errorf("register OIDC server: %w", err)
	}

	if formatter.IsJSON() {
		return formatter.PrintJSON(registered)
	}

	if err := formatter.Println("OIDC server registered successfully!"); err != nil {
		return err
	}
	if registered.ID != "" {
		if err := formatter.PrintText("ID: %s\n", registered.ID); err != nil {
			return err
		}
	}

	return nil
}
