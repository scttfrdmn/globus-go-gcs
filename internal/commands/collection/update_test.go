package collection

import (
	"testing"
)

func TestNewUpdateCmd(t *testing.T) {
	cmd := NewUpdateCmd()

	if cmd == nil {
		t.Fatal("NewUpdateCmd() returned nil")
	}

	if cmd.Use != "update COLLECTION_ID" {
		t.Errorf("NewUpdateCmd() Use = %q, want %q", cmd.Use, "update COLLECTION_ID")
	}

	if cmd.Short == "" {
		t.Error("NewUpdateCmd() Short description is empty")
	}

	if cmd.RunE == nil {
		t.Error("NewUpdateCmd() RunE is nil")
	}

	if cmd.Args == nil {
		t.Error("NewUpdateCmd() Args validation is nil")
	}
}

func TestNewUpdateCmd_Flags(t *testing.T) {
	cmd := NewUpdateCmd()

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
		{
			name:     "display-name flag",
			flagName: "display-name",
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
