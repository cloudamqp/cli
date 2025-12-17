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

// TestListVPCsVCR tests listing all VPCs
func TestListVPCsVCR(t *testing.T) {
	r, err := recorder.New("fixtures/list_vpcs")
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

	vpcs, err := client.ListVPCs()

	require.NoError(t, err)
	require.NotNil(t, vpcs)
	t.Logf("✓ Listed %d VPCs", len(vpcs))
}

// TestVPCLifecycleVCR tests the complete VPC lifecycle: create, get, update, delete
func TestVPCLifecycleVCR(t *testing.T) {
	r, err := recorder.New("fixtures/vpc_lifecycle")
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

	// Step 1: Create VPC
	t.Log("Step 1: Creating VPC")
	createReq := &VPCCreateRequest{
		Name:   "vcr-test-vpc",
		Region: "amazon-web-services::us-east-1",
		Subnet: "10.56.72.0/24",
		Tags:   []string{"test", "vcr"},
	}

	createResp, err := client.CreateVPC(createReq)
	require.NoError(t, err)
	require.NotNil(t, createResp)
	assert.NotZero(t, createResp.ID)
	t.Logf("✓ Created VPC with ID: %d", createResp.ID)

	vpcID := createResp.ID

	// Step 2: Get VPC
	t.Log("\nStep 2: Getting VPC details")
	vpc, err := client.GetVPC(vpcID)
	require.NoError(t, err)
	require.NotNil(t, vpc)
	assert.Equal(t, vpcID, vpc.ID)
	assert.Equal(t, "vcr-test-vpc", vpc.Name)
	t.Logf("✓ Got VPC: %s (subnet: %s)", vpc.Name, vpc.Subnet)

	// Step 3: Update VPC
	t.Log("\nStep 3: Updating VPC")
	updateReq := &VPCUpdateRequest{
		Name: "vcr-test-vpc-updated",
	}

	err = client.UpdateVPC(vpcID, updateReq)
	require.NoError(t, err)
	t.Logf("✓ Updated VPC name")

	// Step 4: Get updated VPC
	t.Log("\nStep 4: Verifying update")
	updatedVPC, err := client.GetVPC(vpcID)
	require.NoError(t, err)
	require.NotNil(t, updatedVPC)
	assert.Equal(t, "vcr-test-vpc-updated", updatedVPC.Name)
	t.Logf("✓ Verified VPC name: %s", updatedVPC.Name)

	// Step 5: Delete VPC
	t.Log("\nStep 5: Deleting VPC")
	err = client.DeleteVPC(vpcID)
	require.NoError(t, err)
	t.Logf("✓ Deleted VPC with ID: %d", vpcID)

	t.Log("\n✓ VPC lifecycle test completed successfully!")
}
