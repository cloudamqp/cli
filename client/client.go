package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var BaseURL = "https://customer.cloudamqp.com/api"

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	version    string
}

func New(apiKey, version string) *Client {
	baseURL := "https://customer.cloudamqp.com/api"
	if envURL := os.Getenv("CLOUDAMQP_URL"); envURL != "" {
		baseURL = envURL
	}
	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: &http.Client{},
		version:    version,
	}
}

func NewWithBaseURL(apiKey, baseURL, version string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: &http.Client{},
		version:    version,
	}
}

// NewWithHTTPClient creates a new client with a custom HTTP client.
// This is useful for testing with tools like go-vcr.
func NewWithHTTPClient(apiKey, baseURL, version string, httpClient *http.Client) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: httpClient,
		version:    version,
	}
}

func (c *Client) makeRequest(method, endpoint string, body any) ([]byte, error) {
	var reqBody io.Reader
	var contentType string

	if body != nil {
		switch v := body.(type) {
		case url.Values:
			contentType = "application/x-www-form-urlencoded"
			reqBody = strings.NewReader(v.Encode())
		default:
			contentType = "application/json"
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
			reqBody = bytes.NewReader(jsonData)
		}
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth("", c.apiKey)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("User-Agent", fmt.Sprintf("cloudamqp-cli/%s", c.version))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errorResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(respBody, &errorResp); err == nil && errorResp.Error != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, errorResp.Error)
		}
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Instance-specific operations using /instances/{id}/ endpoints

// Node management
type Node struct {
	Name               string `json:"name"`
	Hostname           string `json:"hostname"`
	Configured         bool   `json:"configured"`
	HiPE               bool   `json:"hipe"`
	RabbitMQVersion    string `json:"rabbitmq_version"`
	ErlangVersion      string `json:"erlang_version"`
	Running            bool   `json:"running"`
	DiskSize           int    `json:"disk_size"`
	AdditionalDiskSize int    `json:"additional_disk_size"`
	AvailabilityZone   string `json:"availability_zone"`
	HostnameInternal   string `json:"hostname_internal"`
}

