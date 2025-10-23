package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/term"
)

type Config struct {
	MainAPIKey   string            `json:"main_api_key,omitempty"`
	InstanceKeys map[string]string `json:"instance_keys,omitempty"`
}

func getAPIKey() (string, error) {
	// First, check environment variable
	if apiKey := os.Getenv("CLOUDAMQP_APIKEY"); apiKey != "" {
		return apiKey, nil
	}

	// Second, check config file
	config, err := loadConfig()
	if err == nil && config.MainAPIKey != "" {
		return config.MainAPIKey, nil
	}

	// If neither exists, prompt user and save to file
	fmt.Print("CloudAMQP API key not found. Please enter your API key: ")
	apiKey, err := readPassword()
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}

	if err := saveMainAPIKey(apiKey); err != nil {
		fmt.Printf("Warning: failed to save API key to config file: %v\n", err)
	} else {
		configPath, _ := getConfigPath()
		fmt.Printf("API key saved to %s\n", configPath)
	}

	return apiKey, nil
}

func getInstanceAPIKey(instanceID string) (string, error) {
	// First, check environment variable with instance ID suffix
	envKey := "CLOUDAMQP_INSTANCE_" + instanceID + "_APIKEY"
	if apiKey := os.Getenv(envKey); apiKey != "" {
		return apiKey, nil
	}

	// Second, check config file
	config, err := loadConfig()
	if err == nil && config.InstanceKeys != nil {
		if apiKey, exists := config.InstanceKeys[instanceID]; exists && apiKey != "" {
			return apiKey, nil
		}
	}

	return "", fmt.Errorf("instance API key not found for instance %s. Use 'cloudamqp instance get %s' to retrieve it", instanceID, instanceID)
}

func saveInstanceAPIKey(instanceID, apiKey string) error {
	config, err := loadConfig()
	if err != nil {
		config = &Config{}
	}

	if config.InstanceKeys == nil {
		config.InstanceKeys = make(map[string]string)
	}

	config.InstanceKeys[instanceID] = apiKey
	return saveConfig(config)
}

func saveMainAPIKey(apiKey string) error {
	config, err := loadConfig()
	if err != nil {
		config = &Config{}
	}

	config.MainAPIKey = apiKey
	return saveConfig(config)
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".cloudamqprc"), nil
}

func loadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// Handle legacy .cloudamqp file migration
	legacyPath := filepath.Join(filepath.Dir(configPath), ".cloudamqp")
	if _, err := os.Stat(legacyPath); err == nil {
		if err := migrateLegacyConfig(legacyPath, configPath); err != nil {
			fmt.Printf("Warning: failed to migrate legacy config: %v\n", err)
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func saveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

func migrateLegacyConfig(legacyPath, newPath string) error {
	data, err := os.ReadFile(legacyPath)
	if err != nil {
		return err
	}

	apiKey := strings.TrimSpace(string(data))
	if apiKey == "" {
		return fmt.Errorf("empty legacy config")
	}

	config := &Config{
		MainAPIKey:   apiKey,
		InstanceKeys: make(map[string]string),
	}

	if err := saveConfig(config); err != nil {
		return err
	}

	// Remove legacy file after successful migration
	os.Remove(legacyPath)
	fmt.Printf("Migrated legacy config to %s\n", newPath)
	return nil
}

func readPassword() (string, error) {
	if term.IsTerminal(int(syscall.Stdin)) {
		// Terminal input - hide password
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		fmt.Println() // Add newline after hidden input
		return string(bytePassword), nil
	} else {
		// Non-terminal input - read normally
		reader := bufio.NewReader(os.Stdin)
		password, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(password), nil
	}
}
