package cmd

import (
	"fmt"
	"os"

	"cloudamqp-cli/client"
	"cloudamqp-cli/internal/table"
	"github.com/spf13/cobra"
)

var instancePluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Manage RabbitMQ plugins",
	Long:  `List, enable, and disable RabbitMQ plugins for the instance.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		cmd.SilenceUsage = true
		return fmt.Errorf("subcommand required")
	},
}

var instancePluginsListCmd = &cobra.Command{
	Use:     "list --id <instance_id>",
	Short:   "List plugins",
	Long:    `Retrieves all available RabbitMQ plugins.`,
	Example: `  cloudamqp instance plugins list --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		idFlag, _ := cmd.Flags().GetString("id")
		if idFlag == "" {
			return fmt.Errorf("instance ID is required. Use --id flag")
		}

		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		plugins, err := c.ListPlugins(idFlag)
		if err != nil {
			fmt.Printf("Error listing plugins: %v\n", err)
			return err
		}

		if len(plugins) == 0 {
			fmt.Println("No plugins found.")
			return nil
		}

		// Create table and populate data
		t := table.New(os.Stdout, "NAME", "ENABLED")
		for _, plugin := range plugins {
			enabled := "No"
			if plugin.Enabled {
				enabled = "Yes"
			}
			t.AddRow(plugin.Name, enabled)
		}
		t.Print()

		return nil
	},
}

var instancePluginsEnableCmd = &cobra.Command{
	Use:     "enable <plugin_name> --id <instance_id>",
	Short:   "Enable a plugin",
	Long:    `Enables a RabbitMQ plugin on the instance.`,
	Args:    cobra.ExactArgs(1),
	Example: `  cloudamqp instance plugins enable rabbitmq_top --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pluginName := args[0]
		idFlag, _ := cmd.Flags().GetString("id")
		if idFlag == "" {
			return fmt.Errorf("instance ID is required. Use --id flag")
		}

		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		err = c.EnablePlugin(idFlag, pluginName)
		if err != nil {
			fmt.Printf("Error enabling plugin '%s': %v\n", pluginName, err)
			return err
		}

		fmt.Printf("Plugin '%s' enabled successfully.\n", pluginName)
		return nil
	},
}

var instancePluginsDisableCmd = &cobra.Command{
	Use:     "disable <plugin_name> --id <instance_id>",
	Short:   "Disable a plugin",
	Long:    `Disables a RabbitMQ plugin on the instance.`,
	Args:    cobra.ExactArgs(1),
	Example: `  cloudamqp instance plugins disable rabbitmq_top --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pluginName := args[0]
		idFlag, _ := cmd.Flags().GetString("id")
		if idFlag == "" {
			return fmt.Errorf("instance ID is required. Use --id flag")
		}

		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		err = c.DisablePlugin(idFlag, pluginName)
		if err != nil {
			fmt.Printf("Error disabling plugin '%s': %v\n", pluginName, err)
			return err
		}

		fmt.Printf("Plugin '%s' disabled successfully.\n", pluginName)
		return nil
	},
}

func init() {
	// Add --id flag to all plugins commands
	instancePluginsListCmd.Flags().StringP("id", "", "", "Instance ID (required)")
	instancePluginsListCmd.MarkFlagRequired("id")

	instancePluginsEnableCmd.Flags().StringP("id", "", "", "Instance ID (required)")
	instancePluginsEnableCmd.MarkFlagRequired("id")

	instancePluginsDisableCmd.Flags().StringP("id", "", "", "Instance ID (required)")
	instancePluginsDisableCmd.MarkFlagRequired("id")

	// Add all commands to plugins
	instancePluginsCmd.AddCommand(instancePluginsListCmd)
	instancePluginsCmd.AddCommand(instancePluginsEnableCmd)
	instancePluginsCmd.AddCommand(instancePluginsDisableCmd)
}
