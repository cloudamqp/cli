package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var instanceGetCmd = &cobra.Command{
	Use:     "get <id>",
	Short:   "Get details of a specific CloudAMQP instance",
	Long:    `Retrieves and displays detailed information about a specific CloudAMQP instance.`,
	Example: `  cloudamqp instance get 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		instanceID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid instance ID: %v", err)
		}

		c := client.New(apiKey)

		instance, err := c.GetInstance(instanceID)
		if err != nil {
			fmt.Printf("Error getting instance: %v\n", err)
			return err
		}

		// Save instance API key to config if it exists
		if instance.APIKey != "" {
			if err := saveInstanceAPIKey(strconv.Itoa(instanceID), instance.APIKey); err != nil {
				fmt.Printf("Warning: failed to save instance API key to config: %v\n", err)
			} else {
				fmt.Printf("Instance API key saved for instance %d\n", instanceID)
			}
		}

		output, err := json.MarshalIndent(instance, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("Instance details:\n%s\n", string(output))
		return nil
	},
}
