package auth

import (
	"bytes"
	"context"
	"testing"

	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
)

func TestNewWhoamiCmd(t *testing.T) {
	cmd := NewWhoamiCmd()

	if cmd == nil {
		t.Fatal("NewWhoamiCmd() returned nil")
	}

	if cmd.Use != "whoami" {
		t.Errorf("NewWhoamiCmd() Use = %q, want %q", cmd.Use, "whoami")
	}

	if cmd.Short == "" {
		t.Error("NewWhoamiCmd() Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("NewWhoamiCmd() Long description is empty")
	}

	if cmd.RunE == nil {
		t.Error("NewWhoamiCmd() RunE is nil")
	}
}

func TestNewWhoamiCmd_Flags(t *testing.T) {
	cmd := NewWhoamiCmd()

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

			if flag.DefValue != tt.defaultValue {
				t.Errorf("flag %q default = %q, want %q", tt.flagName, flag.DefValue, tt.defaultValue)
			}
		})
	}
}

func TestRunWhoami_NoToken(t *testing.T) {
	ctx := context.Background()
	buf := &bytes.Buffer{}

	// Test with a profile that doesn't exist
	err := runWhoami(ctx, "nonexistent-profile-test", "text", buf)
	if err == nil {
		t.Error("runWhoami() expected error for nonexistent profile, got nil")
	}

	// Output buffer should be empty on error
	if buf.Len() > 0 {
		t.Errorf("runWhoami() wrote to buffer on error: %q", buf.String())
	}
}
