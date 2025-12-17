package cmd

import (
	"encoding/json"
	"fmt"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var (
	updateUserID   string
	updateRole     string
	updateUserTags []string
)

var teamUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update user role and tags",
	Long: `Updates role and tags for a user in the team.

Available roles: admin, devops, member, monitor, billing manager`,
	Example: `  cloudamqp team update --user-id=uuid-here --role=admin
  cloudamqp team update --user-id=uuid-here --tags=production --tags=web-app`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		req := &client.TeamUpdateRequest{
			Role: updateRole,
			Tags: updateUserTags,
		}

		if req.Role == "" && len(req.Tags) == 0 {
			return fmt.Errorf("at least one field (role or tags) must be specified for update")
		}

		resp, err := c.UpdateTeamMember(updateUserID, req)
		if err != nil {
			fmt.Printf("Error updating team member: %v\n", err)
			return err
		}

		output, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("Team member updated:\n%s\n", string(output))
		return nil
	},
}

func init() {
	teamUpdateCmd.Flags().StringVar(&updateUserID, "user-id", "", "User ID (UUID) to update (required)")
	teamUpdateCmd.Flags().StringVar(&updateRole, "role", "", "New role to assign (admin, devops, member, monitor, billing manager)")
	teamUpdateCmd.Flags().StringSliceVar(&updateUserTags, "tags", []string{}, "New tags to assign")
	teamUpdateCmd.MarkFlagRequired("user-id")
}
