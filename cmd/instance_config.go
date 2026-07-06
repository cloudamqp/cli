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
	Use:     "list <instance_id>",
	Short:   "List all configuration settings",
	Long:    `Retrieve and display all current RabbitMQ configuration settings.`,
	Example: `  cloudamqp instance config list 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		config, err := c.GetRabbitMQConfig(args[0])
		if err != nil {
			fmt.Printf("Error getting configuration: %v\n", err)
			return err
		}

		if len(config) == 0 {
			fmt.Println("No configuration found.")
			return nil
		}

		p, err := getPrinter(cmd)
		if err != nil {
			return err
		}

		headers := []string{"KEY", "VALUE"}
		rows := make([][]string, 0, len(config))
		for key, value := range config {
			rows = append(rows, []string{key, fmt.Sprintf("%v", value)})
		}
		p.PrintRecords(headers, rows)

		return nil
	},
}

var instanceConfigGetCmd = &cobra.Command{
	Use:     "get <instance_id> <setting>",
	Short:   "Get a specific configuration setting",
	Long:    `Retrieve a specific RabbitMQ configuration setting by name.`,
	Example: `  cloudamqp instance config get 1234 rabbit.heartbeat`,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		settingName := args[1]

		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		config, err := c.GetRabbitMQConfig(args[0])
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
	Use:   "set <instance_id> <setting> <value>",
	Short: "Set a configuration setting",
	Long:  `Update a RabbitMQ configuration setting. The value will be automatically converted to the appropriate type.`,
	Example: `  cloudamqp instance config set 1234 rabbit.heartbeat 120
  cloudamqp instance config set 1234 rabbit.vm_memory_high_watermark 0.8`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		settingName := args[1]
		settingValue := args[2]

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

		err = c.UpdateRabbitMQConfig(args[0], config)
		if err != nil {
			fmt.Printf("Error updating configuration: %v\n", err)
			return err
		}

		fmt.Printf("Configuration setting '%s' updated to: %v\n", settingName, value)
		return nil
	},
}

func init() {
	instanceConfigListCmd.ValidArgsFunction = completeInstances
	instanceConfigGetCmd.ValidArgsFunction = completeInstances
	instanceConfigSetCmd.ValidArgsFunction = completeInstances

	instanceConfigCmd.AddCommand(instanceConfigListCmd)
	instanceConfigCmd.AddCommand(instanceConfigGetCmd)
	instanceConfigCmd.AddCommand(instanceConfigSetCmd)
}
