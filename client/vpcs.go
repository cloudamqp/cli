package client

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type VPC struct {
	ID         int      `json:"id"`
	Plan       string   `json:"plan"`
	Region     string   `json:"region"`
	Name       string   `json:"name"`
	Tags       []string `json:"tags"`
	ProviderID string   `json:"providerid"`
	Subnet     string   `json:"subnet"`
	Instances  []int    `json:"instances"`
}

type VPCCreateRequest struct {
	Name   string   `json:"name"`
	Region string   `json:"region"`
	Subnet string   `json:"subnet"`
	Tags   []string `json:"tags,omitempty"`
}

type VPCCreateResponse struct {
	ID     int    `json:"id"`
	APIKey string `json:"apikey"`
}

type VPCUpdateRequest struct {
	Name string   `json:"name,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

func (c *Client) ListVPCs() ([]VPC, error) {
	respBody, err := c.makeRequest("GET", "/vpcs", nil)
	if err != nil {
		return nil, err
	}

	var vpcs []VPC
	if err := json.Unmarshal(respBody, &vpcs); err != nil {
		return nil, err
	}

	return vpcs, nil
}

func (c *Client) GetVPC(id int) (*VPC, error) {
	endpoint := "/vpcs/" + strconv.Itoa(id)
	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var vpc VPC
	if err := json.Unmarshal(respBody, &vpc); err != nil {
		return nil, err
	}

	return &vpc, nil
}

func (c *Client) CreateVPC(req *VPCCreateRequest) (*VPCCreateResponse, error) {
	formData := url.Values{}
	formData.Set("name", req.Name)
	formData.Set("region", req.Region)
	formData.Set("subnet", req.Subnet)

	if len(req.Tags) > 0 {
		for _, tag := range req.Tags {
			formData.Add("tags[]", tag)
		}
	}

	respBody, err := c.makeRequest("POST", "/vpcs", formData)
	if err != nil {
		return nil, err
	}

	var createResp VPCCreateResponse
	if err := json.Unmarshal(respBody, &createResp); err != nil {
		return nil, err
	}

	return &createResp, nil
}

func (c *Client) UpdateVPC(id int, req *VPCUpdateRequest) error {
	endpoint := "/vpcs/" + strconv.Itoa(id)

	formData := url.Values{}
	if req.Name != "" {
		formData.Set("name", req.Name)
	}
	if len(req.Tags) > 0 {
		for _, tag := range req.Tags {
			formData.Add("tags[]", tag)
		}
	}

	_, err := c.makeRequest("PUT", endpoint, formData)
	return err
}

func (c *Client) DeleteVPC(id int) error {
	endpoint := "/vpcs/" + strconv.Itoa(id)
	_, err := c.makeRequest("DELETE", endpoint, nil)
	return err
}
