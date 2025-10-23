package cmd

import (
	"github.com/spf13/cobra"
)

var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Manage CloudAMQP instances",
	Long:  `Create, list, update, and delete CloudAMQP instances.`,
}

func init() {
	instanceCmd.AddCommand(instanceCreateCmd)
	instanceCmd.AddCommand(instanceListCmd)
	instanceCmd.AddCommand(instanceGetCmd)
	instanceCmd.AddCommand(instanceUpdateCmd)
	instanceCmd.AddCommand(instanceDeleteCmd)
	instanceCmd.AddCommand(instanceResizeCmd)
	instanceCmd.AddCommand(instanceConfigCmd)
	instanceCmd.AddCommand(instanceNodesCmd)
	instanceCmd.AddCommand(instanceActionsCmd)
	instanceCmd.AddCommand(instanceAccountCmd)
	instanceCmd.AddCommand(instancePluginsCmd)
	instanceCmd.AddCommand(instanceManageCmd)
}
