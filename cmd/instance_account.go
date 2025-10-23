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
}

var rotatePasswordCmd = &cobra.Command{
	Use:     "rotate-password",
	Short:   "Rotate password",
	Long:    `Initiate rotation of the user password on your instance.`,
	Example: `  cloudamqp instance manage 1234 account rotate-password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceID := currentInstanceID

		instanceAPIKey, err := getInstanceAPIKey(instanceID)
		if err != nil {
			return fmt.Errorf("failed to get instance API key: %w", err)
		}

		c := client.NewInstanceAPI(instanceAPIKey)

		err = c.RotatePassword()
		if err != nil {
			fmt.Printf("Error rotating password: %v\n", err)
			return err
		}

		fmt.Println("Password rotation initiated successfully.")
		return nil
	},
}

var rotateInstanceAPIKeyCmd = &cobra.Command{
	Use:     "rotate-apikey",
	Short:   "Rotate Instance API key",
	Long:    `Rotate the Instance API key.`,
	Example: `  cloudamqp instance manage 1234 account rotate-apikey`,
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceID := currentInstanceID

		instanceAPIKey, err := getInstanceAPIKey(instanceID)
		if err != nil {
			return fmt.Errorf("failed to get instance API key: %w", err)
		}

		c := client.NewInstanceAPI(instanceAPIKey)

		err = c.RotateInstanceAPIKey()
		if err != nil {
			fmt.Printf("Error rotating instance API key: %v\n", err)
			return err
		}

		fmt.Println("Instance API key rotation initiated successfully.")
		fmt.Printf("Warning: The local config for instance %s will need to be updated.\n", instanceID)
		fmt.Printf("Run 'cloudamqp instance get %s' to retrieve and save the new API key.\n", instanceID)
		return nil
	},
}

func init() {
	instanceAccountCmd.AddCommand(rotatePasswordCmd)
	instanceAccountCmd.AddCommand(rotateInstanceAPIKeyCmd)
}
