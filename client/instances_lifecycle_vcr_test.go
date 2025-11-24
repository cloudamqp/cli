package client

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v2/cassette"
	"gopkg.in/dnaeon/go-vcr.v2/recorder"
)

// TestInstanceLifecycleVCR tests the complete instance lifecycle:
// 1. Create instance with bunny-1 plan
// 2. Update instance from bunny-1 to hare-1 plan
// 3. Delete instance
func TestInstanceLifecycleVCR(t *testing.T) {
	// Create a VCR recorder for the lifecycle test
	r, err := recorder.New("fixtures/instance_lifecycle")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	// Add filter to sanitize sensitive data in cassettes
	r.AddFilter(func(i *cassette.Interaction) error {
		// Sanitize Authorization header
		delete(i.Request.Headers, "Authorization")

		// Sanitize sensitive data in response body (API keys, passwords)
		i.Response.Body = sanitizeResponseBody(i.Response.Body)

		// Remove session cookies
		delete(i.Response.Headers, "Set-Cookie")

		return nil
	})

	// Get API key from environment
	apiKey := os.Getenv("CLOUDAMQP_APIKEY")
	if apiKey == "" && r.Mode() != recorder.ModeReplaying {
		t.Skip("CLOUDAMQP_APIKEY environment variable not set, skipping test")
	}

	// Create HTTP client with VCR recorder as transport
	httpClient := &http.Client{Transport: r}

	// Create client with VCR HTTP client
	client := NewWithHTTPClient(apiKey, "https://customer.cloudamqp.com/api", httpClient)

	// Step 1: Create instance with bunny-1 plan
	t.Log("Step 1: Creating instance with bunny-1 plan")
	createReq := &InstanceCreateRequest{
		Name:   "vcr-lifecycle-test",
		Plan:   "bunny-1",
		Region: "amazon-web-services::us-east-1",
		Tags:   []string{"test", "vcr", "lifecycle"},
	}

	createResp, err := client.CreateInstance(createReq)
	require.NoError(t, err, "Failed to create instance")
	require.NotNil(t, createResp)
	require.NotZero(t, createResp.ID, "Instance ID should not be zero")
	require.NotEmpty(t, createResp.URL, "Instance URL should not be empty")
	require.NotEmpty(t, createResp.APIKey, "Instance API key should not be empty")

	instanceID := createResp.ID
	t.Logf("✓ Created instance with ID: %d", instanceID)
	t.Logf("  Plan: bunny-1")
	t.Logf("  URL: %s", createResp.URL)

	// Wait for instance to be ready (only in recording mode)
	if r.Mode() != recorder.ModeReplaying {
		t.Log("Waiting for instance to be ready...")
		maxAttempts := 60 // 5 minutes max
		for i := 0; i < maxAttempts; i++ {
			instance, err := client.GetInstance(instanceID)
			if err == nil && instance != nil && instance.Ready {
				t.Logf("✓ Instance is ready after %d seconds", (i+1)*5)
				break
			}
			if i == maxAttempts-1 {
				t.Fatal("Instance did not become ready in time")
			}
			t.Logf("  Attempt %d/%d: Instance not ready yet, waiting...", i+1, maxAttempts)
			time.Sleep(5 * time.Second)
		}
	}

	// Step 2: Get instance details to verify initial state
	t.Log("\nStep 2: Getting instance details to verify bunny-1 plan")
	instance, err := client.GetInstance(instanceID)
	require.NoError(t, err, "Failed to get instance")
	require.NotNil(t, instance)
	assert.Equal(t, instanceID, instance.ID)
	assert.Equal(t, "bunny-1", instance.Plan, "Plan should be bunny-1")
	assert.True(t, instance.Ready, "Instance should be ready")
	t.Logf("✓ Verified instance plan: %s (Ready: %v)", instance.Plan, instance.Ready)

	// Additional wait for cluster to be fully configured (only in recording mode)
	if r.Mode() != recorder.ModeReplaying {
		t.Log("Waiting additional 30 seconds for cluster to be fully configured...")
		time.Sleep(30 * time.Second)
	}

	// Step 3: Update instance from bunny-1 to hare-1
	t.Log("\nStep 3: Updating instance from bunny-1 to hare-1")
	updateReq := &InstanceUpdateRequest{
		Plan: "hare-1",
	}

	err = client.UpdateInstance(instanceID, updateReq)
	require.NoError(t, err, "Failed to update instance")
	t.Logf("✓ Updated instance plan to hare-1")

	// Wait a bit for update to process (only in recording mode)
	if r.Mode() != recorder.ModeReplaying {
		t.Log("Waiting 3 seconds for update to process...")
		time.Sleep(3 * time.Second)
	}

	// Step 4: Get instance details to verify update
	t.Log("\nStep 4: Getting instance details to verify hare-1 plan")
	updatedInstance, err := client.GetInstance(instanceID)
	require.NoError(t, err, "Failed to get updated instance")
	require.NotNil(t, updatedInstance)
	assert.Equal(t, "hare-1", updatedInstance.Plan, "Plan should be updated to hare-1")
	t.Logf("✓ Verified updated instance plan: %s", updatedInstance.Plan)

	// Step 5: Delete instance
	t.Log("\nStep 5: Deleting instance")
	err = client.DeleteInstance(instanceID)
	require.NoError(t, err, "Failed to delete instance")
	t.Logf("✓ Deleted instance with ID: %d", instanceID)

	t.Log("\n✓ Instance lifecycle test completed successfully!")
}

