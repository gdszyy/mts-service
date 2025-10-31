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
	BookmakerID string   `xml:"bookmaker_id,attr"`
	VirtualHost string `xml:"virtual_host,attr"`
}

// FetchBookmakerInfo calls the whoami.xml endpoint to get the Bookmaker ID and VirtualHost
func FetchBookmakerInfo(accessToken string, uofAPIBaseURL string) (string, string, error) {
	// Construct the API URL using the provided base URL
	apiURL := fmt.Sprintf("%s/v1/users/whoami.xml", uofAPIBaseURL)

	// Create HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set x-access-token header (UOF API uses this instead of Authorization)
	req.Header.Set("x-access-token", accessToken)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse XML
	var whoami WhoAmIResponse
	if err := xml.Unmarshal(body, &whoami); err != nil {
		return "", "", fmt.Errorf("failed to parse XML: %w", err)
	}

	// Validate Bookmaker ID and VirtualHost
	if whoami.BookmakerID == "" {
		return "", "", fmt.Errorf("bookmaker_id is empty in response")
	}
	if whoami.VirtualHost == "" {
		return "", "", fmt.Errorf("virtual_host is empty in response")
	}

	return whoami.BookmakerID, whoami.VirtualHost, nil
}

