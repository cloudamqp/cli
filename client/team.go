package client

import (
	"encoding/json"
	"net/url"
)

type TeamMember struct {
	ID             string   `json:"id"`
	Email          string   `json:"email"`
	TFAAuthEnabled bool     `json:"tfa_auth_enabled"`
	Roles          []string `json:"roles"`
}

type TeamInviteRequest struct {
	Email string   `json:"email"`
	Role  string   `json:"role,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

type TeamUpdateRequest struct {
	Role string   `json:"role,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

type TeamRemoveRequest struct {
	Email string `json:"email"`
}

type TeamResponse struct {
	Message string `json:"message"`
}

func (c *Client) ListTeamMembers() ([]TeamMember, error) {
	respBody, err := c.makeRequest("GET", "/team", nil)
	if err != nil {
		return nil, err
	}

	var members []TeamMember
	if err := json.Unmarshal(respBody, &members); err != nil {
		return nil, err
	}

	return members, nil
}

func (c *Client) InviteTeamMember(req *TeamInviteRequest) (*TeamResponse, error) {
	formData := url.Values{}
	formData.Set("email", req.Email)
	if req.Role != "" {
		formData.Set("role", req.Role)
	}
	if len(req.Tags) > 0 {
		for _, tag := range req.Tags {
			formData.Add("tags[]", tag)
		}
	}

	respBody, err := c.makeRequest("POST", "/team/invite", formData)
	if err != nil {
		return nil, err
	}

	var response TeamResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) RemoveTeamMember(email string) (*TeamResponse, error) {
	formData := url.Values{}
	formData.Set("email", email)

	respBody, err := c.makeRequest("POST", "/team/remove", formData)
	if err != nil {
		return nil, err
	}

	var response TeamResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) UpdateTeamMember(userID string, req *TeamUpdateRequest) (*TeamResponse, error) {
	endpoint := "/team/" + userID

	formData := url.Values{}
	if req.Role != "" {
		formData.Set("role", req.Role)
	}
	if len(req.Tags) > 0 {
		for _, tag := range req.Tags {
			formData.Add("tags[]", tag)
		}
	}

	respBody, err := c.makeRequest("PUT", endpoint, formData)
	if err != nil {
		return nil, err
	}

	var response TeamResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
