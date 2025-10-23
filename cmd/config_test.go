package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_LoadAndSave(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cloudamqp-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Test saving config
	config := &Config{
		MainAPIKey: "test-main-key",
		InstanceKeys: map[string]string{
			"1234": "test-instance-key-1234",
			"5678": "test-instance-key-5678",
		},
	}

	err = saveConfig(config)
	assert.NoError(t, err)

	// Verify file exists
	configPath := filepath.Join(tempDir, ".cloudamqprc")
	assert.FileExists(t, configPath)

	// Test loading config
	loadedConfig, err := loadConfig()
	assert.NoError(t, err)
	assert.Equal(t, config.MainAPIKey, loadedConfig.MainAPIKey)
	assert.Equal(t, config.InstanceKeys, loadedConfig.InstanceKeys)
}

func TestConfig_SaveMainAPIKey(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cloudamqp-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Test saving main API key
	err = saveMainAPIKey("new-main-key")
	assert.NoError(t, err)

	// Verify it's saved correctly
	config, err := loadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "new-main-key", config.MainAPIKey)
	assert.Nil(t, config.InstanceKeys) // Should be nil/empty since we only set main key
}

func TestConfig_SaveInstanceAPIKey(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cloudamqp-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// First save main key
	err = saveMainAPIKey("main-key")
	assert.NoError(t, err)

	// Then save instance key
	err = saveInstanceAPIKey("1234", "instance-key-1234")
	assert.NoError(t, err)

	// Verify both are saved
	config, err := loadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "main-key", config.MainAPIKey)
	assert.Equal(t, "instance-key-1234", config.InstanceKeys["1234"])
}

func TestConfig_GetInstanceAPIKey_Found(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cloudamqp-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Save instance key
	err = saveInstanceAPIKey("1234", "test-instance-key")
	assert.NoError(t, err)

	// Test getting instance key
	apiKey, err := getInstanceAPIKey("1234")
	assert.NoError(t, err)
	assert.Equal(t, "test-instance-key", apiKey)
}

func TestConfig_GetInstanceAPIKey_NotFound(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cloudamqp-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Test getting non-existent instance key
	_, err = getInstanceAPIKey("9999")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "instance API key not found for instance 9999")
}

func TestConfig_GetInstanceAPIKey_FromEnv(t *testing.T) {
	// Set environment variable
	os.Setenv("CLOUDAMQP_INSTANCE_1234_APIKEY", "env-instance-key")
	defer os.Unsetenv("CLOUDAMQP_INSTANCE_1234_APIKEY")

	// Test getting instance key from env
	apiKey, err := getInstanceAPIKey("1234")
	assert.NoError(t, err)
	assert.Equal(t, "env-instance-key", apiKey)
}

func TestConfig_GetAPIKey_FromEnv(t *testing.T) {
	// Set environment variable
	os.Setenv("CLOUDAMQP_APIKEY", "env-main-key")
	defer os.Unsetenv("CLOUDAMQP_APIKEY")

	// Test getting main key from env
	apiKey, err := getAPIKey()
	assert.NoError(t, err)
	assert.Equal(t, "env-main-key", apiKey)
}

func TestConfig_GetAPIKey_FromFile(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cloudamqp-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Make sure env var is not set
	os.Unsetenv("CLOUDAMQP_APIKEY")

	// Save main key
	err = saveMainAPIKey("file-main-key")
	assert.NoError(t, err)

	// Test getting main key from file
	apiKey, err := getAPIKey()
	assert.NoError(t, err)
	assert.Equal(t, "file-main-key", apiKey)
}

func TestConfig_MigrateLegacyConfig(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cloudamqp-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create legacy config file
	legacyPath := filepath.Join(tempDir, ".cloudamqp")
	err = os.WriteFile(legacyPath, []byte("legacy-api-key"), 0600)
	assert.NoError(t, err)

	// Test migration by loading config (should trigger migration)
	config, err := loadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "legacy-api-key", config.MainAPIKey)

	// Verify legacy file is removed
	assert.NoFileExists(t, legacyPath)

	// Verify new config file exists
	newConfigPath := filepath.Join(tempDir, ".cloudamqprc")
	assert.FileExists(t, newConfigPath)
}

func TestConfig_InvalidJSON(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cloudamqp-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create invalid JSON config file
	configPath := filepath.Join(tempDir, ".cloudamqprc")
	err = os.WriteFile(configPath, []byte("invalid json"), 0600)
	assert.NoError(t, err)

	// Test loading invalid config
	_, err = loadConfig()
	assert.Error(t, err)
}

func TestConfig_EmptyConfig(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cloudamqp-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create empty config file
	configPath := filepath.Join(tempDir, ".cloudamqprc")
	err = os.WriteFile(configPath, []byte("{}"), 0600)
	assert.NoError(t, err)

	// Test loading empty config
	config, err := loadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "", config.MainAPIKey)
	assert.Nil(t, config.InstanceKeys)
}

func TestConfig_GetConfigPath(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cloudamqp-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Test getting config path
	configPath, err := getConfigPath()
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(tempDir, ".cloudamqprc"), configPath)
}

func TestConfig_FilePermissions(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cloudamqp-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Save config
	config := &Config{MainAPIKey: "test-key"}
	err = saveConfig(config)
	assert.NoError(t, err)

	// Check file permissions
	configPath := filepath.Join(tempDir, ".cloudamqprc")
	info, err := os.Stat(configPath)
	assert.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}
