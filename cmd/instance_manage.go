package cmd

import (
	"github.com/spf13/cobra"
)

var currentInstanceID string

var instanceManageCmd = &cobra.Command{
	Use:   "manage <instance_id>",
	Short: "Manage a specific CloudAMQP instance",
	Long: `Use instance-specific API to manage nodes, plugins, actions, and more.

This command uses the instance API key, not your main API key.
Instance API keys are automatically saved when you run 'cloudamqp instance get <id>'.`,
	Args: cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		currentInstanceID = args[0]
	},
}

func init() {
	instanceManageCmd.AddCommand(instanceNodesCmd)
	instanceManageCmd.AddCommand(instancePluginsCmd)
	instanceManageCmd.AddCommand(instanceActionsCmd)
	instanceManageCmd.AddCommand(instanceAccountCmd)
}
