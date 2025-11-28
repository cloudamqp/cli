package cmd

import (
	"encoding/json"
	"fmt"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var rotateKeyCmd = &cobra.Command{
	Use:     "rotate-key",
	Short:   "Rotate API key",
	Long:    `Removes the current API key and creates a new one with matching permissions.`,
	Example: `  cloudamqp rotate-key`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		resp, err := c.RotateAPIKey()
		if err != nil {
			return err
		}

		output, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Println(string(output))

		// Update local config file with new key
		_ = saveAPIKey(resp.APIKey)

		return nil
	},
}
