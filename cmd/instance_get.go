package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var instanceGetCmd = &cobra.Command{
	Use:     "get --id <id>",
	Short:   "Get details of a specific CloudAMQP instance",
	Long:    `Retrieves and displays detailed information about a specific CloudAMQP instance.`,
	Example: `  cloudamqp instance get --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		idFlag, _ := cmd.Flags().GetString("id")
		if idFlag == "" {
			return fmt.Errorf("instance ID is required. Use --id flag")
		}

		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		instanceID, err := strconv.Atoi(idFlag)
		if err != nil {
			return fmt.Errorf("invalid instance ID: %v", err)
		}

		c := client.New(apiKey, Version)

		instance, err := c.GetInstance(instanceID)
		if err != nil {
			fmt.Printf("Error getting instance: %v\n", err)
			return err
		}

		// Format output as "Name = Value"
		fmt.Printf("Name = %s\n", instance.Name)
		fmt.Printf("Plan = %s\n", instance.Plan)
		fmt.Printf("Region = %s\n", instance.Region)
		fmt.Printf("Tags = %s\n", strings.Join(instance.Tags, ","))
		fmt.Printf("Hostname = %s\n", instance.HostnameExternal)
		ready := "No"
		if instance.Ready {
			ready = "Yes"
		}
		fmt.Printf("Ready = %s\n", ready)

		return nil
	},
}

func init() {
	instanceGetCmd.Flags().StringP("id", "", "", "Instance ID (required)")
	instanceGetCmd.MarkFlagRequired("id")
	instanceGetCmd.RegisterFlagCompletionFunc("id", completeInstanceIDFlag)
}
