package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var vpcGetCmd = &cobra.Command{
	Use:     "get --id <id>",
	Short:   "Get details of a specific CloudAMQP VPC",
	Long:    `Retrieves and displays detailed information about a specific CloudAMQP VPC.`,
	Example: `  cloudamqp vpc get --id 5678`,
	RunE: func(cmd *cobra.Command, args []string) error {
		idFlag, _ := cmd.Flags().GetString("id")
		if idFlag == "" {
			return fmt.Errorf("VPC ID is required. Use --id flag")
		}

		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		vpcID, err := strconv.Atoi(idFlag)
		if err != nil {
			return fmt.Errorf("invalid VPC ID: %v", err)
		}

		c := client.New(apiKey, Version)

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

func init() {
	vpcGetCmd.Flags().StringP("id", "", "", "VPC ID (required)")
	vpcGetCmd.MarkFlagRequired("id")
	vpcGetCmd.RegisterFlagCompletionFunc("id", completeVPCIDFlag)
}
