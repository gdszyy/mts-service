package client

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WhoAmIResponse represents the XML response from whoami.xml
type WhoAmIResponse struct {
	XMLName    xml.Name `xml:"bookmaker_details"`
	BookmakerID string   `xml:"bookmaker_id"`
	ExpireAt   string   `xml:"expire_at"`
	VirtualHost string   `xml:"virtual_host"`
}

// FetchBookmakerID calls the whoami.xml endpoint to get the Bookmaker ID
func FetchBookmakerID(accessToken string, production bool) (string, error) {
	// Determine the API URL based on environment
	apiURL := "https://global.api.betradar.com/v1/users/whoami.xml"
	if production {
		apiURL = "https://api.betradar.com/v1/users/whoami.xml"
	}

	// Create HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse XML
	var whoami WhoAmIResponse
	if err := xml.Unmarshal(body, &whoami); err != nil {
		return "", fmt.Errorf("failed to parse XML: %w", err)
	}

	// Validate Bookmaker ID
	if whoami.BookmakerID == "" {
		return "", fmt.Errorf("bookmaker_id is empty in response")
	}

	return whoami.BookmakerID, nil
}

