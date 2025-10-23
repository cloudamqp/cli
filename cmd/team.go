package cmd

import (
	"github.com/spf13/cobra"
)

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Manage team members",
	Long:  `List, invite, update, and remove team members.`,
}

func init() {
	teamCmd.AddCommand(teamListCmd)
	teamCmd.AddCommand(teamInviteCmd)
	teamCmd.AddCommand(teamRemoveCmd)
	teamCmd.AddCommand(teamUpdateCmd)
}
