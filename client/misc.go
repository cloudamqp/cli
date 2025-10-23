package client

import (
	"encoding/json"
)

type APIKeyRotateResponse struct {
	APIKey string `json:"apikey"`
}

func (c *Client) GetAuditLogCSV(timestamp string) (string, error) {
	endpoint := "/auditlog/csv"
	if timestamp != "" {
		endpoint += "?timestamp=" + timestamp
	}

	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func (c *Client) RotateAPIKey() (*APIKeyRotateResponse, error) {
	respBody, err := c.makeRequest("POST", "/apikeys/rotate-apikey", nil)
	if err != nil {
		return nil, err
	}

	var response APIKeyRotateResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
