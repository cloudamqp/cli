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
	Use:     "rotate-password <instance_id>",
	Short:   "Rotate password",
	Long:    `Initiate rotation of the user password on your instance.`,
	Example: `  cloudamqp instance account rotate-password 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		err = c.RotatePassword(args[0])
		if err != nil {
			fmt.Printf("Error rotating password: %v\n", err)
			return err
		}

		fmt.Println("Password rotation initiated successfully.")
		return nil
	},
}

var rotateInstanceAPIKeyCmd = &cobra.Command{
	Use:     "rotate-apikey <instance_id>",
	Short:   "Rotate Instance API key",
	Long:    `Rotate the Instance API key.`,
	Example: `  cloudamqp instance account rotate-apikey 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		err = c.RotateInstanceAPIKey(args[0])
		if err != nil {
			fmt.Printf("Error rotating instance API key: %v\n", err)
			return err
		}

		fmt.Println("Instance API key rotation initiated successfully.")
		fmt.Printf("Warning: The local config for instance %s will need to be updated.\n", args[0])
		fmt.Printf("Run 'cloudamqp instance get %s' to retrieve and save the new API key.\n", args[0])
		return nil
	},
}

func init() {
	rotatePasswordCmd.ValidArgsFunction = completeInstances
	rotateInstanceAPIKeyCmd.ValidArgsFunction = completeInstances

	instanceAccountCmd.AddCommand(rotatePasswordCmd)
	instanceAccountCmd.AddCommand(rotateInstanceAPIKeyCmd)
}
