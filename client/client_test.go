package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	apiKey := "test-api-key"
	client := New(apiKey, "test")

	assert.NotNil(t, client)
	assert.Equal(t, apiKey, client.apiKey)
	assert.NotNil(t, client.httpClient)
}

func TestNew_WithEnvironmentVariable(t *testing.T) {
	// Save original environment variable
	originalURL := os.Getenv("CLOUDAMQP_URL")
	defer os.Setenv("CLOUDAMQP_URL", originalURL)

	// Test with custom base URL from environment variable
	customURL := "https://custom.example.com/api"
	os.Setenv("CLOUDAMQP_URL", customURL)

	apiKey := "test-api-key"
	client := New(apiKey, "test")

	assert.NotNil(t, client)
	assert.Equal(t, apiKey, client.apiKey)
	assert.Equal(t, customURL, client.baseURL)
	assert.NotNil(t, client.httpClient)

	// Test with empty environment variable (should use default)
	os.Setenv("CLOUDAMQP_URL", "")
	client = New(apiKey, "test")

	assert.NotNil(t, client)
	assert.Equal(t, "https://customer.cloudamqp.com/api", client.baseURL)
}

func TestMakeRequest_GET_Success(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/test", r.URL.Path)

		// Verify auth
		username, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "", username)
		assert.Equal(t, "test-api-key", password)

		// Return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewWithBaseURL("test-api-key", server.URL, "test")

	// Test request
	resp, err := client.makeRequest("GET", "/test", nil)

	assert.NoError(t, err)
	assert.Equal(t, `{"success": true}`, string(resp))
}

func TestMakeRequest_POST_FormData(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		// Parse form data
		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "test-value", r.FormValue("test-key"))

		// Return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 123}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewWithBaseURL("test-api-key", server.URL, "test")

	// Test request with form data
	formData := url.Values{}
	formData.Set("test-key", "test-value")

	resp, err := client.makeRequest("POST", "/test", formData)

	assert.NoError(t, err)
	assert.Equal(t, `{"id": 123}`, string(resp))
}

func TestMakeRequest_POST_JSON(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Read body
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		assert.Equal(t, `{"test":"value"}`, string(body))

		// Return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"created": true}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewWithBaseURL("test-api-key", server.URL, "test")

	// Test request with JSON data
	jsonData := map[string]string{"test": "value"}

	resp, err := client.makeRequest("POST", "/test", jsonData)

	assert.NoError(t, err)
	assert.Equal(t, `{"created": true}`, string(resp))
}

func TestMakeRequest_APIError_JSON(t *testing.T) {
	// Mock server returning error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request"}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewWithBaseURL("test-api-key", server.URL, "test")

	// Test request
	_, err := client.makeRequest("GET", "/test", nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error (400): Invalid request")
}

func TestMakeRequest_APIError_Plain(t *testing.T) {
	// Mock server returning plain text error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not authorized"))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewWithBaseURL("test-api-key", server.URL, "test")

	// Test request
	_, err := client.makeRequest("GET", "/test", nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error (401): Not authorized")
}

func TestMakeRequest_NetworkError(t *testing.T) {
	// Create client with invalid URL
	client := NewWithBaseURL("test-api-key", "http://invalid-url-that-does-not-exist", "test")

	// Test request
	_, err := client.makeRequest("GET", "/test", nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request failed")
}

func TestMakeRequest_InvalidJSON(t *testing.T) {
	client := New("test-api-key", "test")

	// Test with invalid JSON data
	invalidData := make(chan int) // channels can't be marshaled to JSON

	_, err := client.makeRequest("POST", "/test", invalidData)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal request body")
}
