package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var vpcGetCmd = &cobra.Command{
	Use:     "get <id>",
	Short:   "Get details of a specific CloudAMQP VPC",
	Long:    `Retrieves and displays detailed information about a specific CloudAMQP VPC.`,
	Example: `  cloudamqp vpc get 5678`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		vpcID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid VPC ID: %v", err)
		}

		c := client.New(apiKey)

		vpc, err := c.GetVPC(vpcID)
		if err != nil {
			fmt.Printf("Error getting VPC: %v\n", err)
			return err
		}

		output, err := json.MarshalIndent(vpc, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("VPC details:\n%s\n", string(output))
		return nil
	},
}
