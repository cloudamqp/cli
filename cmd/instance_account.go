package cmd

import (
	"fmt"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var instanceAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage instance account operations",
	Long:  `Rotate password and API key for the instance.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		cmd.SilenceUsage = true
		return fmt.Errorf("subcommand required")
	},
}

var rotatePasswordCmd = &cobra.Command{
	Use:     "rotate-password --id <instance_id>",
	Short:   "Rotate password",
	Long:    `Initiate rotation of the user password on your instance.`,
	Example: `  cloudamqp instance account rotate-password --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		idFlag, _ := cmd.Flags().GetString("id")
		if idFlag == "" {
			return fmt.Errorf("instance ID is required. Use --id flag")
		}

		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		return c.RotatePassword(idFlag)
	},
}

var rotateInstanceAPIKeyCmd = &cobra.Command{
	Use:     "rotate-apikey --id <instance_id>",
	Short:   "Rotate Instance API key",
	Long:    `Rotate the Instance API key.`,
	Example: `  cloudamqp instance account rotate-apikey --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		idFlag, _ := cmd.Flags().GetString("id")
		if idFlag == "" {
			return fmt.Errorf("instance ID is required. Use --id flag")
		}

		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		return c.RotateInstanceAPIKey(idFlag)
	},
}

func init() {
	// Add --id flag to both account commands
	rotatePasswordCmd.Flags().StringP("id", "", "", "Instance ID (required)")
	rotatePasswordCmd.MarkFlagRequired("id")

	rotateInstanceAPIKeyCmd.Flags().StringP("id", "", "", "Instance ID (required)")
	rotateInstanceAPIKeyCmd.MarkFlagRequired("id")

	instanceAccountCmd.AddCommand(rotatePasswordCmd)
	instanceAccountCmd.AddCommand(rotateInstanceAPIKeyCmd)
}
