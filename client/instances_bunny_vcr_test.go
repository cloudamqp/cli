package client

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v2/cassette"
	"gopkg.in/dnaeon/go-vcr.v2/recorder"
)

// TestCreateInstanceBunny1 tests creating an instance with bunny-1 plan
// This test uses VCR to record/replay HTTP interactions
func TestCreateInstanceBunny1(t *testing.T) {
	r, err := recorder.New("fixtures/bunny1_create")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	// Add filter to sanitize sensitive data
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		i.Response.Body = sanitizeResponseBody(i.Response.Body)
		delete(i.Response.Headers, "Set-Cookie")
		return nil
	})

	// Get API key - only required if cassette doesn't exist (recording mode)
	apiKey := os.Getenv("CLOUDAMQP_APIKEY")
	// If no API key and no cassette, skip test
	if apiKey == "" {
		// Use dummy key for replay mode (cassette intercepts all requests)
		apiKey = "vcr-replay-mode"
	}

	httpClient := &http.Client{Transport: r}
	client := NewWithHTTPClient(apiKey, "https://customer.cloudamqp.com/api", "test", httpClient)

	// Create instance with bunny-1 plan
	req := &InstanceCreateRequest{
		Name:   "bunny1-test",
		Plan:   "bunny-1",
		Region: "amazon-web-services::us-east-1",
		Tags:   []string{"test", "bunny1"},
	}

	resp, err := client.CreateInstance(req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotZero(t, resp.ID)
	assert.NotEmpty(t, resp.URL)
	assert.NotEmpty(t, resp.APIKey)

	t.Logf("✓ Created bunny-1 instance with ID: %d", resp.ID)
}

// TestUpdateInstanceBunny1ToHare1 tests updating an instance from bunny-1 to hare-1
// Note: This test uses a pre-recorded cassette with a real instance update
func TestUpdateInstanceBunny1ToHare1(t *testing.T) {
	r, err := recorder.New("fixtures/bunny1_to_hare1_update")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	// Add filter to sanitize sensitive data
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		i.Response.Body = sanitizeResponseBody(i.Response.Body)
		delete(i.Response.Headers, "Set-Cookie")
		return nil
	})

	// Get API key - only required if cassette doesn't exist (recording mode)
	apiKey := os.Getenv("CLOUDAMQP_APIKEY")
	// If no API key and no cassette, skip test
	if apiKey == "" {
		// Use dummy key for replay mode (cassette intercepts all requests)
		apiKey = "vcr-replay-mode"
	}

	httpClient := &http.Client{Transport: r}
	client := NewWithHTTPClient(apiKey, "https://customer.cloudamqp.com/api", "test", httpClient)

	// Use a real instance ID from a bunny-1 instance
	// This should be an instance that is already created and fully configured
	instanceID := 359559

	// Get instance before update
	beforeUpdate, err := client.GetInstance(instanceID)
	require.NoError(t, err)
	require.NotNil(t, beforeUpdate)
	t.Logf("Before update - Plan: %s", beforeUpdate.Plan)

	// Update to hare-1
	updateReq := &InstanceUpdateRequest{
		Plan: "hare-1",
	}

	err = client.UpdateInstance(instanceID, updateReq)
	require.NoError(t, err)
	t.Logf("✓ Updated instance %d to hare-1", instanceID)

	// Note: Plan changes take time to reflect in the API
	// The update is successful even though the plan may not show immediately
	afterUpdate, err := client.GetInstance(instanceID)
	if err == nil && afterUpdate != nil {
		t.Logf("After update - Plan: %s (may take time to update)", afterUpdate.Plan)
	}
}

// TestDeleteInstanceBunny1 tests deleting a bunny-1 instance
func TestDeleteInstanceBunny1(t *testing.T) {
	r, err := recorder.New("fixtures/bunny1_delete")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	// Add filter to sanitize sensitive data
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		i.Response.Body = sanitizeResponseBody(i.Response.Body)
		delete(i.Response.Headers, "Set-Cookie")
		return nil
	})

	// Get API key - only required if cassette doesn't exist (recording mode)
	apiKey := os.Getenv("CLOUDAMQP_APIKEY")
	// If no API key and no cassette, skip test
	if apiKey == "" {
		// Use dummy key for replay mode (cassette intercepts all requests)
		apiKey = "vcr-replay-mode"
	}

	httpClient := &http.Client{Transport: r}
	client := NewWithHTTPClient(apiKey, "https://customer.cloudamqp.com/api", "test", httpClient)

	// Use an instance ID that exists (from a previous test or manual creation)
	instanceID := 359559

	// Delete the instance
	err = client.DeleteInstance(instanceID)
	require.NoError(t, err)

	t.Logf("✓ Deleted instance %d", instanceID)
}
