package cmd

import (
	"fmt"
	"strconv"
	"strings"

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

		p, err := getPrinter(cmd)
		if err != nil {
			return err
		}

		instanceIDs := make([]string, len(vpc.Instances))
		for i, id := range vpc.Instances {
			instanceIDs[i] = strconv.Itoa(id)
		}

		p.PrintRecord(
			[]string{"ID", "NAME", "REGION", "SUBNET", "PLAN", "TAGS", "INSTANCES"},
			[]string{
				strconv.Itoa(vpc.ID),
				vpc.Name,
				vpc.Region,
				vpc.Subnet,
				vpc.Plan,
				strings.Join(vpc.Tags, ","),
				strings.Join(instanceIDs, ","),
			},
		)

		return nil
	},
}

func init() {
	vpcGetCmd.Flags().StringP("id", "", "", "VPC ID (required)")
	vpcGetCmd.MarkFlagRequired("id")
	vpcGetCmd.RegisterFlagCompletionFunc("id", completeVPCIDFlag)
}
