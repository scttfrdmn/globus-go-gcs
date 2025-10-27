package role

import (
	"testing"
)

func TestNewDeleteCmd(t *testing.T) {
	cmd := NewDeleteCmd()

	if cmd == nil {
		t.Fatal("NewDeleteCmd() returned nil")
	}

	if cmd.Use != "delete ROLE_ID" {
		t.Errorf("NewDeleteCmd() Use = %q, want %q", cmd.Use, "delete ROLE_ID")
	}

	if cmd.Short == "" {
		t.Error("NewDeleteCmd() Short description is empty")
	}

	if cmd.RunE == nil {
		t.Error("NewDeleteCmd() RunE is nil")
	}

	if cmd.Args == nil {
		t.Error("NewDeleteCmd() Args validation is nil")
	}
}

func TestNewDeleteCmd_Flags(t *testing.T) {
	cmd := NewDeleteCmd()

	tests := []struct {
		name      string
		flagName  string
		shorthand string
	}{
		{
			name:      "profile flag",
			flagName:  "profile",
			shorthand: "p",
		},
		{
			name:      "format flag",
			flagName:  "format",
			shorthand: "f",
		},
		{
			name:     "endpoint flag",
			flagName: "endpoint",
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
		})
	}
}
