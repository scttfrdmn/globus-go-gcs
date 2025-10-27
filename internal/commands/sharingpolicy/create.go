package sharingpolicy

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

// NewCreateCmd creates the sharing-policy create command.
func NewCreateCmd() *cobra.Command {
	var (
		profile            string
		format             string
		endpointFQDN       string
		collectionID       string
		name               string
		description        string
		sharingRestrict    string
		sharingUsersAllow  string
		sharingUsersDeny   string
		sharingGroupsAllow string
		sharingGroupsDeny  string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new sharing policy",
		Long: `Create a new sharing policy for a collection.

Sharing policies control who can create shares on collections and what
restrictions apply to those shares.

Example:
  globus-connect-server sharing-policy create \
    --endpoint example.data.globus.org \
    --collection abc123 \
    --name "Restricted Sharing" \
    --sharing-restrict "private" \
    --sharing-users-allow "user1@globusid.org,user2@globusid.org"

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCreate(cmd.Context(), profile, format, endpointFQDN, collectionID, name, description,
				sharingRestrict, sharingUsersAllow, sharingUsersDeny, sharingGroupsAllow, sharingGroupsDeny,
				cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&collectionID, "collection", "", "Collection ID")
	cmd.Flags().StringVar(&name, "name", "", "Policy name")
	cmd.Flags().StringVar(&description, "description", "", "Policy description")
	cmd.Flags().StringVar(&sharingRestrict, "sharing-restrict", "", "Sharing restriction level")
	cmd.Flags().StringVar(&sharingUsersAllow, "sharing-users-allow", "", "Comma-separated list of allowed user identities")
	cmd.Flags().StringVar(&sharingUsersDeny, "sharing-users-deny", "", "Comma-separated list of denied user identities")
	cmd.Flags().StringVar(&sharingGroupsAllow, "sharing-groups-allow", "", "Comma-separated list of allowed group IDs")
	cmd.Flags().StringVar(&sharingGroupsDeny, "sharing-groups-deny", "", "Comma-separated list of denied group IDs")

	_ = cmd.MarkFlagRequired("endpoint")
	_ = cmd.MarkFlagRequired("collection")

	return cmd
}

// runCreate executes the sharing-policy create command.
func runCreate(ctx context.Context, profile, formatStr, endpointFQDN, collectionID, name, description,
	sharingRestrict, sharingUsersAllow, sharingUsersDeny, sharingGroupsAllow, sharingGroupsDeny string,
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

	// Build sharing policy
	policy := &gcs.SharingPolicy{
		CollectionID:    collectionID,
		Name:            name,
		Description:     description,
		SharingRestrict: sharingRestrict,
	}

	if sharingUsersAllow != "" {
		policy.SharingUsersAllow = strings.Split(sharingUsersAllow, ",")
		for i := range policy.SharingUsersAllow {
			policy.SharingUsersAllow[i] = strings.TrimSpace(policy.SharingUsersAllow[i])
		}
	}

	if sharingUsersDeny != "" {
		policy.SharingUsersDeny = strings.Split(sharingUsersDeny, ",")
		for i := range policy.SharingUsersDeny {
			policy.SharingUsersDeny[i] = strings.TrimSpace(policy.SharingUsersDeny[i])
		}
	}

	if sharingGroupsAllow != "" {
		policy.SharingGroupsAllow = strings.Split(sharingGroupsAllow, ",")
		for i := range policy.SharingGroupsAllow {
			policy.SharingGroupsAllow[i] = strings.TrimSpace(policy.SharingGroupsAllow[i])
		}
	}

	if sharingGroupsDeny != "" {
		policy.SharingGroupsDeny = strings.Split(sharingGroupsDeny, ",")
		for i := range policy.SharingGroupsDeny {
			policy.SharingGroupsDeny[i] = strings.TrimSpace(policy.SharingGroupsDeny[i])
		}
	}

	// Create sharing policy
	created, err := gcsClient.CreateSharingPolicy(ctx, policy)
	if err != nil {
		return fmt.Errorf("create sharing policy: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(created)
	}

	// Text format
	if err := formatter.Println("Sharing policy created successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}
	if created.ID != "" {
		if err := formatter.PrintText("Policy ID: %s\n", created.ID); err != nil {
			return err
		}
	}

	return nil
}
