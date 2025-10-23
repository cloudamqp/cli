package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// InstanceAPIClient for instance-specific operations
type InstanceAPIClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewInstanceAPI(apiKey string) *InstanceAPIClient {
	return &InstanceAPIClient{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *InstanceAPIClient) makeRequest(method, endpoint string, body interface{}) ([]byte, error) {
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

	req, err := http.NewRequest(method, BaseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth("", c.apiKey)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

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

func (c *InstanceAPIClient) ListNodes() ([]Node, error) {
	respBody, err := c.makeRequest("GET", "/nodes", nil)
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

func (c *InstanceAPIClient) ListPlugins() ([]Plugin, error) {
	respBody, err := c.makeRequest("GET", "/plugins", nil)
	if err != nil {
		return nil, err
	}

	var plugins []Plugin
	if err := json.Unmarshal(respBody, &plugins); err != nil {
		return nil, err
	}

	return plugins, nil
}

// Account operations
func (c *InstanceAPIClient) RotatePassword() error {
	_, err := c.makeRequest("POST", "/account/rotate-password", nil)
	return err
}

func (c *InstanceAPIClient) RotateInstanceAPIKey() error {
	_, err := c.makeRequest("POST", "/account/rotate-apikey", nil)
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

func (c *InstanceAPIClient) ToggleHiPE(req *HiPERequest) error {
	_, err := c.makeRequest("PUT", "/actions/hipe", req)
	return err
}

func (c *InstanceAPIClient) ToggleFirehose(req *FirehoseRequest) error {
	_, err := c.makeRequest("PUT", "/actions/firehose", req)
	return err
}

func (c *InstanceAPIClient) RestartRabbitMQ(nodes []string) error {
	req := ActionRequest{Nodes: nodes}
	_, err := c.makeRequest("POST", "/actions/restart", req)
	return err
}

func (c *InstanceAPIClient) RestartCluster() error {
	_, err := c.makeRequest("POST", "/actions/cluster-restart", nil)
	return err
}

func (c *InstanceAPIClient) StopCluster() error {
	_, err := c.makeRequest("POST", "/actions/cluster-stop", nil)
	return err
}

func (c *InstanceAPIClient) StartCluster() error {
	_, err := c.makeRequest("POST", "/actions/cluster-start", nil)
	return err
}

func (c *InstanceAPIClient) RestartManagement(nodes []string) error {
	req := ActionRequest{Nodes: nodes}
	_, err := c.makeRequest("POST", "/actions/mgmt-restart", req)
	return err
}

func (c *InstanceAPIClient) StopInstance(nodes []string) error {
	req := ActionRequest{Nodes: nodes}
	_, err := c.makeRequest("POST", "/actions/stop", req)
	return err
}

func (c *InstanceAPIClient) StartInstance(nodes []string) error {
	req := ActionRequest{Nodes: nodes}
	_, err := c.makeRequest("POST", "/actions/start", req)
	return err
}

func (c *InstanceAPIClient) RebootInstance(nodes []string) error {
	req := ActionRequest{Nodes: nodes}
	_, err := c.makeRequest("POST", "/actions/reboot", req)
	return err
}

func (c *InstanceAPIClient) UpgradeErlang() error {
	_, err := c.makeRequest("POST", "/actions/upgrade-erlang", nil)
	return err
}

func (c *InstanceAPIClient) UpgradeRabbitMQ(version string) error {
	req := UpgradeRequest{Version: version}
	_, err := c.makeRequest("POST", "/actions/upgrade-rabbitmq", req)
	return err
}

func (c *InstanceAPIClient) UpgradeRabbitMQErlang() error {
	_, err := c.makeRequest("POST", "/actions/upgrade-rabbitmq-erlang", nil)
	return err
}

func (c *InstanceAPIClient) GetAvailableVersions() (*VersionInfo, error) {
	respBody, err := c.makeRequest("GET", "/nodes/available-versions", nil)
	if err != nil {
		return nil, err
	}

	var versions VersionInfo
	if err := json.Unmarshal(respBody, &versions); err != nil {
		return nil, err
	}

	return &versions, nil
}

func (c *InstanceAPIClient) GetUpgradeVersions() (map[string]string, error) {
	respBody, err := c.makeRequest("GET", "/actions/new-rabbitmq-erlang-versions", nil)
	if err != nil {
		return nil, err
	}

	var versions map[string]string
	if err := json.Unmarshal(respBody, &versions); err != nil {
		return nil, err
	}

	return versions, nil
}
