package cmd

import (
	"fmt"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var auditTimestamp string

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Get audit log in CSV format",
	Long:  `Returns audit log in CSV format for latest month or for month specified in params.`,
	Example: `  cloudamqp audit
  cloudamqp audit --timestamp=2022-12`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		csv, err := c.GetAuditLogCSV(auditTimestamp)
		if err != nil {
			fmt.Printf("Error getting audit log: %v\n", err)
			return err
		}

		fmt.Print(csv)
		return nil
	},
}

func init() {
	auditCmd.Flags().StringVar(&auditTimestamp, "timestamp", "", "YYYY-MM format (e.g., 2022-12)")
}
