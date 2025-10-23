package cmd

import (
	"encoding/json"
	"fmt"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var backendFilter string

var plansCmd = &cobra.Command{
	Use:   "plans",
	Short: "List available plans",
	Long:  `Retrieves all available subscription plans, optionally filtered by backend.`,
	Example: `  cloudamqp plans
  cloudamqp plans --backend=rabbitmq`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		plans, err := c.ListPlans(backendFilter)
		if err != nil {
			fmt.Printf("Error listing plans: %v\n", err)
			return err
		}

		if len(plans) == 0 {
			fmt.Println("No plans found.")
			return nil
		}

		output, err := json.MarshalIndent(plans, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("Available plans:\n%s\n", string(output))
		return nil
	},
}

func init() {
	plansCmd.Flags().StringVar(&backendFilter, "backend", "", "Filter by specific backend software")
}
