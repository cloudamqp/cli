package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type cacheEntry struct {
	Data      json.RawMessage `json:"data"`
	Timestamp int64           `json:"timestamp"`
}

// getCacheDir returns the cache directory path, creating it if needed
func getCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(homeDir, ".cache", "cloudamqp")
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return "", err
	}

	return cacheDir, nil
}

// formatTTL converts duration to a short string for filename
func formatTTL(ttl time.Duration) string {
	hours := int(ttl.Hours())
	if hours >= 24 {
		days := hours / 24
		if days == 1 {
			return "24h"
		}
		return fmt.Sprintf("%dd", days)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	minutes := int(ttl.Minutes())
	return fmt.Sprintf("%dm", minutes)
}

// getCacheFilename returns the cache filename with TTL information
func getCacheFilename(key string, ttl time.Duration) string {
	return fmt.Sprintf("cache_%s_ttl_%s.json", formatTTL(ttl), key)
}

// getCachedData retrieves cached data if it exists and is not expired
func getCachedData(key string, ttl time.Duration) (json.RawMessage, bool) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, false
	}

	cachePath := filepath.Join(cacheDir, getCacheFilename(key, ttl))
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, false
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	// Check if cache is expired
	if time.Now().Unix()-entry.Timestamp > int64(ttl.Seconds()) {
		return nil, false
	}

	return entry.Data, true
}

// setCachedData stores data in the cache with current timestamp
func setCachedData(key string, ttl time.Duration, data interface{}) error {
	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	entry := cacheEntry{
		Data:      jsonData,
		Timestamp: time.Now().Unix(),
	}

	entryData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	cachePath := filepath.Join(cacheDir, getCacheFilename(key, ttl))
	return os.WriteFile(cachePath, entryData, 0600)
}

// Cache TTL settings
const (
	plansCacheTTL     = 1 * time.Hour   // Plans rarely change
	regionsCacheTTL   = 1 * time.Hour   // Regions rarely change
	versionsCacheTTL  = 1 * time.Hour   // Versions rarely change
	instancesCacheTTL = 1 * time.Minute // Instances change frequently
	vpcsCacheTTL      = 1 * time.Minute // VPCs change frequently
)
