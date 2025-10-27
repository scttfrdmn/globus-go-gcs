package role

import (
	"testing"

	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
)

func TestNewCreateCmd(t *testing.T) {
	cmd := NewCreateCmd()

	if cmd == nil {
		t.Fatal("NewCreateCmd() returned nil")
	}

	if cmd.Use != "create" {
		t.Errorf("NewCreateCmd() Use = %q, want %q", cmd.Use, "create")
	}

	if cmd.Short == "" {
		t.Error("NewCreateCmd() Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("NewCreateCmd() Long description is empty")
	}

	if cmd.RunE == nil {
		t.Error("NewCreateCmd() RunE is nil")
	}
}

func TestNewCreateCmd_Flags(t *testing.T) {
	cmd := NewCreateCmd()

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
		},
		{
			name:         "format flag",
			flagName:     "format",
			shorthand:    "f",
			defaultValue: "text",
		},
		{
			name:     "endpoint flag",
			flagName: "endpoint",
			required: true,
		},
		{
			name:      "collection flag",
			flagName:  "collection",
			shorthand: "c",
			required:  true,
		},
		{
			name:     "principal flag",
			flagName: "principal",
			required: true,
		},
		{
			name:     "role flag",
			flagName: "role",
			required: true,
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
