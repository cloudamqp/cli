package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func getAPIKey() (string, error) {
	// First, check environment variable
	if apiKey := os.Getenv("CLOUDAMQP_APIKEY"); apiKey != "" {
		return apiKey, nil
	}

	// Second, check config file
	apiKey, err := loadAPIKey()
	if err == nil && apiKey != "" {
		return apiKey, nil
	}

	// If neither exists, prompt user and save to file
	fmt.Print("CloudAMQP API key not found. Please enter your API key: ")
	apiKey, err = readPassword()
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}

	if err := saveAPIKey(apiKey); err != nil {
		fmt.Printf("Warning: failed to save API key to config file: %v\n", err)
	} else {
		configPath, _ := getConfigPath()
		fmt.Printf("API key saved to %s\n", configPath)
	}

	return apiKey, nil
}

func saveAPIKey(apiKey string) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, []byte(strings.TrimSpace(apiKey)), 0600)
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".cloudamqprc"), nil
}

func loadAPIKey() (string, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
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