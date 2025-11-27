package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var apiKey string

func getVersionString() string {
	if Version == "dev" {
		return fmt.Sprintf("%s (development build)", Version)
	}
	// Strip leading 'v' if present to avoid double 'v' in URL
	versionTag := strings.TrimPrefix(Version, "v")
	return fmt.Sprintf("%s (%s)\nhttps://github.com/cloudamqp/cli/releases/tag/v%s", Version, BuildDate, versionTag)
}

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
	Version: getVersionString(),
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Set custom version template to match gh style
	rootCmd.SetVersionTemplate("cloudamqp version {{.Version}}\n")

	rootCmd.AddCommand(instanceCmd)
	rootCmd.AddCommand(vpcCmd)
	rootCmd.AddCommand(regionsCmd)
	rootCmd.AddCommand(plansCmd)
	rootCmd.AddCommand(teamCmd)
	rootCmd.AddCommand(auditCmd)
	rootCmd.AddCommand(completionCmd)
}
