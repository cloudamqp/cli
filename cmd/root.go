package cmd

import (
	"github.com/spf13/cobra"
)

var apiKey string

var rootCmd = &cobra.Command{
	Use:   "cloudamqp",
	Short: "CloudAMQP CLI for managing instances and VPCs",
	Long: `A command line interface for the CloudAMQP API.
Use this tool to create, manage, and delete CloudAMQP instances and VPCs.

API Key Configuration:
The CLI will look for your API key in the following order:
1. CLOUDAMQP_APIKEY environment variable
2. ~/.cloudamqprc file (JSON format)
3. If neither exists, you will be prompted to enter it

Instance API keys are automatically saved when using 'instance get' command.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(instanceCmd)
	rootCmd.AddCommand(vpcCmd)
	rootCmd.AddCommand(regionsCmd)
	rootCmd.AddCommand(plansCmd)
	rootCmd.AddCommand(teamCmd)
	rootCmd.AddCommand(auditCmd)
	rootCmd.AddCommand(rotateKeyCmd)
}
