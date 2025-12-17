package cmd

import (
	"fmt"
	"os"

	"cloudamqp-cli/client"
	"cloudamqp-cli/internal/table"
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

		c := client.New(apiKey, Version)

		plans, err := c.ListPlans(backendFilter)
		if err != nil {
			fmt.Printf("Error listing plans: %v\n", err)
			return err
		}

		if len(plans) == 0 {
			fmt.Println("No plans found.")
			return nil
		}

		// Create table and populate data
		t := table.New(os.Stdout, "NAME", "PRICE", "BACKEND", "SHARED")
		for _, plan := range plans {
			shared := "No"
			if plan.Shared {
				shared = "Yes"
			}
			price := fmt.Sprintf("$%.2f", plan.Price)
			if plan.Price == 0 {
				price = "Free"
			}
			t.AddRow(
				plan.Name,
				price,
				plan.Backend,
				shared,
			)
		}
		t.Print()
		return nil
	},
}

func init() {
	plansCmd.Flags().StringVar(&backendFilter, "backend", "", "Filter by specific backend software")
}
