package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var instanceConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage RabbitMQ configuration",
	Long:  `Get and update RabbitMQ configuration settings for the instance.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		cmd.SilenceUsage = true
		return fmt.Errorf("subcommand required")
	},
}

var instanceConfigListCmd = &cobra.Command{
	Use:     "list --id <instance_id>",
	Short:   "List all configuration settings",
	Long:    `Retrieve and display all current RabbitMQ configuration settings.`,
	Example: `  cloudamqp instance config list --id 1234`,
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

		c := client.New(apiKey, Version)

		config, err := c.GetRabbitMQConfig(idFlag)
		if err != nil {
			fmt.Printf("Error getting configuration: %v\n", err)
			return err
		}

		if len(config) == 0 {
			fmt.Println("No configuration found.")
			return nil
		}

		// Print table header
		fmt.Printf("%-40s %-30s\n", "KEY", "VALUE")
		fmt.Printf("%-40s %-30s\n", "---", "-----")

		// Print configuration data
		for key, value := range config {
			valueStr := fmt.Sprintf("%v", value)
			if len(valueStr) > 30 {
				valueStr = valueStr[:27] + "..."
			}
			fmt.Printf("%-40s %-30s\n", key, valueStr)
		}

		return nil
	},
}

var instanceConfigGetCmd = &cobra.Command{
	Use:     "get --id <instance_id> <setting>",
	Short:   "Get a specific configuration setting",
	Long:    `Retrieve a specific RabbitMQ configuration setting by name.`,
	Example: `  cloudamqp instance config get --id 1234 rabbit.heartbeat`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		idFlag, _ := cmd.Flags().GetString("id")
		if idFlag == "" {
			return fmt.Errorf("instance ID is required. Use --id flag")
		}

		settingName := args[0]

		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		config, err := c.GetRabbitMQConfig(idFlag)
		if err != nil {
			fmt.Printf("Error getting configuration: %v\n", err)
			return err
		}

		if value, exists := config[settingName]; exists {
			fmt.Printf("%s: %v\n", settingName, value)
		} else {
			fmt.Printf("Setting '%s' not found\n", settingName)
		}

		return nil
	},
}

var instanceConfigSetCmd = &cobra.Command{
	Use:   "set --id <instance_id> <setting> <value>",
	Short: "Set a configuration setting",
	Long:  `Update a RabbitMQ configuration setting. The value will be automatically converted to the appropriate type.`,
	Example: `  cloudamqp instance config set --id 1234 rabbit.heartbeat 120
  cloudamqp instance config set --id 1234 rabbit.vm_memory_high_watermark 0.8`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		idFlag, _ := cmd.Flags().GetString("id")
		if idFlag == "" {
			return fmt.Errorf("instance ID is required. Use --id flag")
		}

		settingName := args[0]
		settingValue := args[1]

		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		// Convert string value to appropriate type
		var value interface{}
		if strings.ToLower(settingValue) == "true" {
			value = true
		} else if strings.ToLower(settingValue) == "false" {
			value = false
		} else if strings.ToLower(settingValue) == "null" {
			value = nil
		} else if intVal, err := strconv.Atoi(settingValue); err == nil {
			value = intVal
		} else if floatVal, err := strconv.ParseFloat(settingValue, 64); err == nil {
			value = floatVal
		} else {
			value = settingValue
		}

		config := map[string]interface{}{
			settingName: value,
		}

		err = c.UpdateRabbitMQConfig(idFlag, config)
		if err != nil {
			fmt.Printf("Error updating configuration: %v\n", err)
			return err
		}

		fmt.Printf("Configuration setting '%s' updated to: %v\n", settingName, value)
		return nil
	},
}

func init() {
	// Add --id flag to all subcommands
	instanceConfigListCmd.Flags().StringP("id", "", "", "Instance ID (required)")
	instanceConfigListCmd.MarkFlagRequired("id")

	instanceConfigGetCmd.Flags().StringP("id", "", "", "Instance ID (required)")
	instanceConfigGetCmd.MarkFlagRequired("id")

	instanceConfigSetCmd.Flags().StringP("id", "", "", "Instance ID (required)")
	instanceConfigSetCmd.MarkFlagRequired("id")

	instanceConfigCmd.AddCommand(instanceConfigListCmd)
	instanceConfigCmd.AddCommand(instanceConfigGetCmd)
	instanceConfigCmd.AddCommand(instanceConfigSetCmd)
}
