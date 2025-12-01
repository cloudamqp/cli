package cmd

import (
	"encoding/json"
	"fmt"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var removeEmail string

var teamRemoveCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove a user from the team",
	Long:    `Removes a user from the team.`,
	Example: `  cloudamqp team remove --email=user@example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		resp, err := c.RemoveTeamMember(removeEmail)
		if err != nil {
			return err
		}

		output, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	teamRemoveCmd.Flags().StringVar(&removeEmail, "email", "", "Email address of the user to remove (required)")
	teamRemoveCmd.MarkFlagRequired("email")
}
