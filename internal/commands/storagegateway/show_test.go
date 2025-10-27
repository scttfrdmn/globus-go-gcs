package storagegateway

import (
	"bytes"
	"context"
	"testing"

	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
)

func TestNewShowCmd(t *testing.T) {
	cmd := NewShowCmd()

	if cmd == nil {
		t.Fatal("NewShowCmd() returned nil")
	}

	if cmd.Use != "show GATEWAY_ID" {
		t.Errorf("NewShowCmd() Use = %q, want %q", cmd.Use, "show GATEWAY_ID")
	}

	if cmd.Short == "" {
		t.Error("NewShowCmd() Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("NewShowCmd() Long description is empty")
	}

	if cmd.RunE == nil {
		t.Error("NewShowCmd() RunE is nil")
	}

	if cmd.Args == nil {
		t.Error("NewShowCmd() Args validation is nil")
	}
}

func TestNewShowCmd_Flags(t *testing.T) {
	cmd := NewShowCmd()

	tests := []struct {
		name         string
		flagName     string
		shorthand    string
		defaultValue string
	}{
		{
			name:         "profile flag",
			flagName:     "profile",
			shorthand:    "p",
			defaultValue: config.DefaultProfile,
		},
		{
			name:         "format flag",
			flagName:     "format",
			shorthand:    "f",
			defaultValue: "text",
		},
		{
			name:      "endpoint flag",
			flagName:  "endpoint",
			shorthand: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := cmd.Flags().Lookup(tt.flagName)
			if flag == nil {
				t.Fatalf("flag %q not found", tt.flagName)
			}

			if flag.Shorthand != tt.shorthand {
				t.Errorf("flag %q shorthand = %q, want %q", tt.flagName, flag.Shorthand, tt.shorthand)
			}

			if tt.defaultValue != "" && flag.DefValue != tt.defaultValue {
				t.Errorf("flag %q default = %q, want %q", tt.flagName, flag.DefValue, tt.defaultValue)
			}
		})
	}
}

func TestRunShow_NoToken(t *testing.T) {
	ctx := context.Background()
	buf := &bytes.Buffer{}

	// Test with a profile that doesn't exist
	err := runShow(ctx, "nonexistent-profile-test", "text", "test.example.org", "test-gateway-id", buf)
	if err == nil {
		t.Error("runShow() expected error for nonexistent profile, got nil")
	}

	// Output buffer should be empty on error
	if buf.Len() > 0 {
		t.Errorf("runShow() wrote to buffer on error: %q", buf.String())
	}
}
