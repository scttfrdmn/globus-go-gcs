package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewQueryCmd creates the audit query command.
func NewQueryCmd() *cobra.Command {
	var (
		format     string
		startTime  string
		endTime    string
		eventType  string
		identityID string
		action     string
		result     string
		limit      int
	)

	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query audit logs from local database",
		Long: `Query audit logs from the local SQLite database.

This command searches the local audit database with various filters.
Use 'audit load' first to populate the database.

Example:
  # Query all logs
  globus-connect-server audit query

  # Query with filters
  globus-connect-server audit query \
    --event-type transfer \
    --result success \
    --start-time "2025-01-01T00:00:00Z"

  # Limit results
  globus-connect-server audit query --limit 50`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runQuery(cmd.Context(), format, startTime, endTime, eventType,
				identityID, action, result, limit, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&startTime, "start-time", "", "Start time (RFC3339 format)")
	cmd.Flags().StringVar(&endTime, "end-time", "", "End time (RFC3339 format)")
	cmd.Flags().StringVar(&eventType, "event-type", "", "Filter by event type")
	cmd.Flags().StringVar(&identityID, "identity", "", "Filter by identity ID")
	cmd.Flags().StringVar(&action, "action", "", "Filter by action")
	cmd.Flags().StringVar(&result, "result", "", "Filter by result (success, failure)")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum number of results")

	return cmd
}

// runQuery executes the audit query command.
func runQuery(ctx context.Context, formatStr, startTimeStr, endTimeStr, eventType,
	identityID, action, result string, limit int, out interface{ Write([]byte) (int, error) }) error {
	// Create output formatter
	formatter := output.NewFormatter(output.Format(formatStr), out)

	// Parse time parameters
	var startTime, endTime *time.Time
	if startTimeStr != "" {
		t, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return fmt.Errorf("invalid start time: %w", err)
		}
		startTime = &t
	}
	if endTimeStr != "" {
		t, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return fmt.Errorf("invalid end time: %w", err)
		}
		endTime = &t
	}

	// Initialize database
	dbPath, err := getAuditDBPath()
	if err != nil {
		return fmt.Errorf("get database path: %w", err)
	}

	db, err := initAuditDB(dbPath)
	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}
	defer func() { _ = db.Close() }()

	// Build query
	query := "SELECT id, timestamp, event_type, identity_id, username, resource, resource_id, action, result, message, client_ip, metadata FROM audit_logs WHERE 1=1"
	args := []interface{}{}

	if startTime != nil {
		query += " AND timestamp >= ?"
		args = append(args, startTime)
	}
	if endTime != nil {
		query += " AND timestamp <= ?"
		args = append(args, endTime)
	}
	if eventType != "" {
		query += " AND event_type = ?"
		args = append(args, eventType)
	}
	if identityID != "" {
		query += " AND identity_id = ?"
		args = append(args, identityID)
	}
	if action != "" {
		query += " AND action = ?"
		args = append(args, action)
	}
	if result != "" {
		query += " AND result = ?"
		args = append(args, result)
	}

	query += " ORDER BY timestamp DESC LIMIT ?"
	args = append(args, limit)

	// Execute query
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query database: %w", err)
	}
	defer func() { _ = rows.Close() }()

	// Fetch results
	logs, err := scanAuditLogs(rows)
	if err != nil {
		return err
	}

	// Output results
	return formatQueryResults(formatter, logs)
}

// scanAuditLogs scans query results into audit log entries.
func scanAuditLogs(rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}) ([]gcs.AuditLog, error) {
	var logs []gcs.AuditLog
	for rows.Next() {
		var log gcs.AuditLog
		var timestamp string
		var metadataJSON string

		err := rows.Scan(
			&log.ID,
			&timestamp,
			&log.EventType,
			&log.IdentityID,
			&log.Username,
			&log.Resource,
			&log.ResourceID,
			&log.Action,
			&log.Result,
			&log.Message,
			&log.ClientIP,
			&metadataJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		// Parse timestamp
		log.Timestamp, _ = time.Parse(time.RFC3339, timestamp)

		// Parse metadata JSON
		if metadataJSON != "" {
			_ = json.Unmarshal([]byte(metadataJSON), &log.Metadata)
		}

		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return logs, nil
}

// formatQueryResults formats query results for output.
func formatQueryResults(formatter *output.Formatter, logs []gcs.AuditLog) error {
	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(map[string]interface{}{
			"count": len(logs),
			"logs":  logs,
		})
	}

	// Text format
	if len(logs) == 0 {
		return formatter.Println("No audit logs found matching the criteria")
	}

	if err := formatter.PrintText("Found %d audit log entries\n", len(logs)); err != nil {
		return err
	}
	if err := formatter.Println(strings.Repeat("=", 80)); err != nil {
		return err
	}
	if err := formatter.Println(); err != nil {
		return err
	}

	for _, log := range logs {
		if err := formatLogEntry(formatter, log); err != nil {
			return err
		}
	}

	return nil
}

// formatLogEntry formats a single log entry for text output.
func formatLogEntry(formatter *output.Formatter, log gcs.AuditLog) error {
	if err := formatter.PrintText("%-15s %s\n", "Timestamp:", log.Timestamp.Format(time.RFC3339)); err != nil {
		return err
	}
	if log.EventType != "" {
		if err := formatter.PrintText("%-15s %s\n", "Event Type:", log.EventType); err != nil {
			return err
		}
	}
	if log.Username != "" {
		if err := formatter.PrintText("%-15s %s\n", "Username:", log.Username); err != nil {
			return err
		}
	}
	if log.Resource != "" {
		if err := formatter.PrintText("%-15s %s\n", "Resource:", log.Resource); err != nil {
			return err
		}
	}
	if log.Action != "" {
		if err := formatter.PrintText("%-15s %s\n", "Action:", log.Action); err != nil {
			return err
		}
	}
	if log.Result != "" {
		if err := formatter.PrintText("%-15s %s\n", "Result:", log.Result); err != nil {
			return err
		}
	}
	if log.Message != "" {
		if err := formatter.PrintText("%-15s %s\n", "Message:", log.Message); err != nil {
			return err
		}
	}
	return formatter.Println()
}
