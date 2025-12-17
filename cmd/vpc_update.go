package cmd

import (
	"fmt"
	"strconv"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var (
	updateVPCID   string
	updateVPCName string
	updateVPCTags []string
)

var vpcUpdateCmd = &cobra.Command{
	Use:   "update --id <id>",
	Short: "Update a CloudAMQP VPC",
	Long: `Update an existing CloudAMQP VPC with new configuration.

You can update the following fields:
  --name: VPC name
  --tags: VPC tags (replaces existing tags)`,
	Example: `  cloudamqp vpc update --id 5678 --name=new-vpc-name
  cloudamqp vpc update --id 5678 --tags=production --tags=updated`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		if updateVPCID == "" {
			return fmt.Errorf("--id is required")
		}

		vpcID, err := strconv.Atoi(updateVPCID)
		if err != nil {
			return fmt.Errorf("invalid VPC ID: %v", err)
		}

		c := client.New(apiKey, Version)

		req := &client.VPCUpdateRequest{
			Name: updateVPCName,
			Tags: updateVPCTags,
		}

		if req.Name == "" && len(req.Tags) == 0 {
			return fmt.Errorf("at least one field must be specified for update")
		}

		err = c.UpdateVPC(vpcID, req)
		if err != nil {
			fmt.Printf("Error updating VPC: %v\n", err)
			return err
		}

		fmt.Printf("VPC %d updated successfully.\n", vpcID)
		return nil
	},
}

func init() {
	vpcUpdateCmd.Flags().StringVar(&updateVPCID, "id", "", "VPC ID (required)")
	vpcUpdateCmd.Flags().StringVar(&updateVPCName, "name", "", "New VPC name")
	vpcUpdateCmd.Flags().StringSliceVar(&updateVPCTags, "tags", []string{}, "New VPC tags")
	vpcUpdateCmd.MarkFlagRequired("id")
	vpcUpdateCmd.RegisterFlagCompletionFunc("id", completeVPCArgs)
}
