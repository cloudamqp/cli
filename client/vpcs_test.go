package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListVPCs(t *testing.T) {
	// Mock server
	expectedVPCs := []VPC{
		{
			ID:     5678,
			Name:   "test-vpc",
			Plan:   "vpc",
			Region: "amazon-web-services::us-east-1",
			Subnet: "10.0.0.0/24",
			Tags:   []string{"production"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/vpcs", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedVPCs)
	}))
	defer server.Close()

	// Test
	client := NewWithBaseURL("test-api-key", server.URL, "test")

	vpcs, err := client.ListVPCs()

	assert.NoError(t, err)
	assert.Len(t, vpcs, 1)
	assert.Equal(t, expectedVPCs[0].ID, vpcs[0].ID)
	assert.Equal(t, expectedVPCs[0].Name, vpcs[0].Name)
	assert.Equal(t, expectedVPCs[0].Subnet, vpcs[0].Subnet)
}

func TestGetVPC(t *testing.T) {
	// Mock server
	expectedVPC := VPC{
		ID:        5678,
		Name:      "test-vpc",
		Plan:      "vpc",
		Region:    "amazon-web-services::us-east-1",
		Subnet:    "10.0.0.0/24",
		Tags:      []string{"production"},
		Instances: []int{1234, 5678},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/vpcs/5678", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedVPC)
	}))
	defer server.Close()

	// Test
	client := NewWithBaseURL("test-api-key", server.URL, "test")

	vpc, err := client.GetVPC(5678)

	assert.NoError(t, err)
	assert.Equal(t, expectedVPC.ID, vpc.ID)
	assert.Equal(t, expectedVPC.Subnet, vpc.Subnet)
	assert.Equal(t, expectedVPC.Instances, vpc.Instances)
}

func TestCreateVPC(t *testing.T) {
	// Mock server
	expectedResponse := VPCCreateResponse{
		ID:     5678,
		APIKey: "new-vpc-api-key",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/vpcs", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		// Parse form data
		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "test-vpc", r.FormValue("name"))
		assert.Equal(t, "amazon-web-services::us-east-1", r.FormValue("region"))
		assert.Equal(t, "10.0.0.0/24", r.FormValue("subnet"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	// Test
	client := NewWithBaseURL("test-api-key", server.URL, "test")

	req := &VPCCreateRequest{
		Name:   "test-vpc",
		Region: "amazon-web-services::us-east-1",
		Subnet: "10.0.0.0/24",
	}

	response, err := client.CreateVPC(req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.APIKey, response.APIKey)
}

func TestCreateVPC_WithTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		assert.NoError(t, err)

		// Check that tags are properly encoded as array
		tags := r.Form["tags[]"]
		assert.Len(t, tags, 2)
		assert.Contains(t, tags, "production")
		assert.Contains(t, tags, "network")

		response := VPCCreateResponse{ID: 5678}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	req := &VPCCreateRequest{
		Name:   "test-vpc",
		Region: "amazon-web-services::us-east-1",
		Subnet: "10.0.0.0/24",
		Tags:   []string{"production", "network"},
	}

	_, err := client.CreateVPC(req)
	assert.NoError(t, err)
}

func TestUpdateVPC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/vpcs/5678", r.URL.Path)

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "updated-vpc-name", r.FormValue("name"))

		// Check tags
		tags := r.Form["tags[]"]
		assert.Len(t, tags, 1)
		assert.Contains(t, tags, "updated")

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	req := &VPCUpdateRequest{
		Name: "updated-vpc-name",
		Tags: []string{"updated"},
	}

	err := client.UpdateVPC(5678, req)
	assert.NoError(t, err)
}

func TestDeleteVPC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/vpcs/5678", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	err := client.DeleteVPC(5678)
	assert.NoError(t, err)
}

func TestVPCError_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "VPC not found"}`))
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	_, err := client.GetVPC(9999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error (404): VPC not found")
}