// TestCreateInstanceBunny1VCR tests creating an instance with bunny-1 plan
func TestCreateInstanceBunny1VCR(t *testing.T) {
	// Create a VCR recorder
	r, err := recorder.New("fixtures/create_instance_bunny1")
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

	// Get API key from environment
	apiKey := os.Getenv("CLOUDAMQP_APIKEY")
	if apiKey == "" && r.Mode() != recorder.ModeReplaying {
		t.Skip("CLOUDAMQP_APIKEY environment variable not set, skipping test")
	}

	// Create HTTP client with VCR recorder as transport
	httpClient := &http.Client{Transport: r}
	client := NewWithHTTPClient(apiKey, "https://customer.cloudamqp.com/api", httpClient)

	// Create instance request with bunny-1 plan
	req := &InstanceCreateRequest{
		Name:   "vcr-bunny1-test",
		Plan:   "bunny-1",
		Region: "amazon-web-services::us-east-1",
		Tags:   []string{"test", "vcr", "bunny1"},
	}

	// Execute the create instance request
	resp, err := client.CreateInstance(req)

	// Verify the response
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotZero(t, resp.ID)
	assert.NotEmpty(t, resp.URL)
	assert.NotEmpty(t, resp.APIKey)

	t.Logf("Created instance with ID: %d", resp.ID)
	t.Logf("Plan: bunny-1")
	t.Logf("Instance URL: %s", resp.URL)
}

