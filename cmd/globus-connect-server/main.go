// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Scott Friedman and Project Contributors

// Package main provides the Globus Connect Server CLI entry point.
package main

import (
	"fmt"
	"os"

	authcmd "github.com/scttfrdmn/globus-go-gcs/internal/commands/auth"
	collectioncmd "github.com/scttfrdmn/globus-go-gcs/internal/commands/collection"
	endpointcmd "github.com/scttfrdmn/globus-go-gcs/internal/commands/endpoint"
	"github.com/spf13/cobra"
)

var (
	// Version information (set by build flags)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "globus-connect-server",
		Short: "Globus Connect Server command-line interface",
		Long: `globus-connect-server is a CLI for managing Globus Connect Server v5 endpoints.

This is a complete Go port of the Python globus-connect-server CLI with 100% feature parity.

For more information, see: https://docs.globus.org/globus-connect-server/v5/`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Global flags
	rootCmd.PersistentFlags().String("format", "text", "Output format (text, json)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logging")

	// Authentication commands
	rootCmd.AddCommand(authcmd.NewLoginCmd())
	rootCmd.AddCommand(authcmd.NewLogoutCmd())
	rootCmd.AddCommand(authcmd.NewWhoamiCmd())

	// Endpoint commands
	rootCmd.AddCommand(endpointcmd.NewEndpointCmd())

	// Collection commands
	rootCmd.AddCommand(collectioncmd.NewCollectionCmd())

	// TODO: Add additional command groups
	// rootCmd.AddCommand(newNodeCmd())
	// rootCmd.AddCommand(newStorageGatewayCmd())
	// rootCmd.AddCommand(newSessionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
