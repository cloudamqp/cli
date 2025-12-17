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

// TestListInstancesVCR tests listing all instances
func TestListInstancesVCR(t *testing.T) {
	r, err := recorder.New("fixtures/list_instances")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	// Add sanitization filters
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		i.Response.Body = sanitizeResponseBody(i.Response.Body)
		delete(i.Response.Headers, "Set-Cookie")
		return nil
	})

	apiKey := os.Getenv("CLOUDAMQP_APIKEY")
	if apiKey == "" {
		apiKey = "vcr-replay-mode"
	}

	httpClient := &http.Client{Transport: r}
	client := NewWithHTTPClient(apiKey, "https://customer.cloudamqp.com/api", "test", httpClient)

	instances, err := client.ListInstances()

	require.NoError(t, err)
	require.NotNil(t, instances)
	// Should have at least some instances (or empty list is OK)
	t.Logf("✓ Listed %d instances", len(instances))
}

// TestGetInstanceVCR tests getting a specific instance
func TestGetInstanceVCR(t *testing.T) {
	r, err := recorder.New("fixtures/get_instance")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		i.Response.Body = sanitizeResponseBody(i.Response.Body)
		delete(i.Response.Headers, "Set-Cookie")
		return nil
	})

	apiKey := os.Getenv("CLOUDAMQP_APIKEY")
	if apiKey == "" {
		apiKey = "vcr-replay-mode"
	}

	httpClient := &http.Client{Transport: r}
	client := NewWithHTTPClient(apiKey, "https://customer.cloudamqp.com/api", "test", httpClient)

	// Use an existing instance ID (should match cassette)
	instanceID := 359563

	instance, err := client.GetInstance(instanceID)

	require.NoError(t, err)
	require.NotNil(t, instance)
	assert.Equal(t, instanceID, instance.ID)
	t.Logf("✓ Got instance %d: %s (plan: %s)", instance.ID, instance.Name, instance.Plan)
}

// TestListRegionsVCR tests listing all regions
func TestListRegionsVCR(t *testing.T) {
	r, err := recorder.New("fixtures/list_regions")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		i.Response.Body = sanitizeResponseBody(i.Response.Body)
		delete(i.Response.Headers, "Set-Cookie")
		return nil
	})

	apiKey := os.Getenv("CLOUDAMQP_APIKEY")
	if apiKey == "" {
		apiKey = "vcr-replay-mode"
	}

	httpClient := &http.Client{Transport: r}
	client := NewWithHTTPClient(apiKey, "https://customer.cloudamqp.com/api", "test", httpClient)

	regions, err := client.ListRegions("")

	require.NoError(t, err)
	require.NotEmpty(t, regions)
	t.Logf("✓ Listed %d regions", len(regions))

	// Verify some expected fields
	if len(regions) > 0 {
		assert.NotEmpty(t, regions[0].Name)
		assert.NotEmpty(t, regions[0].Provider)
		t.Logf("  Example: %s (%s)", regions[0].Name, regions[0].Provider)
	}
}

// TestListPlansVCR tests listing all plans
func TestListPlansVCR(t *testing.T) {
	r, err := recorder.New("fixtures/list_plans")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		i.Response.Body = sanitizeResponseBody(i.Response.Body)
		delete(i.Response.Headers, "Set-Cookie")
		return nil
	})

	apiKey := os.Getenv("CLOUDAMQP_APIKEY")
	if apiKey == "" {
		apiKey = "vcr-replay-mode"
	}

	httpClient := &http.Client{Transport: r}
	client := NewWithHTTPClient(apiKey, "https://customer.cloudamqp.com/api", "test", httpClient)

	plans, err := client.ListPlans("")

	require.NoError(t, err)
	require.NotEmpty(t, plans)
	t.Logf("✓ Listed %d plans", len(plans))

	// Verify some expected fields and find bunny-1
	bunnyFound := false
	for _, plan := range plans {
		assert.NotEmpty(t, plan.Name)
		if plan.Name == "bunny-1" {
			bunnyFound = true
			t.Logf("  Found bunny-1 plan")
		}
	}
	assert.True(t, bunnyFound, "Should find bunny-1 plan")
}
