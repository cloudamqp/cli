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

var (
	deleteVPCID    string
	forceDeleteVPC bool
)

var vpcDeleteCmd = &cobra.Command{
	Use:   "delete --id <id>",
	Short: "Delete a CloudAMQP VPC",
	Long: `Delete a CloudAMQP VPC permanently.

WARNING: This action cannot be undone. All instances in the VPC must be deleted first.`,
	Example: `  cloudamqp vpc delete --id 5678
  cloudamqp vpc delete --id 5678 --force`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		if deleteVPCID == "" {
			return fmt.Errorf("--id is required")
		}

		vpcID, err := strconv.Atoi(deleteVPCID)
		if err != nil {
			return fmt.Errorf("invalid VPC ID: %v", err)
		}

		if !forceDeleteVPC {
			fmt.Printf("Are you sure you want to delete VPC %d? This action cannot be undone. (y/N): ", vpcID)
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

		err = c.DeleteVPC(vpcID)
		if err != nil {
			fmt.Printf("Error deleting VPC: %v\n", err)
			return err
		}

		fmt.Printf("VPC %d deleted successfully.\n", vpcID)
		return nil
	},
}

func init() {
	vpcDeleteCmd.Flags().StringVar(&deleteVPCID, "id", "", "VPC ID (required)")
	vpcDeleteCmd.Flags().BoolVar(&forceDeleteVPC, "force", false, "Skip confirmation prompt")
	vpcDeleteCmd.MarkFlagRequired("id")
	vpcDeleteCmd.RegisterFlagCompletionFunc("id", completeVPCArgs)
}
