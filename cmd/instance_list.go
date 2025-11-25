package cmd

import (
	"fmt"
	"strconv"

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

		// Calculate column widths
		idWidth := len("ID")
		nameWidth := len("NAME")
		planWidth := len("PLAN")
		regionWidth := len("REGION")

		for _, instance := range instances {
			idLen := len(strconv.Itoa(instance.ID))
			if idLen > idWidth {
				idWidth = idLen
			}
			if len(instance.Name) > nameWidth {
				nameWidth = len(instance.Name)
			}
			if len(instance.Plan) > planWidth {
				planWidth = len(instance.Plan)
			}
			if len(instance.Region) > regionWidth {
				regionWidth = len(instance.Region)
			}
		}

		// Add padding
		idWidth += 2
		nameWidth += 2
		planWidth += 2
		regionWidth += 2

		// Create format strings
		headerFormat := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%-%ds\n", idWidth, nameWidth, planWidth, regionWidth)
		rowFormat := fmt.Sprintf("%%-%dd %%-%ds %%-%ds %%-%ds\n", idWidth, nameWidth, planWidth, regionWidth)

		// Print table header
		fmt.Printf(headerFormat, "ID", "NAME", "PLAN", "REGION")
		fmt.Printf(headerFormat, "--", "----", "----", "------")

		// Print instance data
		for _, instance := range instances {
			fmt.Printf(rowFormat,
				instance.ID,
				instance.Name,
				instance.Plan,
				instance.Region)
		}

		return nil
	},
}
