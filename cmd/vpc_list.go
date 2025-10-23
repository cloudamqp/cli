package cmd

import (
	"encoding/json"
	"fmt"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var vpcListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all CloudAMQP VPCs",
	Long:    `Retrieves and displays all CloudAMQP VPCs in your account.`,
	Example: `  cloudamqp vpc list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		vpcs, err := c.ListVPCs()
		if err != nil {
			fmt.Printf("Error listing VPCs: %v\n", err)
			return err
		}

		if len(vpcs) == 0 {
			fmt.Println("No VPCs found.")
			return nil
		}

		output, err := json.MarshalIndent(vpcs, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("VPCs:\n%s\n", string(output))
		return nil
	},
}
