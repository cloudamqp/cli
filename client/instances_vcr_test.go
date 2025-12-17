package client

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v2/cassette"
	"gopkg.in/dnaeon/go-vcr.v2/recorder"
)

// sanitizeResponseBody removes sensitive data from API responses
func sanitizeResponseBody(body string) string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return body
	}

	// Sanitize sensitive fields
	if _, ok := data["apikey"]; ok {
		data["apikey"] = "REDACTED"
	}

	// Replace credentials in AMQP/AMQPS URLs
	credRegex := regexp.MustCompile(`://([^:]+):([^@]+)@`)

	// Sanitize single URL field
	if urlStr, ok := data["url"].(string); ok {
		data["url"] = credRegex.ReplaceAllString(urlStr, "://REDACTED:REDACTED@")
	}

	// Sanitize urls object (contains external/internal URLs)
	if urls, ok := data["urls"].(map[string]interface{}); ok {
		for key, val := range urls {
			if urlStr, ok := val.(string); ok {
				urls[key] = credRegex.ReplaceAllString(urlStr, "://REDACTED:REDACTED@")
			}
		}
	}

	sanitized, err := json.Marshal(data)
	if err != nil {
		return body
	}

	return string(sanitized)
}

// TestCreateInstanceVCR tests the CreateInstance method using VCR to record/replay HTTP interactions
func TestCreateInstanceVCR(t *testing.T) {
	// Create a VCR recorder
	r, err := recorder.New("fixtures/create_instance")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	// Add filter to sanitize sensitive data in cassettes
	r.AddFilter(func(i *cassette.Interaction) error {
		// Sanitize Authorization header
		delete(i.Request.Headers, "Authorization")

		// Sanitize sensitive data in response body (API keys, passwords)
		i.Response.Body = sanitizeResponseBody(i.Response.Body)

		// Remove session cookies
		delete(i.Response.Headers, "Set-Cookie")

		return nil
	})

	// Get API key - only required if cassette doesn't exist (recording mode)
	apiKey := os.Getenv("CLOUDAMQP_APIKEY")
	// If no API key, use dummy (cassette will be used if it exists)
	if apiKey == "" {
		apiKey = "vcr-replay-mode"
	}

	// Create HTTP client with VCR recorder as transport
	httpClient := &http.Client{Transport: r}

	// Create client with VCR HTTP client
	client := NewWithHTTPClient(apiKey, "https://customer.cloudamqp.com/api", "test", httpClient)

	// Create instance request
	req := &InstanceCreateRequest{
		Name:   "vcr-test-instance",
		Plan:   "lemur",
		Region: "amazon-web-services::us-east-1",
		Tags:   []string{"test", "vcr"},
	}

	// Execute the create instance request
	resp, err := client.CreateInstance(req)

	// Verify the response
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotZero(t, resp.ID)
	assert.NotEmpty(t, resp.URL)
	assert.NotEmpty(t, resp.APIKey)

	t.Logf("Created instance with ID: %d", resp.ID)
	t.Logf("Instance URL: %s", resp.URL)
}
