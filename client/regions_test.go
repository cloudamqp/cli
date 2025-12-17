package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListRegions(t *testing.T) {
	expectedRegions := []Region{
		{
			Provider:       "amazon-web-services",
			Region:         "us-east-1",
			Name:           "US East (N. Virginia)",
			HasSharedPlans: true,
		},
		{
			Provider:       "google-compute-engine",
			Region:         "us-central1",
			Name:           "US Central",
			HasSharedPlans: false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/regions", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedRegions)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	regions, err := client.ListRegions("")

	assert.NoError(t, err)
	assert.Len(t, regions, 2)
	assert.Equal(t, expectedRegions[0].Provider, regions[0].Provider)
	assert.Equal(t, expectedRegions[1].HasSharedPlans, regions[1].HasSharedPlans)
}

func TestListRegions_WithProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/regions", r.URL.Path)
		assert.Equal(t, "amazon-web-services", r.URL.Query().Get("provider"))

		regions := []Region{
			{
				Provider: "amazon-web-services",
				Region:   "us-east-1",
				Name:     "US East (N. Virginia)",
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(regions)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	regions, err := client.ListRegions("amazon-web-services")

	assert.NoError(t, err)
	assert.Len(t, regions, 1)
	assert.Equal(t, "amazon-web-services", regions[0].Provider)
}

func TestListPlans(t *testing.T) {
	expectedPlans := []Plan{
		{
			Name:    "lemming",
			Price:   0,
			Backend: "rabbitmq",
			Shared:  true,
		},
		{
			Name:    "bunny-1",
			Price:   19,
			Backend: "rabbitmq",
			Shared:  false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/plans", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedPlans)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	plans, err := client.ListPlans("")

	assert.NoError(t, err)
	assert.Len(t, plans, 2)
	assert.Equal(t, expectedPlans[0].Name, plans[0].Name)
	assert.Equal(t, expectedPlans[1].Price, plans[1].Price)
}

func TestListPlans_WithBackend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/plans", r.URL.Path)
		assert.Equal(t, "rabbitmq", r.URL.Query().Get("backend"))

		plans := []Plan{
			{
				Name:    "bunny-1",
				Backend: "rabbitmq",
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(plans)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-api-key", server.URL, "test")

	plans, err := client.ListPlans("rabbitmq")

	assert.NoError(t, err)
	assert.Len(t, plans, 1)
	assert.Equal(t, "rabbitmq", plans[0].Backend)
}
