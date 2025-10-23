package cmd

import (
	"github.com/spf13/cobra"
)

var vpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "Manage CloudAMQP VPCs",
	Long:  `Create, list, update, and delete CloudAMQP VPCs.`,
}

func init() {
	vpcCmd.AddCommand(vpcCreateCmd)
	vpcCmd.AddCommand(vpcListCmd)
	vpcCmd.AddCommand(vpcGetCmd)
	vpcCmd.AddCommand(vpcUpdateCmd)
	vpcCmd.AddCommand(vpcDeleteCmd)
}
