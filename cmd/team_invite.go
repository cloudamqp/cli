package cmd

import (
	"encoding/json"
	"fmt"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var (
	inviteEmail string
	inviteRole  string
	inviteTags  []string
)

var teamInviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Invite a new user to the team",
	Long: `Invites a user to join the team with specified role.

Available roles: admin, devops, member, monitor, billing manager
Default role: member`,
	Example: `  cloudamqp team invite --email=user@example.com
  cloudamqp team invite --email=user@example.com --role=admin --tags=production`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		req := &client.TeamInviteRequest{
			Email: inviteEmail,
			Role:  inviteRole,
			Tags:  inviteTags,
		}

		resp, err := c.InviteTeamMember(req)
		if err != nil {
			fmt.Printf("Error inviting team member: %v\n", err)
			return err
		}

		output, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("Team member invited:\n%s\n", string(output))
		return nil
	},
}

func init() {
	teamInviteCmd.Flags().StringVar(&inviteEmail, "email", "", "Email address of the user to invite (required)")
	teamInviteCmd.Flags().StringVar(&inviteRole, "role", "member", "Role to assign (admin, devops, member, monitor, billing manager)")
	teamInviteCmd.Flags().StringSliceVar(&inviteTags, "tags", []string{}, "Tags to assign")
	teamInviteCmd.MarkFlagRequired("email")
}
