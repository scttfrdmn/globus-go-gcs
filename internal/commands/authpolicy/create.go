package authpolicy

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

// NewCreateCmd creates the auth-policy create command.
func NewCreateCmd() *cobra.Command {
	var (
		profile              string
		format               string
		endpointFQDN         string
		name                 string
		description          string
		requireMFA           bool
		requireHighAssurance bool
		allowedDomains       string
		blockedDomains       string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new authentication policy",
		Long: `Create a new authentication policy for the endpoint.

Authentication policies define security requirements such as MFA, high assurance,
and domain-based access controls.

Example:
  globus-connect-server auth-policy create \
    --endpoint example.data.globus.org \
    --name "High Security Policy" \
    --description "Requires MFA for all access" \
    --require-mfa \
    --allowed-domains "example.edu,example.org"

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCreate(cmd.Context(), profile, format, endpointFQDN, name, description,
				requireMFA, requireHighAssurance, allowedDomains, blockedDomains, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&name, "name", "", "Policy name")
	cmd.Flags().StringVar(&description, "description", "", "Policy description")
	cmd.Flags().BoolVar(&requireMFA, "require-mfa", false, "Require multi-factor authentication")
	cmd.Flags().BoolVar(&requireHighAssurance, "require-high-assurance", false, "Require high assurance authentication")
	cmd.Flags().StringVar(&allowedDomains, "allowed-domains", "", "Comma-separated list of allowed domains")
	cmd.Flags().StringVar(&blockedDomains, "blocked-domains", "", "Comma-separated list of blocked domains")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

// runCreate executes the auth-policy create command.
func runCreate(ctx context.Context, profile, formatStr, endpointFQDN, name, description string,
	requireMFA, requireHighAssurance bool, allowedDomains, blockedDomains string,
	out interface{ Write([]byte) (int, error) }) error {

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

	// Build policy
	policy := &gcs.AuthPolicy{
		Name:                 name,
		Description:          description,
		RequireMFA:           requireMFA,
		RequireHighAssurance: requireHighAssurance,
	}

	// Parse domains
	if allowedDomains != "" {
		domains := strings.Split(allowedDomains, ",")
		for i, d := range domains {
			domains[i] = strings.TrimSpace(d)
		}
		policy.AllowedDomains = domains
	}

	if blockedDomains != "" {
		domains := strings.Split(blockedDomains, ",")
		for i, d := range domains {
			domains[i] = strings.TrimSpace(d)
		}
		policy.BlockedDomains = domains
	}

	// Create policy
	created, err := gcsClient.CreateAuthPolicy(ctx, policy)
	if err != nil {
		return fmt.Errorf("create auth policy: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(created)
	}

	// Text format
	if err := formatter.Println("Authentication policy created successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}
	if created.ID != "" {
		if err := formatter.PrintText("%-20s%s\n", "ID:", created.ID); err != nil {
			return err
		}
	}
	if created.Name != "" {
		if err := formatter.PrintText("%-20s%s\n", "Name:", created.Name); err != nil {
			return err
		}
	}

	return nil
}
