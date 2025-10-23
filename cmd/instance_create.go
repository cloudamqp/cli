package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var (
	instanceName      string
	instancePlan      string
	instanceRegion    string
	instanceTags      []string
	instanceVPCSubnet string
	instanceVPCID     string
)

var instanceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new CloudAMQP instance",
	Long: `Create a new CloudAMQP instance with the specified configuration.

Required flags:
  --name: Name of the instance
  --plan: Subscription plan (e.g., lemming, bunny-1, rabbit-1)
  --region: Region identifier (e.g., amazon-web-services::us-east-1)

Optional flags:
  --tags: Instance tags (can be specified multiple times)
  --vpc-subnet: VPC subnet for dedicated VPC
  --vpc-id: ID of existing VPC to add instance to`,
	Example: `  cloudamqp instance create --name=my-instance --plan=bunny-1 --region=amazon-web-services::us-east-1
  cloudamqp instance create --name=my-instance --plan=bunny-1 --region=amazon-web-services::us-east-1 --tags=production --tags=web-app`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		req := &client.InstanceCreateRequest{
			Name:   instanceName,
			Plan:   instancePlan,
			Region: instanceRegion,
			Tags:   instanceTags,
		}

		if instanceVPCSubnet != "" {
			req.VPCSubnet = instanceVPCSubnet
		}

		if instanceVPCID != "" {
			vpcID, err := strconv.Atoi(instanceVPCID)
			if err != nil {
				return fmt.Errorf("invalid VPC ID: %v", err)
			}
			req.VPCID = &vpcID
		}

		resp, err := c.CreateInstance(req)
		if err != nil {
			fmt.Printf("Error creating instance: %v\n", err)
			return err
		}

		output, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("Instance created successfully:\n%s\n", string(output))
		return nil
	},
}

func init() {
	instanceCreateCmd.Flags().StringVar(&instanceName, "name", "", "Name of the instance (required)")
	instanceCreateCmd.Flags().StringVar(&instancePlan, "plan", "", "Subscription plan (required)")
	instanceCreateCmd.Flags().StringVar(&instanceRegion, "region", "", "Region identifier (required)")
	instanceCreateCmd.Flags().StringSliceVar(&instanceTags, "tags", []string{}, "Instance tags")
	instanceCreateCmd.Flags().StringVar(&instanceVPCSubnet, "vpc-subnet", "", "VPC subnet")
	instanceCreateCmd.Flags().StringVar(&instanceVPCID, "vpc-id", "", "VPC ID")

	instanceCreateCmd.MarkFlagRequired("name")
	instanceCreateCmd.MarkFlagRequired("plan")
	instanceCreateCmd.MarkFlagRequired("region")
}
