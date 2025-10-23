package cmd

import (
	"encoding/json"
	"fmt"

	"cloudamqp-cli/client"
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

		c := client.New(apiKey)

		members, err := c.ListTeamMembers()
		if err != nil {
			fmt.Printf("Error listing team members: %v\n", err)
			return err
		}

		if len(members) == 0 {
			fmt.Println("No team members found.")
			return nil
		}

		output, err := json.MarshalIndent(members, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("Team members:\n%s\n", string(output))
		return nil
	},
}
