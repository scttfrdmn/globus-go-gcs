package endpoint

import (
	"testing"
)

func TestNewUpdateCmd(t *testing.T) {
	cmd := NewUpdateCmd()

	if cmd == nil {
		t.Fatal("NewUpdateCmd() returned nil")
	}

	if cmd.Use != "update" {
		t.Errorf("NewUpdateCmd() Use = %q, want %q", cmd.Use, "update")
	}

	if cmd.Short == "" {
		t.Error("NewUpdateCmd() Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("NewUpdateCmd() Long description is empty")
	}

	if cmd.RunE == nil {
		t.Error("NewUpdateCmd() RunE is nil")
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
		{
			name:     "organization flag",
			flagName: "organization",
		},
		{
			name:     "department flag",
			flagName: "department",
		},
		{
			name:     "description flag",
			flagName: "description",
		},
		{
			name:     "contact-email flag",
			flagName: "contact-email",
		},
		{
			name:     "keywords flag",
			flagName: "keywords",
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