func (c *Client) ListNodes(instanceID string) ([]Node, error) {
	endpoint := "/instances/" + instanceID + "/nodes"
	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var nodes []Node
	if err := json.Unmarshal(respBody, &nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

// Plugin management
type Plugin struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func (c *Client) ListPlugins(instanceID string) ([]Plugin, error) {
	endpoint := "/instances/" + instanceID + "/plugins"
	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var plugins []Plugin
	if err := json.Unmarshal(respBody, &plugins); err != nil {
		return nil, err
	}

	return plugins, nil
}

func (c *Client) EnablePlugin(instanceID, pluginName string) error {
	endpoint := "/instances/" + instanceID + "/plugins"

	requestBody := map[string]string{
		"plugin_name": pluginName,
	}

	_, err := c.makeRequest("POST", endpoint, requestBody)
	return err
}

func (c *Client) DisablePlugin(instanceID, pluginName string) error {
	endpoint := "/instances/" + instanceID + "/plugins/" + pluginName
	_, err := c.makeRequest("DELETE", endpoint, nil)
	return err
}

// Account operations
func (c *Client) RotatePassword(instanceID string) error {
	endpoint := "/instances/" + instanceID + "/account/rotate-password"
	_, err := c.makeRequest("POST", endpoint, nil)
	return err
}

func (c *Client) RotateInstanceAPIKey(instanceID string) error {
	endpoint := "/instances/" + instanceID + "/account/rotate-apikey"
	_, err := c.makeRequest("POST", endpoint, nil)
	return err
}

// Action operations
type ActionRequest struct {
	Nodes []string `json:"nodes,omitempty"`
}

type HiPERequest struct {
	Enable bool     `json:"enable"`
	Nodes  []string `json:"nodes,omitempty"`
}

type FirehoseRequest struct {
	Enable bool   `json:"enable"`
	VHost  string `json:"vhost"`
}

type UpgradeRequest struct {
	Version string `json:"version"`
}

type VersionInfo struct {
	RabbitMQVersions []string `json:"rabbitmq_versions"`
	ErlangVersions   []string `json:"erlang_versions"`
}

func (c *Client) ToggleHiPE(instanceID string, req *HiPERequest) error {
	endpoint := "/instances/" + instanceID + "/actions/hipe"
	_, err := c.makeRequest("PUT", endpoint, req)
	return err
}

func (c *Client) ToggleFirehose(instanceID string, req *FirehoseRequest) error {
	endpoint := "/instances/" + instanceID + "/actions/firehose"
	_, err := c.makeRequest("PUT", endpoint, req)
	return err
}

func (c *Client) RestartRabbitMQ(instanceID string, nodes []string) error {
	endpoint := "/instances/" + instanceID + "/actions/restart"
	req := ActionRequest{Nodes: nodes}
	_, err := c.makeRequest("POST", endpoint, req)
	return err
}

func (c *Client) RestartCluster(instanceID string) error {
	endpoint := "/instances/" + instanceID + "/actions/cluster-restart"
	_, err := c.makeRequest("POST", endpoint, nil)
	return err
}

func (c *Client) StopCluster(instanceID string) error {
	endpoint := "/instances/" + instanceID + "/actions/cluster-stop"
	_, err := c.makeRequest("POST", endpoint, nil)
	return err
}

func (c *Client) StartCluster(instanceID string) error {
	endpoint := "/instances/" + instanceID + "/actions/cluster-start"
	_, err := c.makeRequest("POST", endpoint, nil)
	return err
}

func (c *Client) RestartManagement(instanceID string, nodes []string) error {
	endpoint := "/instances/" + instanceID + "/actions/mgmt-restart"
	req := ActionRequest{Nodes: nodes}
	_, err := c.makeRequest("POST", endpoint, req)
	return err
}

func (c *Client) StopInstance(instanceID string, nodes []string) error {
	endpoint := "/instances/" + instanceID + "/actions/stop"
	req := ActionRequest{Nodes: nodes}
	_, err := c.makeRequest("POST", endpoint, req)
	return err
}

func (c *Client) StartInstance(instanceID string, nodes []string) error {
	endpoint := "/instances/" + instanceID + "/actions/start"
	req := ActionRequest{Nodes: nodes}
	_, err := c.makeRequest("POST", endpoint, req)
	return err
}

func (c *Client) RebootInstance(instanceID string, nodes []string) error {
	endpoint := "/instances/" + instanceID + "/actions/reboot"
	req := ActionRequest{Nodes: nodes}
	_, err := c.makeRequest("POST", endpoint, req)
	return err
}

func (c *Client) UpgradeErlang(instanceID string) error {
	endpoint := "/instances/" + instanceID + "/actions/upgrade-erlang"
	_, err := c.makeRequest("POST", endpoint, nil)
	return err
}

func (c *Client) UpgradeRabbitMQ(instanceID string, version string) error {
	endpoint := "/instances/" + instanceID + "/actions/upgrade-rabbitmq"
	req := UpgradeRequest{Version: version}
	_, err := c.makeRequest("POST", endpoint, req)
	return err
}

func (c *Client) UpgradeRabbitMQErlang(instanceID string) error {
	endpoint := "/instances/" + instanceID + "/actions/upgrade-rabbitmq-erlang"
	_, err := c.makeRequest("POST", endpoint, nil)
	return err
}

func (c *Client) GetAvailableVersions(instanceID string) (*VersionInfo, error) {
	endpoint := "/instances/" + instanceID + "/nodes/available-versions"
	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var versions VersionInfo
	if err := json.Unmarshal(respBody, &versions); err != nil {
		return nil, err
	}

	return &versions, nil
}

func (c *Client) GetUpgradeVersions(instanceID string) (map[string]string, error) {
	endpoint := "/instances/" + instanceID + "/actions/new-rabbitmq-erlang-versions"
	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var versions map[string]string
	if err := json.Unmarshal(respBody, &versions); err != nil {
		return nil, err
	}

	return versions, nil
}

// RabbitMQ Config operations
func (c *Client) GetRabbitMQConfig(instanceID string) (map[string]interface{}, error) {
	endpoint := "/instances/" + instanceID + "/config"
	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if err := json.Unmarshal(respBody, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Client) UpdateRabbitMQConfig(instanceID string, config map[string]interface{}) error {
	endpoint := "/instances/" + instanceID + "/config"
	_, err := c.makeRequest("PUT", endpoint, config)
	return err
}
