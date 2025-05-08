package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func ExchangeCodeForToken(clientID, projectID, code, codeVerifier string) (*TokenResponse, error) {
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
