// Package secureinput provides secure methods for reading sensitive data (passwords, secrets).
//
// This package implements NIST 800-53 IA-5(7) compliance by preventing secrets
// from appearing in:
//   - Process listings (ps, /proc)
//   - Shell history
//   - Log files
//   - Command-line argument capture
//
// Three secure input methods are provided:
//  1. Interactive prompt (default, recommended)
//  2. Read from stdin (--secret-stdin)
//  3. Read from environment variable (--secret-env)
package secureinput

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// ReadSecretOptions configures how to read a secret value.
type ReadSecretOptions struct {
	// PromptMessage is displayed when prompting interactively
	PromptMessage string

	// UseStdin indicates to read from stdin instead of prompting
	UseStdin bool

	// EnvVar is the environment variable name to read from
	EnvVar string

	// AllowEmpty allows empty/blank secrets (default: false)
	AllowEmpty bool
}

// ReadSecret reads a secret value using one of three secure methods.
//
// Priority order:
//  1. If EnvVar is set, read from that environment variable
//  2. If UseStdin is true, read from stdin
//  3. Otherwise, prompt interactively (with hidden input)
//
// Returns an error if:
//   - The secret is empty (unless AllowEmpty is true)
//   - The environment variable doesn't exist
//   - Reading from stdin/terminal fails
func ReadSecret(opts ReadSecretOptions) (string, error) {
	var secret string
	var err error

	// Select input method based on options
	switch {
	case opts.EnvVar != "":
		// Method 1: Read from environment variable
		secret, err = readFromEnv(opts.EnvVar)
		if err != nil {
			return "", err
		}
	case opts.UseStdin:
		// Method 2: Read from stdin
		secret, err = readFromStdin()
		if err != nil {
			return "", err
		}
	default:
		// Method 3: Interactive prompt (default, most secure)
		secret, err = readFromPrompt(opts.PromptMessage)
		if err != nil {
			return "", err
		}
	}

	// Trim whitespace
	secret = strings.TrimSpace(secret)

	// Validate not empty (unless explicitly allowed)
	if !opts.AllowEmpty && secret == "" {
		return "", fmt.Errorf("secret cannot be empty")
	}

	return secret, nil
}

// readFromEnv reads a secret from an environment variable.
func readFromEnv(envVar string) (string, error) {
	value, exists := os.LookupEnv(envVar)
	if !exists {
		return "", fmt.Errorf("environment variable %q not set", envVar)
	}

	return value, nil
}

// readFromStdin reads a secret from stdin (one line).
//
// This is used with --secret-stdin flag, typically in scripts:
//
//	echo "my-secret" | globus-connect-server ... --secret-stdin
func readFromStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("read from stdin: %w", err)
	}

	return strings.TrimSpace(line), nil
}

// readFromPrompt reads a secret interactively with hidden input.
//
// This is the most secure method as it:
//   - Doesn't echo characters to the terminal
//   - Doesn't appear in shell history
//   - Doesn't appear in process listings
//   - Provides clear user feedback
func readFromPrompt(promptMessage string) (string, error) {
	// Default prompt if none provided
	if promptMessage == "" {
		promptMessage = "Enter secret"
	}

	// Print prompt
	fmt.Fprintf(os.Stderr, "%s: ", promptMessage)

	// Read with hidden input (no echo)
	secretBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("read password: %w", err)
	}

	// Print newline after hidden input
	fmt.Fprintln(os.Stderr)

	return string(secretBytes), nil
}

// ValidateSecret performs basic validation on a secret value.
//
// Returns an error if the secret:
//   - Is empty or only whitespace
//   - Is too short (less than minLength)
//   - Is too long (more than maxLength, if specified)
func ValidateSecret(secret string, minLength, maxLength int) error {
	if strings.TrimSpace(secret) == "" {
		return fmt.Errorf("secret cannot be empty")
	}

	if len(secret) < minLength {
		return fmt.Errorf("secret too short (min %d characters)", minLength)
	}

	if maxLength > 0 && len(secret) > maxLength {
		return fmt.Errorf("secret too long (max %d characters)", maxLength)
	}

	return nil
}

// SecureString is a string wrapper that prevents accidental logging/printing of secrets.
//
// Use this to wrap secret values to prevent them from being accidentally
// exposed in logs, error messages, or debug output.
type SecureString struct {
	value string
}

// NewSecureString creates a SecureString from a plain string.
func NewSecureString(value string) SecureString {
	return SecureString{value: value}
}

// String implements fmt.Stringer and returns a redacted value.
// This prevents accidental exposure in logs/debug output.
func (s SecureString) String() string {
	if s.value == "" {
		return "(empty)"
	}
	return "(redacted)"
}

// Value returns the actual secret value.
// Only call this when you actually need to use the secret.
func (s SecureString) Value() string {
	return s.value
}

// Clear zeros out the secret value in memory.
// Call this when done with the secret to reduce exposure time.
func (s *SecureString) Clear() {
	// Overwrite with zeros
	if s.value != "" {
		b := []byte(s.value)
		for i := range b {
			b[i] = 0
		}
		s.value = ""
	}
}
