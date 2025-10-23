package cmd

import (
	"encoding/json"
	"fmt"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var instancePluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Manage RabbitMQ plugins",
	Long:  `List available RabbitMQ plugins for the instance.`,
}

var instancePluginsListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List plugins",
	Long:    `Retrieves all available RabbitMQ plugins.`,
	Example: `  cloudamqp instance manage 1234 plugins list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceID := currentInstanceID

		instanceAPIKey, err := getInstanceAPIKey(instanceID)
		if err != nil {
			return fmt.Errorf("failed to get instance API key: %w", err)
		}

		c := client.NewInstanceAPI(instanceAPIKey)

		plugins, err := c.ListPlugins()
		if err != nil {
			fmt.Printf("Error listing plugins: %v\n", err)
			return err
		}

		if len(plugins) == 0 {
			fmt.Println("No plugins found.")
			return nil
		}

		output, err := json.MarshalIndent(plugins, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("Plugins:\n%s\n", string(output))
		return nil
	},
}

func init() {
	instancePluginsCmd.AddCommand(instancePluginsListCmd)
}
