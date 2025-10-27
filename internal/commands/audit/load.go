package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-gcs/pkg/gcs"
	"github.com/scttfrdmn/globus-go-gcs/pkg/output"
	"github.com/spf13/cobra"
)

// NewLoadCmd creates the audit load command.
func NewLoadCmd() *cobra.Command {
	var (
		profile      string
		format       string
		endpointFQDN string
		startTime    string
		endTime      string
		eventType    string
		limit        int
	)

	cmd := &cobra.Command{
		Use:   "load",
		Short: "Load audit logs into local database",
		Long: `Load audit logs from the GCS Manager API into a local SQLite database.

The local database enables fast searching and filtering of audit logs.
You can specify time ranges and filters to load specific subsets of logs.

Example:
  # Load last 24 hours of logs
  globus-connect-server audit load \
    --endpoint example.data.globus.org

  # Load logs with filters
  globus-connect-server audit load \
    --endpoint example.data.globus.org \
    --start-time "2025-01-01T00:00:00Z" \
    --end-time "2025-01-02T00:00:00Z" \
    --event-type transfer

Requires an active authentication session (use 'login' first).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runLoad(cmd.Context(), profile, format, endpointFQDN, startTime, endTime,
				eventType, limit, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&endpointFQDN, "endpoint", "", "Endpoint FQDN (e.g., abc.def.data.globus.org)")
	cmd.Flags().StringVar(&startTime, "start-time", "", "Start time (RFC3339 format)")
	cmd.Flags().StringVar(&endTime, "end-time", "", "End time (RFC3339 format)")
	cmd.Flags().StringVar(&eventType, "event-type", "", "Filter by event type")
	cmd.Flags().IntVar(&limit, "limit", 1000, "Maximum number of logs to load")

	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

// runLoad executes the audit load command.
func runLoad(ctx context.Context, profile, formatStr, endpointFQDN, startTimeStr, endTimeStr,
	eventType string, limit int, out interface{ Write([]byte) (int, error) }) error {
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

	// Create GCS client
	gcsClient, err := gcs.NewClient(
		endpointFQDN,
		gcs.WithAccessToken(token.AccessToken),
	)
	if err != nil {
		return fmt.Errorf("create GCS client: %w", err)
	}

	// Build query parameters
	params := &gcs.AuditQueryParams{
		StartTime: startTime,
		EndTime:   endTime,
		EventType: eventType,
		Limit:     limit,
	}

	// Fetch audit logs from API
	logs, err := gcsClient.GetAuditLogs(ctx, params)
	if err != nil {
		return fmt.Errorf("fetch audit logs: %w", err)
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

	// Insert logs into database
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.Prepare(`
		INSERT OR REPLACE INTO audit_logs
		(id, timestamp, event_type, identity_id, username, resource, resource_id,
		 action, result, message, client_ip, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	loaded := 0
	for _, log := range logs.Data {
		metadataJSON, _ := json.Marshal(log.Metadata)

		_, err := stmt.Exec(
			log.ID,
			log.Timestamp,
			log.EventType,
			log.IdentityID,
			log.Username,
			log.Resource,
			log.ResourceID,
			log.Action,
			log.Result,
			log.Message,
			log.ClientIP,
			string(metadataJSON),
		)
		if err != nil {
			return fmt.Errorf("insert log: %w", err)
		}
		loaded++
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	// Output based on format
	if formatter.IsJSON() {
		return formatter.PrintJSON(map[string]interface{}{
			"loaded": loaded,
			"database": dbPath,
		})
	}

	// Text format
	if err := formatter.PrintText("Loaded %d audit log entries into database\n", loaded); err != nil {
		return err
	}
	if err := formatter.PrintText("Database: %s\n", dbPath); err != nil {
		return err
	}

	return nil
}
