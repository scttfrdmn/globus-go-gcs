package role

import (
	"bytes"
	"context"
	"testing"

	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
)

func TestNewListCmd(t *testing.T) {
	cmd := NewListCmd()

	if cmd == nil {
		t.Fatal("NewListCmd() returned nil")
	}

	if cmd.Use != "list" {
		t.Errorf("NewListCmd() Use = %q, want %q", cmd.Use, "list")
	}

	if cmd.Short == "" {
		t.Error("NewListCmd() Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("NewListCmd() Long description is empty")
	}

	if cmd.RunE == nil {
		t.Error("NewListCmd() RunE is nil")
	}
}

func TestNewListCmd_Flags(t *testing.T) {
	cmd := NewListCmd()

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
		{
			name:      "collection flag",
			flagName:  "collection",
			shorthand: "c",
		},
		{
			name:      "principal flag",
			flagName:  "principal",
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

func TestRunList_NoToken(t *testing.T) {
	ctx := context.Background()
	buf := &bytes.Buffer{}

	// Test with a profile that doesn't exist
	err := runList(ctx, "nonexistent-profile-test", "text", "test.example.org", "", "", buf)
	if err == nil {
		t.Error("runList() expected error for nonexistent profile, got nil")
	}

	// Output buffer should be empty on error
	if buf.Len() > 0 {
		t.Errorf("runList() wrote to buffer on error: %q", buf.String())
	}
}
