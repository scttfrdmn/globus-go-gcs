package endpoint

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewUpgradeCmd creates the endpoint upgrade command.
func NewUpgradeCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		force        bool
		check        bool
	)

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade endpoint to latest version",
		Long: `Upgrade the Globus Connect Server endpoint to the latest version.

This command performs version compatibility checks, pre-upgrade validation,
and post-upgrade verification. Use --check to see available upgrades without
performing the upgrade.

Example:
  # Check for available upgrades
  globus-connect-server endpoint upgrade \
    --endpoint example.data.globus.org \
    --check

  # Perform upgrade
  globus-connect-server endpoint upgrade \
    --endpoint example.data.globus.org

  # Skip confirmation prompt
  globus-connect-server endpoint upgrade \
    --endpoint example.data.globus.org \
    --force

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runUpgrade(cmd.Context(), profile, format, endpointFQDN, force, check, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")
	cmd.Flags().BoolVar(&check, "check", false, "Check for available upgrades without performing upgrade")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runUpgrade executes the endpoint upgrade command.
func runUpgrade(ctx context.Context, profile, formatStr, endpointFQDN string, force, check bool, out interface{ Write([]byte) (int, error) }) error {
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

	// Check for available upgrades
	upgradeInfo, err := gcsClient.CheckEndpointUpgrade(ctx)
	if err != nil {
		return fmt.Errorf("check endpoint upgrade: %w", err)
	}

	// If --check flag, just display upgrade information
	if check {
		return displayUpgradeInfo(formatter, upgradeInfo)
	}

	// Check if upgrade is needed
	if !upgradeInfo.UpgradeRequired {
		return handleNoUpgradeNeeded(formatter, upgradeInfo)
	}

	// Display upgrade information and confirm
	if err := displayUpgradePrompt(formatter, upgradeInfo); err != nil {
		return err
	}

	// Confirmation prompt unless --force
	if !force {
		if err := confirmUpgrade(upgradeInfo); err != nil {
			return err
		}
	}

	// Perform upgrade
	result, err := gcsClient.UpgradeEndpoint(ctx)
	if err != nil {
		return fmt.Errorf("upgrade endpoint: %w", err)
	}

	// Display results
	return displayUpgradeResult(formatter, result)
}

// handleNoUpgradeNeeded handles the case when no upgrade is needed.
func handleNoUpgradeNeeded(formatter *output.Formatter, info *gcs.UpgradeInfo) error {
	if formatter.IsJSON() {
		return formatter.PrintJSON(map[string]interface{}{
			"message":         "Endpoint is already at the latest version",
			"current_version": info.CurrentVersion,
		})
	}
	return formatter.PrintText("Endpoint is already at the latest version: %s\n", info.CurrentVersion)
}

// displayUpgradePrompt displays upgrade information before performing the upgrade.
func displayUpgradePrompt(formatter *output.Formatter, info *gcs.UpgradeInfo) error {
	if formatter.IsJSON() {
		return nil // Don't display in JSON mode
	}

	if err := formatter.Println("Upgrade Available"); err != nil {
		return err
	}
	if err := formatter.Println("================="); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}
	if err := formatter.PrintText("Current Version: %s\n", info.CurrentVersion); err != nil {
		return err
	}
	if err := formatter.PrintText("Latest Version:  %s\n", info.LatestVersion); err != nil {
		return err
	}
	if info.Compatible {
		if err := formatter.Println("Compatibility:   Compatible"); err != nil {
			return err
		}
	} else {
		if err := formatter.Println("Compatibility:   Incompatible (review release notes)"); err != nil {
			return err
		}
	}
	if len(info.UpgradePath) > 0 {
		if err := formatter.PrintText("Upgrade Path:    %s\n", strings.Join(info.UpgradePath, " -> ")); err != nil {
			return err
		}
	}
	if info.ReleaseNotes != "" {
		if err := formatter.Println(); err != nil {
			return err
		}
		if err := formatter.Println("Release Notes:"); err != nil {
			return err
		}
		if err := formatter.PrintText("%s\n", info.ReleaseNotes); err != nil {
			return err
		}
	}
	return formatter.Println()
}

// confirmUpgrade prompts the user for confirmation.
func confirmUpgrade(info *gcs.UpgradeInfo) error {
	fmt.Fprintf(os.Stderr, "Do you want to upgrade from %s to %s? (yes/no): ",
		info.CurrentVersion, info.LatestVersion)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("read confirmation: %w", err)
	}
	response = strings.TrimSpace(strings.ToLower(response))
	if response != "yes" && response != "y" {
		return fmt.Errorf("upgrade cancelled")
	}
	return nil
}

// displayUpgradeResult displays the result of the upgrade operation.
func displayUpgradeResult(formatter *output.Formatter, result *gcs.UpgradeResult) error {
	if formatter.IsJSON() {
		return formatter.PrintJSON(result)
	}

	// Text format
	if !result.Success {
		if err := formatter.Println("Endpoint upgrade failed"); err != nil {
			return err
		}
		if result.Message != "" {
			if err := formatter.Println(); err != nil {
				return err
			}
			if err := formatter.PrintText("Error: %s\n", result.Message); err != nil {
				return err
			}
		}
		return fmt.Errorf("upgrade failed")
	}

	if err := formatter.Println("Endpoint upgraded successfully!"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}
	if result.PreviousVersion != "" {
		if err := formatter.PrintText("Previous Version: %s\n", result.PreviousVersion); err != nil {
			return err
		}
	}
	if result.NewVersion != "" {
		if err := formatter.PrintText("New Version:      %s\n", result.NewVersion); err != nil {
			return err
		}
	}
	if result.Message != "" {
		if err := formatter.Println(); err != nil {
			return err
		}
		if err := formatter.PrintText("%s\n", result.Message); err != nil {
			return err
		}
	}
	if result.RollbackAvailable {
		if err := formatter.Println(); err != nil {
			return err
		}
		if err := formatter.Println("Note: Rollback is available if needed"); err != nil {
			return err
		}
	}
	return nil
}

// displayUpgradeInfo displays upgrade information without performing the upgrade.
func displayUpgradeInfo(formatter *output.Formatter, info *gcs.UpgradeInfo) error {
	if formatter.IsJSON() {
		return formatter.PrintJSON(info)
	}

	// Text format
	if err := formatter.Println("Endpoint Upgrade Information"); err != nil {
		return err
	}
	if err := formatter.Println("============================"); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	if info.CurrentVersion != "" {
		if err := formatter.PrintText("%-20s%s\n", "Current Version:", info.CurrentVersion); err != nil {
			return err
		}
	}
	if info.LatestVersion != "" {
		if err := formatter.PrintText("%-20s%s\n", "Latest Version:", info.LatestVersion); err != nil {
			return err
		}
	}
	if info.UpgradeRequired {
		if err := formatter.PrintText("%-20s%s\n", "Upgrade Required:", "Yes"); err != nil {
			return err
		}
	} else {
		if err := formatter.PrintText("%-20s%s\n", "Upgrade Required:", "No"); err != nil {
			return err
		}
	}
	if info.Compatible {
		if err := formatter.PrintText("%-20s%s\n", "Compatible:", "Yes"); err != nil {
			return err
		}
	} else {
		if err := formatter.PrintText("%-20s%s\n", "Compatible:", "No"); err != nil {
			return err
		}
	}
	if len(info.UpgradePath) > 0 {
		if err := formatter.PrintText("%-20s%s\n", "Upgrade Path:", strings.Join(info.UpgradePath, " -> ")); err != nil {
			return err
		}
	}

	if info.ReleaseNotes != "" {
		if err := formatter.Println(); err != nil {
			return err
		}
		if err := formatter.Println("Release Notes:"); err != nil {
			return err
		}
		if err := formatter.PrintText("%s\n", info.ReleaseNotes); err != nil {
			return err
		}
	}

	return nil
}
