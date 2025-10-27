// Package audit provides commands for managing audit logs.
package audit

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	_ "modernc.org/sqlite" // SQLite driver
)

// NewAuditCmd creates the audit command with subcommands.
func NewAuditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Manage audit logs",
		Long: `Commands for managing audit logs.

Audit logs track all activities on the endpoint including transfers,
access events, and authentication. Logs can be loaded into a local
SQLite database for searching and analysis.`,
	}

	// Add subcommands
	cmd.AddCommand(NewLoadCmd())
	cmd.AddCommand(NewQueryCmd())
	cmd.AddCommand(NewDumpCmd())

	return cmd
}

// getAuditDBPath returns the path to the audit database.
func getAuditDBPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}

	dbDir := filepath.Join(homeDir, ".globus-connect-server", "audit")
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return "", fmt.Errorf("create audit directory: %w", err)
	}

	return filepath.Join(dbDir, "audit.db"), nil
}

// initAuditDB initializes the audit database.
func initAuditDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Create audit_logs table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS audit_logs (
		id TEXT PRIMARY KEY,
		timestamp DATETIME NOT NULL,
		event_type TEXT,
		identity_id TEXT,
		username TEXT,
		resource TEXT,
		resource_id TEXT,
		action TEXT,
		result TEXT,
		message TEXT,
		client_ip TEXT,
		metadata TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON audit_logs(timestamp);
	CREATE INDEX IF NOT EXISTS idx_event_type ON audit_logs(event_type);
	CREATE INDEX IF NOT EXISTS idx_identity ON audit_logs(identity_id);
	CREATE INDEX IF NOT EXISTS idx_result ON audit_logs(result);
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create table: %w", err)
	}

	return db, nil
}
