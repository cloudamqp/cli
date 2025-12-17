package cmd

import (
	"fmt"
	"os"
	"strings"

	"cloudamqp-cli/client"
	"cloudamqp-cli/internal/table"
	"github.com/spf13/cobra"
)

var teamListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List team members",
	Long:    `Retrieves all team members.`,
	Example: `  cloudamqp team list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		members, err := c.ListTeamMembers()
		if err != nil {
			fmt.Printf("Error listing team members: %v\n", err)
			return err
		}

		if len(members) == 0 {
			fmt.Println("No team members found.")
			return nil
		}

		// Create table and populate data
		t := table.New(os.Stdout, "EMAIL", "ROLES", "2FA")
		for _, member := range members {
			roles := strings.Join(member.Roles, ", ")
			if roles == "" {
				roles = "-"
			}
			tfa := "No"
			if member.TFAAuthEnabled {
				tfa = "Yes"
			}
			t.AddRow(member.Email, roles, tfa)
		}
		t.Print()

		return nil
	},
}
