// Package auth provides authentication commands for the GCS CLI.
package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/scttfrdmn/globus-go-gcs/internal/auth"
	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	globusauth "github.com/scttfrdmn/globus-go-sdk/v3/pkg/services/auth"
	"github.com/spf13/cobra"
)

const (
	// Default OAuth2 scopes for GCS CLI
	defaultScopes = "openid profile email " +
		"urn:globus:auth:scope:auth.globus.org:view_identities " +
		"urn:globus:auth:scope:transfer.api.globus.org:all"

	// Local callback server settings
	callbackPort = "8080"
	callbackPath = "/callback"
)

// NewLoginCmd creates the login command.
func NewLoginCmd() *cobra.Command {
	var (
		profile string
		scopes  string
		noLocal bool
	)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Globus Auth",
		Long: `Authenticate with Globus Auth using OAuth2.

This command opens your browser to authenticate with Globus Auth.
After authentication, your tokens are securely stored locally.

The tokens are stored in: ~/.globus-connect-server/tokens/<profile>.json

By default, uses a local callback server to receive the OAuth code.
Use --no-local-server to manually copy/paste the authorization code.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runLogin(cmd.Context(), profile, scopes, noLocal)
		},
	}

	cmd.Flags().StringVarP(&profile, "profile", "p", config.DefaultProfile, "Profile name")
	cmd.Flags().StringVar(&scopes, "scopes", defaultScopes, "OAuth2 scopes (space-separated)")
	cmd.Flags().BoolVar(&noLocal, "no-local-server", false, "Disable local callback server (manual code entry)")

	return cmd
}

// runLogin executes the login flow.
func runLogin(ctx context.Context, profile, scopes string, noLocal bool) error {
	// Load client configuration
	cfg, err := config.LoadClientConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Create auth client
	authClient, err := globusauth.NewClient(
		globusauth.WithClientID(cfg.ClientID),
		globusauth.WithClientSecret(cfg.ClientSecret),
	)
	if err != nil {
		return fmt.Errorf("create auth client: %w", err)
	}

	// Set redirect URI
	redirectURI := fmt.Sprintf("http://localhost:%s%s", callbackPort, callbackPath)
	authClient.RedirectURL = redirectURI

	// Generate state for CSRF protection
	state := generateState()

	// Get authorization URL
	authURL := authClient.GetAuthorizationURL(state, scopes)

	fmt.Println("Please authenticate by visiting this URL:")
	fmt.Println()
	fmt.Println(authURL)
	fmt.Println()

	// Get authorization code
	var code string
	if noLocal {
		code, err = getCodeManual()
	} else {
		code, err = getCodeViaCallback(ctx, state, callbackPort, callbackPath)
	}
	if err != nil {
		return fmt.Errorf("get authorization code: %w", err)
	}

	// Exchange code for tokens
	fmt.Println("Exchanging authorization code for tokens...")
	tokenResp, err := authClient.ExchangeAuthorizationCode(ctx, code)
	if err != nil {
		return fmt.Errorf("exchange code: %w", err)
	}

	// Save tokens
	tokenInfo := auth.TokenFromAuthResponse(tokenResp)
	if err := auth.SaveToken(profile, tokenInfo); err != nil {
		return fmt.Errorf("save token: %w", err)
	}

	fmt.Println("✓ Login successful!")
	fmt.Printf("Profile: %s\n", profile)
	fmt.Printf("Token expires: %s\n", tokenInfo.ExpiresAt.Format(time.RFC3339))

	return nil
}

// getCodeManual prompts the user to manually enter the authorization code.
func getCodeManual() (string, error) {
	fmt.Print("Enter authorization code: ")
	var code string
	if _, err := fmt.Scanln(&code); err != nil {
		return "", fmt.Errorf("read code: %w", err)
	}

	return code, nil
}

// getCodeViaCallback starts a local HTTP server to receive the OAuth callback.
func getCodeViaCallback(ctx context.Context, expectedState, port, path string) (string, error) {
	// Create channel to receive code
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Create HTTP handler
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		// Check state for CSRF protection
		state := r.URL.Query().Get("state")
		if state != expectedState {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			errChan <- fmt.Errorf("invalid state parameter")
			return
		}

		// Get authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "No code in response", http.StatusBadRequest)
			errChan <- fmt.Errorf("no code in response")
			return
		}

		// Send success response
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintf(w, `
<html>
<head><title>Authentication Successful</title></head>
<body>
<h1>✓ Authentication Successful</h1>
<p>You have successfully authenticated with Globus.</p>
<p>You may close this window and return to the CLI.</p>
</body>
</html>
`)

		// Send code to channel
		codeChan <- code
	})

	// Create server
	server := &http.Server{
		Addr:              "localhost:" + port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("callback server: %w", err)
		}
	}()

	fmt.Printf("Waiting for authentication on http://localhost:%s%s\n", port, path)
	fmt.Println("(If browser doesn't open automatically, copy the URL above)")
	fmt.Println()

	// Wait for code or error with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	var code string
	select {
	case code = <-codeChan:
		// Success - shutdown server
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		_ = server.Shutdown(shutdownCtx)
		return code, nil

	case err := <-errChan:
		// Error from handler
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		_ = server.Shutdown(shutdownCtx)
		return "", err

	case <-ctx.Done():
		// Timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		_ = server.Shutdown(shutdownCtx)
		return "", fmt.Errorf("authentication timeout (waited 5 minutes)")
	}
}

// generateState generates a random state parameter for CSRF protection.
func generateState() string {
	// Use timestamp-based state (simple but effective for CLI use)
	return fmt.Sprintf("gcs-cli-%d", time.Now().UnixNano())
}
