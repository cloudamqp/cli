package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var forceDelete bool

var instanceDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a CloudAMQP instance",
	Long: `Delete a CloudAMQP instance permanently.

WARNING: This action cannot be undone. All data will be lost.`,
	Example: `  cloudamqp instance delete 1234
  cloudamqp instance delete 1234 --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		instanceID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid instance ID: %v", err)
		}

		if !forceDelete {
			fmt.Printf("Are you sure you want to delete instance %d? This action cannot be undone. (y/N): ", instanceID)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read confirmation: %v", err)
			}

			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Delete operation cancelled.")
				return nil
			}
		}

		c := client.New(apiKey)

		err = c.DeleteInstance(instanceID)
		if err != nil {
			fmt.Printf("Error deleting instance: %v\n", err)
			return err
		}

		fmt.Printf("Instance %d deleted successfully.\n", instanceID)
		return nil
	},
}

func init() {
	instanceDeleteCmd.Flags().BoolVar(&forceDelete, "force", false, "Skip confirmation prompt")
}
