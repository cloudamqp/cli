package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInstanceAPI(t *testing.T) {
	apiKey := "test-instance-api-key"
	client := NewInstanceAPI(apiKey)

	assert.NotNil(t, client)
	assert.Equal(t, apiKey, client.apiKey)
	assert.NotNil(t, client.httpClient)
}

func TestInstanceAPI_ListNodes(t *testing.T) {
	// Mock server
	expectedNodes := []Node{
		{
			Name:            "node-01",
			Hostname:        "node-01.example.com",
			Running:         true,
			RabbitMQVersion: "3.10.5",
			ErlangVersion:   "24.1.7",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/nodes", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedNodes)
	}))
	defer server.Close()

	// Test
	client := NewInstanceAPI("test-api-key")
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	nodes, err := client.ListNodes()

	assert.NoError(t, err)
	assert.Len(t, nodes, 1)
	assert.Equal(t, expectedNodes[0].Name, nodes[0].Name)
	assert.Equal(t, expectedNodes[0].Running, nodes[0].Running)
}

func TestInstanceAPI_ListPlugins(t *testing.T) {
	// Mock server
	expectedPlugins := []Plugin{
		{
			Name:        "rabbitmq_management",
			Version:     "3.10.5",
			Description: "RabbitMQ Management Console",
			Enabled:     true,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/plugins", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedPlugins)
	}))
	defer server.Close()

	// Test
	client := NewInstanceAPI("test-api-key")
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	plugins, err := client.ListPlugins()

	assert.NoError(t, err)
	assert.Len(t, plugins, 1)
	assert.Equal(t, expectedPlugins[0].Name, plugins[0].Name)
	assert.Equal(t, expectedPlugins[0].Enabled, plugins[0].Enabled)
}

func TestInstanceAPI_RotatePassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/account/rotate-password", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewInstanceAPI("test-api-key")
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	err := client.RotatePassword()
	assert.NoError(t, err)
}

func TestInstanceAPI_RotateInstanceAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/account/rotate-apikey", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewInstanceAPI("test-api-key")
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	err := client.RotateInstanceAPIKey()
	assert.NoError(t, err)
}

func TestInstanceAPI_RestartRabbitMQ(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/actions/restart", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req ActionRequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, []string{"node1", "node2"}, req.Nodes)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewInstanceAPI("test-api-key")
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	err := client.RestartRabbitMQ([]string{"node1", "node2"})
	assert.NoError(t, err)
}

func TestInstanceAPI_RestartCluster(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/actions/cluster-restart", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewInstanceAPI("test-api-key")
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	err := client.RestartCluster()
	assert.NoError(t, err)
}

func TestInstanceAPI_ToggleHiPE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/actions/hipe", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req HiPERequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.True(t, req.Enable)
		assert.Equal(t, []string{"node1"}, req.Nodes)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewInstanceAPI("test-api-key")
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	req := &HiPERequest{
		Enable: true,
		Nodes:  []string{"node1"},
	}

	err := client.ToggleHiPE(req)
	assert.NoError(t, err)
}

func TestInstanceAPI_ToggleFirehose(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/actions/firehose", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req FirehoseRequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.True(t, req.Enable)
		assert.Equal(t, "/", req.VHost)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewInstanceAPI("test-api-key")
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	req := &FirehoseRequest{
		Enable: true,
		VHost:  "/",
	}

	err := client.ToggleFirehose(req)
	assert.NoError(t, err)
}

func TestInstanceAPI_UpgradeRabbitMQ(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/actions/upgrade-rabbitmq", r.URL.Path)

		var req UpgradeRequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "3.10.7", req.Version)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewInstanceAPI("test-api-key")
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	err := client.UpgradeRabbitMQ("3.10.7")
	assert.NoError(t, err)
}

func TestInstanceAPI_GetAvailableVersions(t *testing.T) {
	expectedVersions := VersionInfo{
		RabbitMQVersions: []string{"3.10.7", "3.10.5"},
		ErlangVersions:   []string{"24.1.7", "24.0.4"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/nodes/available-versions", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedVersions)
	}))
	defer server.Close()

	client := NewInstanceAPI("test-api-key")
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	versions, err := client.GetAvailableVersions()

	assert.NoError(t, err)
	assert.Equal(t, expectedVersions.RabbitMQVersions, versions.RabbitMQVersions)
	assert.Equal(t, expectedVersions.ErlangVersions, versions.ErlangVersions)
}

func TestInstanceAPI_GetUpgradeVersions(t *testing.T) {
	expectedVersions := map[string]string{
		"rabbitmq_version": "3.10.7",
		"erlang_version":   "24.1.7",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/actions/new-rabbitmq-erlang-versions", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedVersions)
	}))
	defer server.Close()

	client := NewInstanceAPI("test-api-key")
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	versions, err := client.GetUpgradeVersions()

	assert.NoError(t, err)
	assert.Equal(t, expectedVersions, versions)
}

func TestInstanceAPI_MultipleActions(t *testing.T) {
	actionsCalled := make(map[string]bool)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actionsCalled[r.URL.Path] = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewInstanceAPI("test-api-key")
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	// Test multiple actions
	err := client.StopCluster()
	assert.NoError(t, err)

	err = client.StartCluster()
	assert.NoError(t, err)

	err = client.UpgradeErlang()
	assert.NoError(t, err)

	err = client.UpgradeRabbitMQErlang()
	assert.NoError(t, err)

	// Verify all actions were called
	assert.True(t, actionsCalled["/actions/cluster-stop"])
	assert.True(t, actionsCalled["/actions/cluster-start"])
	assert.True(t, actionsCalled["/actions/upgrade-erlang"])
	assert.True(t, actionsCalled["/actions/upgrade-rabbitmq-erlang"])
}

func TestInstanceAPI_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Upgrade already in progress"}`))
	}))
	defer server.Close()

	client := NewInstanceAPI("test-api-key")
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	err := client.UpgradeErlang()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error (400): Upgrade already in progress")
}
