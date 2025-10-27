package audit

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewDumpCmd creates the audit dump command.
func NewDumpCmd() *cobra.Command {
	var (
		format     string
		outputFile string
		startTime  string
		endTime    string
		eventType  string
		identityID string
		action     string
		result     string
	)

	cmd := &cobra.Command{
		Use:   "dump",
		Short: "Export audit logs to file",
		Long: `Export audit logs from the local database to a file.

Supports JSON and CSV export formats. You can apply filters to export
specific subsets of logs.

Example:
  # Export to JSON
  globus-connect-server audit dump \
    --output audit-logs.json \
    --format json

  # Export to CSV with filters
  globus-connect-server audit dump \
    --output audit-logs.csv \
    --format csv \
    --event-type transfer \
    --result success`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDump(cmd.Context(), format, outputFile, startTime, endTime,
				eventType, identityID, action, result, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "json", "Export format (json, csv)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	cmd.Flags().StringVar(&startTime, "start-time", "", "Start time (RFC3339 format)")
	cmd.Flags().StringVar(&endTime, "end-time", "", "End time (RFC3339 format)")
	cmd.Flags().StringVar(&eventType, "event-type", "", "Filter by event type")
	cmd.Flags().StringVar(&identityID, "identity", "", "Filter by identity ID")
	cmd.Flags().StringVar(&action, "action", "", "Filter by action")
	cmd.Flags().StringVar(&result, "result", "", "Filter by result (success, failure)")

	_ = cmd.MarkFlagRequired("output")

	return cmd
}

// runDump executes the audit dump command.
func runDump(ctx context.Context, format, outputFile, startTimeStr, endTimeStr, eventType,
	identityID, action, result string, out interface{ Write([]byte) (int, error) }) error {
	// Create output formatter
	formatter := output.NewFormatter(output.Format("text"), out)

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

	query += " ORDER BY timestamp DESC"

	// Execute query
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query database: %w", err)
	}
	defer func() { _ = rows.Close() }()

	// Fetch results
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
			return fmt.Errorf("scan row: %w", err)
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
		return fmt.Errorf("iterate rows: %w", err)
	}

	// Create output file
	// #nosec G304 - outputFile path is user-provided via flag, which is expected behavior
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Export based on format
	switch format {
	case "json":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(logs); err != nil {
			return fmt.Errorf("encode JSON: %w", err)
		}

	case "csv":
		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write header
		header := []string{
			"ID", "Timestamp", "EventType", "IdentityID", "Username",
			"Resource", "ResourceID", "Action", "Result", "Message", "ClientIP",
		}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("write CSV header: %w", err)
		}

		// Write rows
		for _, log := range logs {
			row := []string{
				log.ID,
				log.Timestamp.Format(time.RFC3339),
				log.EventType,
				log.IdentityID,
				log.Username,
				log.Resource,
				log.ResourceID,
				log.Action,
				log.Result,
				log.Message,
				log.ClientIP,
			}
			if err := writer.Write(row); err != nil {
				return fmt.Errorf("write CSV row: %w", err)
			}
		}

	default:
		return fmt.Errorf("unsupported format: %s (use json or csv)", format)
	}

	// Output success message
	if err := formatter.PrintText("Exported %d audit log entries to %s\n", len(logs), outputFile); err != nil {
		return err
	}

	return nil
}
