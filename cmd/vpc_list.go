package cmd

import (
	"fmt"
	"strconv"

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

		c := client.New(apiKey, Version)

		vpcs, err := c.ListVPCs()
		if err != nil {
			fmt.Printf("Error listing VPCs: %v\n", err)
			return err
		}

		if len(vpcs) == 0 {
			fmt.Println("No VPCs found.")
			return nil
		}

		p, err := getPrinter(cmd)
		if err != nil {
			return err
		}

		headers := []string{"ID", "NAME", "SUBNET", "REGION"}
		rows := make([][]string, len(vpcs))
		for i, vpc := range vpcs {
			rows[i] = []string{
				strconv.Itoa(vpc.ID),
				vpc.Name,
				vpc.Subnet,
				vpc.Region,
			}
		}
		p.PrintRecords(headers, rows)

		return nil
	},
}
