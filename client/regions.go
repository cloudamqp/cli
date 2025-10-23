package client

import (
	"encoding/json"
)

type Region struct {
	Provider       string `json:"provider"`
	Region         string `json:"region"`
	Name           string `json:"name"`
	HasSharedPlans bool   `json:"has_shared_plans"`
}

type Plan struct {
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	Backend string  `json:"backend"`
	Shared  bool    `json:"shared"`
}

func (c *Client) ListRegions(provider string) ([]Region, error) {
	endpoint := "/regions"
	if provider != "" {
		endpoint += "?provider=" + provider
	}

	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var regions []Region
	if err := json.Unmarshal(respBody, &regions); err != nil {
		return nil, err
	}

	return regions, nil
}

func (c *Client) ListPlans(backend string) ([]Plan, error) {
	endpoint := "/plans"
	if backend != "" {
		endpoint += "?backend=" + backend
	}

	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var plans []Plan
	if err := json.Unmarshal(respBody, &plans); err != nil {
		return nil, err
	}

	return plans, nil
}
