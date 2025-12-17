package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListInstances(t *testing.T) {
	// Mock server
	expectedInstances := []Instance{
		{
			ID:     1234,
			Name:   "test-instance",
			Plan:   "bunny-1",
			Region: "amazon-web-services::us-east-1",
			Ready:  true,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/instances", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedInstances)
	}))
	defer server.Close()

	// Test
	client := NewWithBaseURL("test-api-key", server.URL, "test")

	instances, err := client.ListInstances()

	assert.NoError(t, err)
	assert.Len(t, instances, 1)
	assert.Equal(t, expectedInstances[0].ID, instances[0].ID)
	assert.Equal(t, expectedInstances[0].Name, instances[0].Name)
}

func TestGetInstance(t *testing.T) {
	// Mock server
	expectedInstance := Instance{
		ID:     1234,
		Name:   "test-instance",
		Plan:   "bunny-1",
		Region: "amazon-web-services::us-east-1",
		APIKey: "instance-api-key",
		Ready:  true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/instances/1234", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedInstance)
	}))
	defer server.Close()

	// Test
	client := NewWithBaseURL("test-api-key", server.URL, "test")

	instance, err := client.GetInstance(1234)

	assert.NoError(t, err)
	assert.Equal(t, expectedInstance.ID, instance.ID)
	assert.Equal(t, expectedInstance.APIKey, instance.APIKey)
}

func TestCreateInstance(t *testing.T) {
	// Mock server
	expectedResponse := InstanceCreateResponse{
		ID:     1234,
		URL:    "amqp://user:pass@host/vhost",
		APIKey: "new-instance-api-key",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/instances", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		// Parse form data
		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "test-instance", r.FormValue("name"))
		assert.Equal(t, "bunny-1", r.FormValue("plan"))
		assert.Equal(t, "amazon-web-services::us-east-1", r.FormValue("region"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	// Test
	client := NewWithBaseURL("test-api-key", server.URL, "test")

	req := &InstanceCreateRequest{
		Name:   "test-instance",
		Plan:   "bunny-1",
		Region: "amazon-web-services::us-east-1",
	}

	response, err := client.CreateInstance(req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.APIKey, response.APIKey)
}

func TestCreateInstance_WithTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		assert.NoError(t, err)

		// Check that tags are properly encoded as array
		tags := r.Form["tags[]"]
		assert.Len(t, tags, 2)
		assert.Contains(t, tags, "production")
		assert.Contains(t, tags, "web-app")

		response := InstanceCreateResponse{ID: 1234}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	req := &InstanceCreateRequest{
		Name:   "test-instance",
		Plan:   "bunny-1",
		Region: "amazon-web-services::us-east-1",
		Tags:   []string{"production", "web-app"},
	}

	_, err := client.CreateInstance(req)
	assert.NoError(t, err)
}

func TestCreateInstance_WithVPC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		assert.NoError(t, err)

		assert.Equal(t, "10.0.0.0/24", r.FormValue("vpc_subnet"))
		assert.Equal(t, "5678", r.FormValue("vpc_id"))

		response := InstanceCreateResponse{ID: 1234}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	vpcID := 5678
	req := &InstanceCreateRequest{
		Name:      "test-instance",
		Plan:      "bunny-1",
		Region:    "amazon-web-services::us-east-1",
		VPCSubnet: "10.0.0.0/24",
		VPCID:     &vpcID,
	}

	_, err := client.CreateInstance(req)
	assert.NoError(t, err)
}

func TestUpdateInstance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/instances/1234", r.URL.Path)

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "updated-name", r.FormValue("name"))
		assert.Equal(t, "rabbit-1", r.FormValue("plan"))

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	req := &InstanceUpdateRequest{
		Name: "updated-name",
		Plan: "rabbit-1",
	}

	err := client.UpdateInstance(1234, req)
	assert.NoError(t, err)
}

func TestDeleteInstance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/instances/1234", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	err := client.DeleteInstance(1234)
	assert.NoError(t, err)
}

func TestResizeInstanceDisk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/instances/1234/disk", r.URL.Path)

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "100", r.FormValue("extra_disk_size"))
		assert.Equal(t, "true", r.FormValue("allow_downtime"))

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	req := &DiskResizeRequest{
		ExtraDiskSize: 100,
		AllowDowntime: true,
	}

	err := client.ResizeInstanceDisk(1234, req)
	assert.NoError(t, err)
}

func TestResizeInstanceDisk_NoDowntime(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "50", r.FormValue("extra_disk_size"))
		// allow_downtime should not be set when false
		assert.Empty(t, r.FormValue("allow_downtime"))

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	req := &DiskResizeRequest{
		ExtraDiskSize: 50,
		AllowDowntime: false,
	}

	err := client.ResizeInstanceDisk(1234, req)
	assert.NoError(t, err)
}
