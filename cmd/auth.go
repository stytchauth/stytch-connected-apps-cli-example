/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var (
	clientID     string
	projectID    string
	redirectURI  = "http://localhost:8080/callback"
	authorizeURL = "http://localhost:3000/authorize"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
	RequestID    string `json:"request_id"`
	StatusCode   int    `json:"status_code"`
}

// generateCodeVerifier generates a random code verifier for PKCE
func generateCodeVerifier() (string, error) {
	// Generate 32 random bytes
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Base64URL encode the bytes
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// generateCodeChallenge generates a code challenge from a verifier
func generateCodeChallenge(verifier string) string {
	// Hash the verifier with SHA-256
	h := sha256.New()
	h.Write([]byte(verifier))
	// Base64URL encode the hash
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with your Stytch-connected account",
	Run: func(cmd *cobra.Command, args []string) {
		// Generate PKCE values
		codeVerifier, err := generateCodeVerifier()
		if err != nil {
			fmt.Printf("Error generating code verifier: %v\n", err)
			return
		}
		codeChallenge := generateCodeChallenge(codeVerifier)

		// Start local server to receive the callback
		server := &http.Server{
			Addr: ":8080",
		}

		// Channel to receive the auth code
		codeChan := make(chan string)
		errorChan := make(chan error)

		// Set up the callback handler
		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code == "" {
				errorChan <- fmt.Errorf("no code received in callback")
				return
			}

			// Send success response to browser
			w.Write([]byte("Authentication successful! You can close this window."))
			codeChan <- code
		})

		// Start the server in a goroutine
		go func() {
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				errorChan <- err
			}
		}()

		// Construct the auth URL with PKCE parameters
		authURL := fmt.Sprintf(authorizeURL+"?client_id=%s&redirect_uri=%s&response_type=code&code_challenge=%s&code_challenge_method=S256",
			clientID, redirectURI, codeChallenge)

		fmt.Println("Opening browser for authentication...")

		// Open the browser with the auth URL
		err = openBrowser(authURL)
		if err != nil {
			fmt.Println("Please open the following URL in your browser:", authURL)
		}

		// Wait for either the code or an error
		var code string
		select {
		case code = <-codeChan:
			fmt.Println("Received authorization code")
		case err := <-errorChan:
			fmt.Printf("Error: %v\n", err)
			return
		case <-time.After(5 * time.Minute):
			fmt.Println("Timeout waiting for authentication")
			return
		}

		// Exchange the code for a token
		token, err := exchangeCodeForToken(code, codeVerifier)
		if err != nil {
			fmt.Printf("Error exchanging code for token: %v\n", err)
			return
		}

		fmt.Printf("Successfully obtained access token: %s\n", token.AccessToken)
		fmt.Printf("Token expires in: %d seconds\n", token.ExpiresIn)
		fmt.Printf("Refresh token: %s\n", token.RefreshToken)

		// Shutdown the server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	},
}

func exchangeCodeForToken(code, codeVerifier string) (*TokenResponse, error) {
	url := fmt.Sprintf("https://test.stytch.com/v1/public/%s/oauth2/token", projectID)

	// Prepare the request body
	body := map[string]string{
		"client_id":     clientID,
		"code_verifier": codeVerifier,
		"code":          code,
		"grant_type":    "authorization_code",
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Parse the response
	var tokenResp TokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	if tokenResp.StatusCode != 200 {
		return nil, fmt.Errorf("error from API: %s", string(respBody))
	}

	return &tokenResp, nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default: // linux, freebsd, etc.
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}

func init() {
	rootCmd.AddCommand(authCmd)

	// Add flags for client credentials
	authCmd.Flags().StringVar(&clientID, "client-id", "", "OAuth client ID")
	authCmd.Flags().StringVar(&projectID, "project-id", "", "Stytch project ID")
	authCmd.MarkFlagRequired("client-id")
	authCmd.MarkFlagRequired("project-id")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
