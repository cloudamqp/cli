package client

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Instance struct {
	ID               int      `json:"id"`
	Plan             string   `json:"plan"`
	Region           string   `json:"region"`
	Name             string   `json:"name"`
	Tags             []string `json:"tags"`
	ProviderID       string   `json:"providerid"`
	VPCID            *int     `json:"vpc_id"`
	URL              string   `json:"url"`
	APIKey           string   `json:"apikey"`
	Ready            bool     `json:"ready"`
	RMQVersion       string   `json:"rmq_version"`
	HostnameExternal string   `json:"hostname_external"`
	HostnameInternal string   `json:"hostname_internal"`
}

type InstanceCreateRequest struct {
	Name      string   `json:"name"`
	Plan      string   `json:"plan"`
	Region    string   `json:"region"`
	Tags      []string `json:"tags,omitempty"`
	VPCSubnet string   `json:"vpc_subnet,omitempty"`
	VPCID     *int     `json:"vpc_id,omitempty"`
}

type InstanceCreateResponse struct {
	ID     int    `json:"id"`
	URL    string `json:"url"`
	APIKey string `json:"apikey"`
}

type InstanceUpdateRequest struct {
	Name string   `json:"name,omitempty"`
	Plan string   `json:"plan,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

func (c *Client) ListInstances() ([]Instance, error) {
	respBody, err := c.makeRequest("GET", "/instances", nil)
	if err != nil {
		return nil, err
	}

	var instances []Instance
	if err := json.Unmarshal(respBody, &instances); err != nil {
		return nil, err
	}

	return instances, nil
}

func (c *Client) GetInstance(id int) (*Instance, error) {
	endpoint := "/instances/" + strconv.Itoa(id)
	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var instance Instance
	if err := json.Unmarshal(respBody, &instance); err != nil {
		return nil, err
	}

	return &instance, nil
}

func (c *Client) CreateInstance(req *InstanceCreateRequest) (*InstanceCreateResponse, error) {
	formData := url.Values{}
	formData.Set("name", req.Name)
	formData.Set("plan", req.Plan)
	formData.Set("region", req.Region)

	if len(req.Tags) > 0 {
		for _, tag := range req.Tags {
			formData.Add("tags[]", tag)
		}
	}

	if req.VPCSubnet != "" {
		formData.Set("vpc_subnet", req.VPCSubnet)
	}

	if req.VPCID != nil {
		formData.Set("vpc_id", strconv.Itoa(*req.VPCID))
	}

	respBody, err := c.makeRequest("POST", "/instances", formData)
	if err != nil {
		return nil, err
	}

	var createResp InstanceCreateResponse
	if err := json.Unmarshal(respBody, &createResp); err != nil {
		return nil, err
	}

	return &createResp, nil
}

func (c *Client) UpdateInstance(id int, req *InstanceUpdateRequest) error {
	endpoint := "/instances/" + strconv.Itoa(id)

	formData := url.Values{}
	if req.Name != "" {
		formData.Set("name", req.Name)
	}
	if req.Plan != "" {
		formData.Set("plan", req.Plan)
	}
	if len(req.Tags) > 0 {
		for _, tag := range req.Tags {
			formData.Add("tags[]", tag)
		}
	}

	_, err := c.makeRequest("PUT", endpoint, formData)
	return err
}

func (c *Client) DeleteInstance(id int) error {
	endpoint := "/instances/" + strconv.Itoa(id)
	_, err := c.makeRequest("DELETE", endpoint, nil)
	return err
}

type DiskResizeRequest struct {
	ExtraDiskSize int  `json:"extra_disk_size"`
	AllowDowntime bool `json:"allow_downtime,omitempty"`
}

func (c *Client) ResizeInstanceDisk(id int, req *DiskResizeRequest) error {
	endpoint := "/instances/" + strconv.Itoa(id) + "/disk"

	formData := url.Values{}
	formData.Set("extra_disk_size", strconv.Itoa(req.ExtraDiskSize))
	if req.AllowDowntime {
		formData.Set("allow_downtime", "true")
	}

	_, err := c.makeRequest("PUT", endpoint, formData)
	return err
}
