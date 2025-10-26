package auth

import (
	"testing"

	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
)

func TestNewLoginCmd(t *testing.T) {
	cmd := NewLoginCmd()

	if cmd == nil {
		t.Fatal("NewLoginCmd() returned nil")
	}

	if cmd.Use != "login" {
		t.Errorf("NewLoginCmd() Use = %q, want %q", cmd.Use, "login")
	}

	if cmd.Short == "" {
		t.Error("NewLoginCmd() Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("NewLoginCmd() Long description is empty")
	}

	if cmd.RunE == nil {
		t.Error("NewLoginCmd() RunE is nil")
	}
}

func TestNewLoginCmd_Flags(t *testing.T) {
	cmd := NewLoginCmd()

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
			name:         "scopes flag",
			flagName:     "scopes",
			shorthand:    "",
			defaultValue: defaultScopes,
		},
		{
			name:      "no-local-server flag",
			flagName:  "no-local-server",
			shorthand: "",
			// Bool flags have "false" as default
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

func TestGenerateState(t *testing.T) {
	state1 := generateState()
	state2 := generateState()

	if state1 == "" {
		t.Error("generateState() returned empty string")
	}

	if state1 == state2 {
		t.Error("generateState() returned same value twice (should be unique)")
	}

	// Check that it has the expected prefix
	const expectedPrefix = "gcs-cli-"
	if len(state1) < len(expectedPrefix) {
		t.Errorf("generateState() returned value too short: %q", state1)
	}

	actualPrefix := state1[:len(expectedPrefix)]
	if actualPrefix != expectedPrefix {
		t.Errorf("generateState() prefix = %q, want %q", actualPrefix, expectedPrefix)
	}
}
