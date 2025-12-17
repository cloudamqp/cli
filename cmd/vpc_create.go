package cmd

import (
	"encoding/json"
	"fmt"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var (
	vpcName   string
	vpcRegion string
	vpcSubnet string
	vpcTags   []string
)

var vpcCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new CloudAMQP VPC",
	Long: `Create a new CloudAMQP VPC with the specified configuration.

Required flags:
  --name: Name of the VPC
  --region: Region identifier (e.g., amazon-web-services::us-east-1)
  --subnet: VPC subnet (e.g., 10.56.72.0/24)

Optional flags:
  --tags: VPC tags (can be specified multiple times)`,
	Example: `  cloudamqp vpc create --name=my-vpc --region=amazon-web-services::us-east-1 --subnet=10.56.72.0/24
  cloudamqp vpc create --name=my-vpc --region=amazon-web-services::us-east-1 --subnet=10.56.72.0/24 --tags=production --tags=web-app`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		req := &client.VPCCreateRequest{
			Name:   vpcName,
			Region: vpcRegion,
			Subnet: vpcSubnet,
			Tags:   vpcTags,
		}

		resp, err := c.CreateVPC(req)
		if err != nil {
			fmt.Printf("Error creating VPC: %v\n", err)
			return err
		}

		output, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("VPC created successfully:\n%s\n", string(output))
		return nil
	},
}

func init() {
	vpcCreateCmd.Flags().StringVar(&vpcName, "name", "", "Name of the VPC (required)")
	vpcCreateCmd.Flags().StringVar(&vpcRegion, "region", "", "Region identifier (required)")
	vpcCreateCmd.Flags().StringVar(&vpcSubnet, "subnet", "", "VPC subnet (required)")
	vpcCreateCmd.Flags().StringSliceVar(&vpcTags, "tags", []string{}, "VPC tags")

	vpcCreateCmd.MarkFlagRequired("name")
	vpcCreateCmd.MarkFlagRequired("region")
	vpcCreateCmd.MarkFlagRequired("subnet")

	vpcCreateCmd.RegisterFlagCompletionFunc("region", completeRegions)
}
