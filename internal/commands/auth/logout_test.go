package auth

import (
	"testing"

	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
)

func TestNewLogoutCmd(t *testing.T) {
	cmd := NewLogoutCmd()

	if cmd == nil {
		t.Fatal("NewLogoutCmd() returned nil")
	}

	if cmd.Use != "logout" {
		t.Errorf("NewLogoutCmd() Use = %q, want %q", cmd.Use, "logout")
	}

	if cmd.Short == "" {
		t.Error("NewLogoutCmd() Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("NewLogoutCmd() Long description is empty")
	}

	if cmd.RunE == nil {
		t.Error("NewLogoutCmd() RunE is nil")
	}
}

func TestNewLogoutCmd_Flags(t *testing.T) {
	cmd := NewLogoutCmd()

	flag := cmd.Flags().Lookup("profile")
	if flag == nil {
		t.Fatal("profile flag not found")
	}

	if flag.Shorthand != "p" {
		t.Errorf("profile flag shorthand = %q, want %q", flag.Shorthand, "p")
	}

	if flag.DefValue != config.DefaultProfile {
		t.Errorf("profile flag default = %q, want %q", flag.DefValue, config.DefaultProfile)
	}
}

func TestRunLogout_NoToken(t *testing.T) {
	// Test with a profile that doesn't exist
	err := runLogout("nonexistent-profile-test")
	if err == nil {
		t.Error("runLogout() expected error for nonexistent profile, got nil")
	}
}
