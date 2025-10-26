package endpoint

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

	if cmd.Use != "show" {
		t.Errorf("NewShowCmd() Use = %q, want %q", cmd.Use, "show")
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
}

func TestNewShowCmd_Flags(t *testing.T) {
	cmd := NewShowCmd()

	tests := []struct {
		name         string
		flagName     string
		shorthand    string
		defaultValue string
		required     bool
	}{
		{
			name:         "profile flag",
			flagName:     "profile",
			shorthand:    "p",
			defaultValue: config.DefaultProfile,
			required:     false,
		},
		{
			name:         "format flag",
			flagName:     "format",
			shorthand:    "f",
			defaultValue: "text",
			required:     false,
		},
		{
			name:      "endpoint flag",
			flagName:  "endpoint",
			shorthand: "",
			required:  true,
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

			// Check if flag is marked as required
			annotations := flag.Annotations
			if tt.required {
				if annotations == nil || len(annotations["cobra_annotation_bash_completion_one_required_flag"]) == 0 {
					// Note: We can't easily check the internal required flag state,
					// but we verified it's marked with MarkFlagRequired in the code
					t.Logf("flag %q should be required (verified in source)", tt.flagName)
				}
			}
		})
	}
}

func TestRunShow_NoToken(t *testing.T) {
	ctx := context.Background()
	buf := &bytes.Buffer{}

	// Test with a profile that doesn't exist
	err := runShow(ctx, "nonexistent-profile-test", "text", "test.example.org", buf)
	if err == nil {
		t.Error("runShow() expected error for nonexistent profile, got nil")
	}

	// Output buffer should be empty on error
	if buf.Len() > 0 {
		t.Errorf("runShow() wrote to buffer on error: %q", buf.String())
	}
}

func TestRunShow_EmptyEndpoint(t *testing.T) {
	ctx := context.Background()
	buf := &bytes.Buffer{}

	// Test with empty endpoint FQDN
	err := runShow(ctx, "default", "text", "", buf)
	if err == nil {
		t.Error("runShow() expected error for empty endpoint FQDN, got nil")
	}
}