// TestUpdateInstancePlanVCR tests updating an instance plan from bunny-1 to hare-1
func TestUpdateInstancePlanVCR(t *testing.T) {
	// Create a VCR recorder
	r, err := recorder.New("fixtures/update_instance_plan")
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

	// Get API key from environment
	apiKey := os.Getenv("CLOUDAMQP_APIKEY")
	if apiKey == "" && r.Mode() != recorder.ModeReplaying {
		t.Skip("CLOUDAMQP_APIKEY environment variable not set, skipping test")
	}

	// Create HTTP client with VCR recorder as transport
	httpClient := &http.Client{Transport: r}
	client := NewWithHTTPClient(apiKey, "https://customer.cloudamqp.com/api", httpClient)

	// First create an instance with bunny-1 (only in recording mode)
	var instanceID int
	if r.Mode() != recorder.ModeReplaying {
		createReq := &InstanceCreateRequest{
			Name:   "vcr-update-test",
			Plan:   "bunny-1",
			Region: "amazon-web-services::us-east-1",
		}
		createResp, err := client.CreateInstance(createReq)
		require.NoError(t, err)
		instanceID = createResp.ID
		t.Logf("Created instance with ID: %d (bunny-1)", instanceID)

		// Wait for instance to be ready and fully configured
		t.Log("Waiting for instance to be ready and fully configured...")
		maxAttempts := 120 // 10 minutes max
		var readyInstance *Instance
		for i := 0; i < maxAttempts; i++ {
			inst, err := client.GetInstance(instanceID)
			if err == nil && inst != nil && inst.Ready {
				readyInstance = inst
				t.Logf("✓ Instance is ready after %d seconds", (i+1)*5)
				break
			}
			if i == maxAttempts-1 {
				t.Fatal("Instance did not become ready in time")
			}
			if i%6 == 0 { // Log every 30 seconds
				t.Logf("  Still waiting... (%d seconds elapsed)", (i+1)*5)
			}
			time.Sleep(5 * time.Second)
		}

		// Additional wait to ensure cluster is fully configured for updates
		t.Log("Waiting additional 60 seconds for cluster to be fully configured for updates...")
		time.Sleep(60 * time.Second)

		// Verify instance plan before update
		t.Logf("Instance before update - Plan: %s, Ready: %v", readyInstance.Plan, readyInstance.Ready)
	} else {
		// When replaying, use the ID from the recording
		// This will be determined by reading the cassette interactions
		instanceID = 359289 // This should match your recorded instance ID
	}

	// Update instance to hare-1 plan
	updateReq := &InstanceUpdateRequest{
		Plan: "hare-1",
	}

	err = client.UpdateInstance(instanceID, updateReq)
	assert.NoError(t, err)

	t.Logf("Updated instance %d from bunny-1 to hare-1", instanceID)

	// Verify the update
	instance, err := client.GetInstance(instanceID)
	assert.NoError(t, err)
	assert.NotNil(t, instance)
	if instance != nil {
		assert.Equal(t, "hare-1", instance.Plan)
		t.Logf("Verified instance plan: %s", instance.Plan)
	}

	// Clean up (only in recording mode)
	if r.Mode() != recorder.ModeReplaying {
		time.Sleep(3 * time.Second)
		err = client.DeleteInstance(instanceID)
		if err != nil {
			t.Logf("Warning: Failed to delete instance %d: %v", instanceID, err)
		} else {
			t.Logf("Cleaned up instance %d", instanceID)
		}
	}
}

// TestDeleteInstanceVCR tests deleting an instance
func TestDeleteInstanceVCR(t *testing.T) {
	// Create a VCR recorder
	r, err := recorder.New("fixtures/delete_instance")
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

	// Get API key from environment
	apiKey := os.Getenv("CLOUDAMQP_APIKEY")
	if apiKey == "" && r.Mode() != recorder.ModeReplaying {
		t.Skip("CLOUDAMQP_APIKEY environment variable not set, skipping test")
	}

	// Create HTTP client with VCR recorder as transport
	httpClient := &http.Client{Transport: r}
	client := NewWithHTTPClient(apiKey, "https://customer.cloudamqp.com/api", httpClient)

	// First create an instance (only in recording mode)
	var instanceID int
	if r.Mode() != recorder.ModeReplaying {
		createReq := &InstanceCreateRequest{
			Name:   "vcr-delete-test",
			Plan:   "lemur",
			Region: "amazon-web-services::us-east-1",
		}
		createResp, err := client.CreateInstance(createReq)
		require.NoError(t, err)
		instanceID = createResp.ID
		t.Logf("Created instance with ID: %d for deletion test", instanceID)

		// Wait for instance to be ready
		time.Sleep(5 * time.Second)
	} else {
		// When replaying, use the ID from the recording
		instanceID = 359287 // This should match your recorded instance ID
	}

	// Delete the instance
	err = client.DeleteInstance(instanceID)
	assert.NoError(t, err)

	t.Logf("✓ Deleted instance with ID: %d", instanceID)
}
