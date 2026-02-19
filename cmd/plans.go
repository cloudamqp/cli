package cmd

import (
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

		p, err := getPrinter(cmd)
		if err != nil {
			return err
		}

		headers := []string{"NAME", "PRICE", "BACKEND", "SHARED"}
		rows := make([][]string, len(plans))
		for i, plan := range plans {
			shared := "No"
			if plan.Shared {
				shared = "Yes"
			}
			price := fmt.Sprintf("$%.2f", plan.Price)
			if plan.Price == 0 {
				price = "Free"
			}
			rows[i] = []string{plan.Name, price, plan.Backend, shared}
		}
		p.PrintRecords(headers, rows)

		return nil
	},
}

func init() {
	plansCmd.Flags().StringVar(&backendFilter, "backend", "", "Filter by specific backend software")
}
