package cmd

import (
	"fmt"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var instanceListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all CloudAMQP instances",
	Long:    `Retrieves and displays all CloudAMQP instances in your account.`,
	Example: `  cloudamqp instance list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		instances, err := c.ListInstances()
		if err != nil {
			fmt.Printf("Error listing instances: %v\n", err)
			return err
		}

		if len(instances) == 0 {
			fmt.Println("No instances found.")
			return nil
		}

		// Print table header
		fmt.Printf("%-10s %-20s %-15s %-30s\n", "ID", "NAME", "PLAN", "REGION")
		fmt.Printf("%-10s %-20s %-15s %-30s\n", "--", "----", "----", "------")

		// Print instance data
		for _, instance := range instances {
			fmt.Printf("%-10d %-20s %-15s %-30s\n", 
				instance.ID, 
				instance.Name, 
				instance.Plan, 
				instance.Region)
		}

		return nil
	},
}
