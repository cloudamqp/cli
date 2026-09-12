package cmd

import (
	"fmt"

	"cloudamqp-cli/client"
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
	Use:     "list <instance_id>",
	Short:   "List plugins",
	Long:    `Retrieves all available RabbitMQ plugins.`,
	Example: `  cloudamqp instance plugins list 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		plugins, err := c.ListPlugins(args[0])
		if err != nil {
			fmt.Printf("Error listing plugins: %v\n", err)
			return err
		}

		if len(plugins) == 0 {
			fmt.Println("No plugins found.")
			return nil
		}

		p, err := getPrinter(cmd)
		if err != nil {
			return err
		}

		headers := []string{"NAME", "ENABLED"}
		rows := make([][]string, len(plugins))
		for i, plugin := range plugins {
			enabled := "No"
			if plugin.Enabled {
				enabled = "Yes"
			}
			rows[i] = []string{plugin.Name, enabled}
		}
		p.PrintRecords(headers, rows)

		return nil
	},
}

var instancePluginsEnableCmd = &cobra.Command{
	Use:     "enable <instance_id> <plugin_name>",
	Short:   "Enable a plugin",
	Long:    `Enables a RabbitMQ plugin on the instance.`,
	Args:    cobra.ExactArgs(2),
	Example: `  cloudamqp instance plugins enable 1234 rabbitmq_top`,
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceID := args[0]
		pluginName := args[1]

		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		err = c.EnablePlugin(instanceID, pluginName)
		if err != nil {
			fmt.Printf("Error enabling plugin '%s': %v\n", pluginName, err)
			return err
		}

		fmt.Printf("Plugin '%s' enabled successfully.\n", pluginName)
		return nil
	},
}

var instancePluginsDisableCmd = &cobra.Command{
	Use:     "disable <instance_id> <plugin_name>",
	Short:   "Disable a plugin",
	Long:    `Disables a RabbitMQ plugin on the instance.`,
	Args:    cobra.ExactArgs(2),
	Example: `  cloudamqp instance plugins disable 1234 rabbitmq_top`,
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceID := args[0]
		pluginName := args[1]

		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		err = c.DisablePlugin(instanceID, pluginName)
		if err != nil {
			fmt.Printf("Error disabling plugin '%s': %v\n", pluginName, err)
			return err
		}

		fmt.Printf("Plugin '%s' disabled successfully.\n", pluginName)
		return nil
	},
}

func init() {
	instancePluginsListCmd.ValidArgsFunction = completeInstances
	instancePluginsEnableCmd.ValidArgsFunction = completeInstances
	instancePluginsDisableCmd.ValidArgsFunction = completeInstances

	// Add all commands to plugins
	instancePluginsCmd.AddCommand(instancePluginsListCmd)
	instancePluginsCmd.AddCommand(instancePluginsEnableCmd)
	instancePluginsCmd.AddCommand(instancePluginsDisableCmd)
}
