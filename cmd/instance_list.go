package cmd

import (
	"fmt"
	"os"
	"strconv"

	"cloudamqp-cli/client"
	"cloudamqp-cli/internal/table"
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

		c := client.New(apiKey, Version)

		instances, err := c.ListInstances()
		if err != nil {
			fmt.Printf("Error listing instances: %v\n", err)
			return err
		}

		if len(instances) == 0 {
			fmt.Println("No instances found.")
			return nil
		}

		// Create table and populate data
		t := table.New(os.Stdout, "ID", "NAME", "PLAN", "REGION")
		for _, instance := range instances {
			t.AddRow(
				strconv.Itoa(instance.ID),
				instance.Name,
				instance.Plan,
				instance.Region,
			)
		}
		t.Print()

		return nil
	},
}
